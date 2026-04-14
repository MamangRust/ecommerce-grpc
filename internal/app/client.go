package app

import (
	"context"
	_ "ecommerce/docs"
	"ecommerce/internal/cache"
	"ecommerce/internal/handler/api"
	response_api "ecommerce/internal/mapper"
	"ecommerce/internal/middlewares"
	"ecommerce/pkg/auth"
	"ecommerce/pkg/dotenv"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/pkg/otel"
	"ecommerce/pkg/upload_image"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	defaultClientAddr             = "localhost:50051"
	defaultServerPort             = ":5000"
	defaultWindowSizeClient       = 16 * 1024 * 1024
	defaultKeepaliveTimeClient    = 20 * time.Second
	defaultKeepaliveTimeoutClient = 5 * time.Second
	shutdownTimeoutClient         = 10 * time.Second
)

var (
	addr = flag.String("addr", defaultClientAddr, "the gRPC server address to connect to")
)

// @title Ecommerce gRPC Gateway
// @version 1.0
// @description REST Gateway for gRPC based Ecommerce service
// @host localhost:5000
// @BasePath /api/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
type Client struct {
	Logger       logger.LoggerInterface
	Echo         *echo.Echo
	GRPCConn     *grpc.ClientConn
	TokenManager *auth.Manager
	Telemetry    *otel.Telemetry
	Config       *ClientConfig
	Redis        *redis.Client

	cancelTasks context.CancelFunc
	tasksDone   []<-chan struct{}
}

type ClientConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OtelEndpoint   string
	GRPCAddr       string
	ServerPort     string
	AllowedOrigins []string
}

type CacheManager struct {
	cache  *cache.CacheStore
	logger logger.LoggerInterface
}

func NewCacheManager(cache *cache.CacheStore, logger logger.LoggerInterface) *CacheManager {
	return &CacheManager{
		cache:  cache,
		logger: logger,
	}
}

func (cm *CacheManager) StartMonitoring(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		ticker := time.NewTicker(monitoringInterval)
		defer ticker.Stop()

		cm.logger.Info("Cache monitoring task started",
			zap.Duration("interval", monitoringInterval),
		)

		for {
			select {
			case <-ctx.Done():
				cm.logger.Info("Cache monitoring task stopped")
				return
			case <-ticker.C:
				cm.monitor(ctx)
			}
		}
	}()

	return done
}

func (cm *CacheManager) monitor(ctx context.Context) {
	refCount := cm.cache.GetRefCount()

	stats, err := cm.cache.GetStats(ctx)
	if err != nil {
		cm.logger.Error("Failed to get cache stats", zap.Error(err))
		return
	}

	logLevel := zap.InfoLevel
	if refCount > cacheRefCountThreshold {
		logLevel = zap.WarnLevel
	}

	if ce := cm.logger.Check(logLevel, "Cache statistics"); ce != nil {
		ce.Write(
			zap.Int64("ref_count", refCount),
			zap.Int64("total_keys", stats.TotalKeys),
			zap.Float64("hit_rate", stats.HitRate),
			zap.String("memory_used", stats.MemoryUsedHuman),
			zap.Bool("high_ref_count", refCount > cacheRefCountThreshold),
		)
	}
}

func (cm *CacheManager) StartCleanup(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		cm.logger.Info("Cache cleanup task started",
			zap.Duration("interval", cleanupInterval),
		)

		for {
			select {
			case <-ctx.Done():
				cm.logger.Info("Cache cleanup task stopped")
				return
			case <-ticker.C:
				cm.cleanup(ctx)
			}
		}
	}()

	return done
}

func (cm *CacheManager) cleanup(ctx context.Context) {
	cm.logger.Info("Starting periodic cache cleanup")

	statsBefore, err := cm.cache.GetStats(ctx)
	if err != nil {
		cm.logger.Error("Failed to get cache stats before cleanup", zap.Error(err))
		statsBefore = nil
	}

	scanned, err := cm.cache.ClearExpired(ctx)
	if err != nil {
		cm.logger.Error("Cache cleanup failed", zap.Error(err))
		return
	}

	statsAfter, err := cm.cache.GetStats(ctx)
	if err != nil {
		cm.logger.Error("Failed to get cache stats after cleanup", zap.Error(err))
		statsAfter = nil
	}

	logFields := []zap.Field{
		zap.Int64("scanned_keys", scanned),
		zap.Int64("ref_count", cm.cache.GetRefCount()),
	}

	if statsBefore != nil && statsAfter != nil {
		keysRemoved := statsBefore.TotalKeys - statsAfter.TotalKeys
		logFields = append(logFields,
			zap.Int64("keys_before", statsBefore.TotalKeys),
			zap.Int64("keys_after", statsAfter.TotalKeys),
			zap.Int64("keys_removed", keysRemoved),
			zap.String("memory_before", statsBefore.MemoryUsedHuman),
			zap.String("memory_after", statsAfter.MemoryUsedHuman),
		)
	}

	cm.logger.Info("Cache cleanup completed", logFields...)
}

