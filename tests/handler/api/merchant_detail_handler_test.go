package api_test

import (
	"bytes"
	"context"
	"ecommerce/internal/cache"
	api_merchantdetail_cache "ecommerce/internal/cache/api/merchant_detail"
	merchantdetail_cache "ecommerce/internal/cache/merchant_detail"
	"ecommerce/internal/handler/api"
	"ecommerce/internal/handler/gapi"
	response_api "ecommerce/internal/mapper"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/pkg/upload_image"
	"ecommerce/tests"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MerchantDetailApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	echo        *echo.Echo
	client      pb.MerchantDetailServiceClient
	conn        *grpc.ClientConn
	merchantID  int32
}

func (s *MerchantDetailApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	queries := db.New(pool)
	repos := repository.NewRepositories(queries)

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	log, _ := logger.NewLogger("test-merchant-detail-api", lp)
	obs, _ := observability.NewObservability("test-merchant-detail-api", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-merchant-detail-api")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	
	detailCacheSrv := merchantdetail_cache.NewMerchantDetailMencache(cacheStore)
	detailCacheApi := api_merchantdetail_cache.NewMerchantDetailMencache(cacheStore)

	detailService := service.NewMerchantDetailService(service.MerchantDetailServiceDeps{
		MerchantDetailRepository:     repos.MerchantDetail,
		MerchantSocialLinkRepository: repos.MerchantSocialLink,
		Logger:                       log,
		Cache:                        detailCacheSrv,
		Observability:                obs,
	})

	// Start gRPC Server
	detailGapi := gapi.NewMerchantDetailHandleGrpc(detailService)
	server := grpc.NewServer()
	pb.RegisterMerchantDetailServiceServer(server, detailGapi)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// gRPC Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewMerchantDetailServiceClient(conn)

	// Echo Setup
	s.echo = echo.New()
	mapping := response_api.NewMerchantDetailResponseMapper()
	mappingMerchant := response_api.NewMerchantResponseMapper()
	imgUpload := upload_image.NewImageUpload(log)
	apiHandler := errors.NewApiHandler(obs, log)

	api.NewHandlerMerchantDetail(s.echo, s.client, log, mapping, mappingMerchant, imgUpload, apiHandler, detailCacheApi)

	// Prerequisite: Create Merchant
	ctx := context.Background()
	_ = pool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING user_id",
		"Detail", "Api", "detail-api@example.com", "password").Scan(new(int32))
	err = pool.QueryRow(ctx, "INSERT INTO merchants (user_id, name) VALUES (1, $1) RETURNING merchant_id", "Detail Merchant").Scan(&s.merchantID)
	s.Require().NoError(err)
}

func (s *MerchantDetailApiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *MerchantDetailApiTestSuite) TestMerchantDetailApiLifecycle() {
	// 1. Create (Multipart)
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("merchant_id", fmt.Sprintf("%d", s.merchantID))
	_ = writer.WriteField("display_name", "Display Name")
	_ = writer.WriteField("short_description", "Short Description")
	_ = writer.WriteField("website_url", "https://website.com")
	
	// Social links JSON
	socialLinks := []map[string]interface{}{
		{"platform": "Instagram", "url": "https://insta.am/test", "merchant_detail_id": 0},
	}
	socialBytes, _ := json.Marshal(socialLinks)
	_ = writer.WriteField("social_links", string(socialBytes))

	part, _ := writer.CreateFormFile("cover_image_url", "cover.jpg")
	_, _ = part.Write([]byte("fake cover content"))
	
	part2, _ := writer.CreateFormFile("logo_url", "logo.jpg")
	_, _ = part2.Write([]byte("fake logo content"))
	
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/merchant-detail/create", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var createRes map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &createRes)
	data := createRes["data"].(map[string]interface{})
	detailID := int(data["id"].(float64))

	// 2. Find By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/merchant-detail/%d", detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. Update (Multipart)
	body = new(bytes.Buffer)
	writer = multipart.NewWriter(body)
	_ = writer.WriteField("display_name", "Updated Display Name")
	_ = writer.WriteField("short_description", "Updated Short Description")
	_ = writer.WriteField("website_url", "https://updated-website.com")
	
	// Social links
	socialLinks = []map[string]interface{}{
		{"id": 1, "platform": "Twitter", "url": "https://twitter.com/test", "merchant_detail_id": detailID},
	}
	socialBytes, _ = json.Marshal(socialLinks)
	_ = writer.WriteField("social_links", string(socialBytes))
	_ = writer.Close()

	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-detail/update/%d", detailID), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/merchant-detail/trashed/%d", detailID), nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestMerchantDetailApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDetailApiTestSuite))
}
