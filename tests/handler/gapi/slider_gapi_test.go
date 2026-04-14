package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	slider_cache "ecommerce/internal/cache/slider"
	"ecommerce/internal/handler/gapi"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"net"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type SliderGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	repos       *repository.Repositories
	client      pb.SliderServiceClient
	conn        *grpc.ClientConn
	listener    *bufconn.Listener
}

func (s *SliderGapiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	// DB Setup
	pool, err := pgxpool.New(s.ts.Ctx, s.ts.DBURL)
	s.Require().NoError(err)
	s.dbPool = pool

	// Redis Setup
	opt, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opt)

	// Repositories
	queries := db.New(pool)
	s.repos = repository.NewRepositories(queries)

	// Logging & Observability
	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	l, err := logger.NewLogger("test-slider-gapi", lp)
	s.Require().NoError(err)

	obs, err := observability.NewObservability("test-slider-gapi", l)
	s.Require().NoError(err)

	// Cache Setup
	cacheMetrics, err := observability.NewCacheMetrics("test-slider-gapi")
	s.Require().NoError(err)
	cacheStore := cache.NewCacheStore(s.redisClient, l, cacheMetrics)
	slCache := slider_cache.NewSliderMencache(cacheStore)

	// Service
	slService := service.NewSliderService(service.SliderServiceDeps{
		SliderRepository: s.repos.Slider,
		Logger:           l,
		Observability:    obs,
		Cache:            slCache,
	})

	// gRPC Server Setup
	s.listener = bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	hGrpc := gapi.NewSliderHandleGrpc(slService)
	pb.RegisterSliderServiceServer(server, hGrpc)

	go func() {
		if err := server.Serve(s.listener); err != nil {
			panic(err)
		}
	}()

	// gRPC Client Setup
	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return s.listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewSliderServiceClient(conn)
}

func (s *SliderGapiTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.dbPool != nil {
		s.dbPool.Close()
	}
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
}

func (s *SliderGapiTestSuite) TestSliderGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pb.CreateSliderRequest{
		Name:  "Gapi Test Slider",
		Image: "path/to/gapi_image.jpg",
	}

	resCreate, err := s.client.Create(ctx, createReq)
	s.NoError(err)
	s.NotNil(resCreate)
	s.Equal(createReq.Name, resCreate.Data.Name)
	sliderID := resCreate.Data.Id

	// 2. Find All
	resAll, err := s.client.FindAll(ctx, &pb.FindAllSliderRequest{
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(resAll)
	s.GreaterOrEqual(resAll.Pagination.TotalRecords, int32(1))

	// 3. Update
	updateReq := &pb.UpdateSliderRequest{
		Id:    sliderID,
		Name:  "Updated Gapi Test Slider",
		Image: "path/to/updated_gapi_image.jpg",
	}
	resUpdate, err := s.client.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(resUpdate)
	s.Equal(updateReq.Name, resUpdate.Data.Name)

	// 4. Trash
	resTrash, err := s.client.TrashedSlider(ctx, &pb.FindByIdSliderRequest{Id: sliderID})
	s.NoError(err)
	s.NotNil(resTrash)
	s.NotNil(resTrash.Data.DeletedAt)

	// 5. Restore
	resRestore, err := s.client.RestoreSlider(ctx, &pb.FindByIdSliderRequest{Id: sliderID})
	s.NoError(err)
	s.NotNil(resRestore)
	s.Nil(resRestore.Data.DeletedAt)

	// 6. Delete Permanent
	// Trash again
	_, err = s.client.TrashedSlider(ctx, &pb.FindByIdSliderRequest{Id: sliderID})
	s.NoError(err)

	resDelete, err := s.client.DeleteSliderPermanent(ctx, &pb.FindByIdSliderRequest{Id: sliderID})
	s.NoError(err)
	s.NotNil(resDelete)
	s.Equal("success", resDelete.Status)
}

func TestSliderGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(SliderGapiTestSuite))
}