func NewClient(cfg *ClientConfig) (*Client, error) {
	flag.Parse()

	if err := initPyroscope(); err != nil {
		log.Fatal("Failed to initialize pyroscope:", err)
	}

	if cfg == nil {
		cfg = &ClientConfig{
			ServiceName:    "client",
			ServiceVersion: "1.0.0",
			Environment:    getEnv("ENVIRONMENT", "production"),
			OtelEndpoint:   getEnv("OTEL_ENDPOINT", "otel-collector:4317"),
			GRPCAddr:       *addr,
			ServerPort:     defaultServerPort,
			AllowedOrigins: []string{"http://localhost:1420"},
		}
	}

	telemetry, err := initTelemetryClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize telemetry: %w", err)
	}

	cacheMetrics, err := observability.NewCacheMetrics("cache")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache metrics: %w", err)
	}

	logger, err := logger.NewLogger(cfg.ServiceName, telemetry.GetLogger())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to create token manager: %w", err)
	}

	grpcConn, err := connectToGRPC(cfg.GRPCAddr, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	echoServer := createEchoServer(cfg)

	mapper := response_api.NewResponseApiMapper()
	imageUpload := upload_image.NewImageUpload(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisClient, err := initRedisClient(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	cacheStore := cache.NewCacheStore(redisClient, logger, cacheMetrics)

	tasksCtx, cancelTasks := context.WithCancel(context.Background())
	cacheManager := NewCacheManager(cacheStore, logger)

	tasksDone := []<-chan struct{}{
		cacheManager.StartMonitoring(tasksCtx),
		cacheManager.StartCleanup(tasksCtx),
	}

	handlerDeps := api.Deps{
		Conn:    grpcConn,
		Token:   tokenManager,
		E:       echoServer,
		Logger:  logger,
		Image:   imageUpload,
		Mapping: mapper,
		Cache:   cacheStore,
	}
	api.NewHandler(handlerDeps)

	client := &Client{
		Logger:       logger,
		Echo:         echoServer,
		GRPCConn:     grpcConn,
		TokenManager: tokenManager,
		Telemetry:    telemetry,
		Config:       cfg,
		Redis:        redisClient,
		cancelTasks:  cancelTasks,
		tasksDone:    tasksDone,
	}

	logger.Info("Client initialized successfully",
		zap.String("service", cfg.ServiceName),
		zap.String("version", cfg.ServiceVersion),
		zap.String("grpc_addr", cfg.GRPCAddr),
		zap.String("server_port", cfg.ServerPort),
	)

	return client, nil
}

func (c *Client) Run() error {
	defer c.Cleanup()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	errChan := make(chan error, 1)
	go func() {
		c.Logger.Info("HTTP server starting",
			zap.String("port", c.Config.ServerPort),
			zap.String("swagger", "http://localhost"+c.Config.ServerPort+"/swagger/index.html"),
		)
		if err := c.Echo.Start(c.Config.ServerPort); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start server: %w", err)
		}
	}()

	select {
	case sig := <-quit:
		c.Logger.Info("Received shutdown signal",
			zap.String("signal", sig.String()),
		)
	case err := <-errChan:
		c.Logger.Error("Server error", zap.Error(err))
		return err
	}

	return c.gracefulShutdown()
}

func (c *Client) gracefulShutdown() error {
	c.Logger.Info("Starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := c.Echo.Shutdown(ctx); err != nil {
		c.Logger.Error("Echo shutdown error", zap.Error(err))
		return fmt.Errorf("failed to shutdown echo server: %w", err)
	}

	c.Logger.Info("HTTP server stopped gracefully")
	return nil
}

func (c *Client) Cleanup() {
	c.Logger.Info("Cleaning up resources...")

	if c.cancelTasks != nil {
		c.Logger.Info("Stopping background tasks...")
		c.cancelTasks()

		for i, done := range c.tasksDone {
			c.Logger.Debug("Waiting for background task to complete",
				zap.Int("task_index", i),
			)
			<-done
		}
		c.Logger.Info("All background tasks stopped")
	}

	if c.Redis != nil {
		if err := c.Redis.Close(); err != nil {
			c.Logger.Error("Failed to close Redis connection", zap.Error(err))
		} else {
			c.Logger.Info("Redis connection closed")
		}
	}

	if c.GRPCConn != nil {
		if err := c.GRPCConn.Close(); err != nil {
			c.Logger.Error("Failed to close gRPC connection", zap.Error(err))
		} else {
			c.Logger.Info("gRPC connection closed")
		}
	}

	if c.Telemetry != nil {
		if err := c.Telemetry.Shutdown(context.Background()); err != nil {
			c.Logger.Error("Failed to shutdown telemetry", zap.Error(err))
		} else {
			c.Logger.Info("Telemetry shutdown successfully")
		}
	}

	if c.Logger != nil {
		if err := c.Logger.Sync(); err != nil {
			// Ignore sync errors
			_ = err
		}
	}

	c.Logger.Info("Cleanup completed")
}

func initPyroscope() error {
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "client",
		ServerAddress:   os.Getenv("PYROSCOPE_SERVER"),

		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},

		Tags: map[string]string{
			"service": "grpc-client-echo",
			"env":     os.Getenv("ENV"),
			"version": os.Getenv("VERSION"),
		},
	})
	return err
}

func initTelemetryClient(cfg *ClientConfig) (*otel.Telemetry, error) {
	telemetry := otel.NewTelemetry(otel.Config{
		ServiceName:            cfg.ServiceName,
		ServiceVersion:         cfg.ServiceVersion,
		Environment:            cfg.Environment,
		Endpoint:               cfg.OtelEndpoint,
		Insecure:               true,
		EnableRuntimeMetrics:   true,
		RuntimeMetricsInterval: 15 * time.Second,
	})

	if err := telemetry.Init(context.Background()); err != nil {
		return nil, err
	}

	return telemetry, nil
}

func initRedisClient(ctx context.Context, logger logger.LoggerInterface) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST_CLIENT"), viper.GetString("REDIS_PORT_CLIENT")),
		Password:     viper.GetString("REDIS_PASSWORD_CLIENT"),
		DB:           viper.GetInt("REDIS_DB_CLIENT"),
		DialTimeout:  redisDialTimeout,
		ReadTimeout:  redisReadTimeout,
		WriteTimeout: redisWriteTimeout,
		PoolSize:     redisPoolSize,
		MinIdleConns: redisMinIdleConns,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	logger.Info("Redis connection established",
		zap.String("addr", fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST_CLIENT"), viper.GetString("REDIS_PORT_CLIENT"))),
		zap.Int("db", viper.GetInt("REDIS_DB_CLIENT")),
	)

	return client, nil
}

func connectToGRPC(addr string, logger logger.LoggerInterface) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Connecting to gRPC server", zap.String("address", addr))

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithInitialConnWindowSize(defaultWindowSizeClient),
		grpc.WithInitialWindowSize(defaultWindowSizeClient),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                defaultKeepaliveTimeClient,
			Timeout:             defaultKeepaliveTimeoutClient,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	select {
	case <-ctx.Done():
		conn.Close()
		return nil, fmt.Errorf("connection timeout: %w", ctx.Err())
	default:
		logger.Info("gRPC connection established", zap.String("address", addr))
		return conn, nil
	}
}

func createEchoServer(cfg *ClientConfig) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(createLoggerMiddleware())
	e.Use(middlewares.PyroscopeMiddleware())
	e.Use(createCORSMiddleware(cfg.AllowedOrigins))
	e.Use(middleware.Gzip())
	e.Use(createSecureMiddleware())

	middlewares.WebSecurityConfig(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "healthy",
			"service": cfg.ServiceName,
			"version": cfg.ServiceVersion,
			"time":    time.Now().UTC(),
		})
	})

	return e
}

func createLoggerMiddleware() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","status":${status},` +
			`"error":"${error}","latency":${latency},"latency_human":"${latency_human}",` +
			`"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	})
}

func createCORSMiddleware(allowedOrigins []string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	})
}

func createSecureMiddleware() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		HSTSMaxAge:         31536000,
		ContentSecurityPolicy: "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' https://cdnjs.cloudflare.com; " +
			"style-src 'self' 'unsafe-inline' https://cdnjs.cloudflare.com; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data: https://cdnjs.cloudflare.com; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none';",
	})
}
