package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	category_cache "ecommerce/internal/cache/category"
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
)

type CategoryGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.CategoryServiceClient
	conn        *grpc.ClientConn
}

func (s *CategoryGapiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-gapi", lp)
	obs, _ := observability.NewObservability("test-gapi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	catCache := category_cache.NewCategoryMencache(cacheStore)

	categoryService := service.NewCategoryService(service.CategoryServiceDeps{
		CategoryRepository: repos.Category,
		Logger:             log,
		Observability:      obs,
		Cache:              catCache,
	})

	// Start gRPC Server
	categoryHandler := gapi.NewCategoryHandleGrpc(categoryService)
	server := grpc.NewServer()
	pb.RegisterCategoryServiceServer(server, categoryHandler)
	s.grpcServer = server

	lis, err := net.Listen("tcp", ":0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	// Create Client
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn
	s.client = pb.NewCategoryServiceClient(conn)
}

func (s *CategoryGapiTestSuite) TearDownSuite() {
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

func (s *CategoryGapiTestSuite) TestCategoryLifecycle() {
	ctx := context.Background()

	// 1. Create
	slug := "electronics-gapi"
	createReq := &pb.CreateCategoryRequest{
		Name:          "Electronics Gapi",
		Description:   "Electronic gadgets via gRPC",
		SlugCategory:  slug,
		ImageCategory: "electronics_gapi.jpg",
	}


	res, err := s.client.Create(ctx, createReq)
	s.NoError(err)
	s.Equal(createReq.Name, res.Data.Name)
	categoryID := res.Data.Id

	// 2. Find By ID
	found, err := s.client.FindById(ctx, &pb.FindByIdCategoryRequest{Id: categoryID})
	s.NoError(err)
	s.Equal(categoryID, found.Data.Id)

	// 3. Update
	updateReq := &pb.UpdateCategoryRequest{
		CategoryId:    categoryID,
		Name:          "Electronics Gapi v2",
		Description:   "Updated via gRPC",
		ImageCategory: "electronics_gapi_v2.jpg",
	}
	updated, err := s.client.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.Name, updated.Data.Name)

	// 4. Trash
	_, err = s.client.TrashedCategory(ctx, &pb.FindByIdCategoryRequest{Id: categoryID})
	s.NoError(err)

	// 5. Restore
	_, err = s.client.RestoreCategory(ctx, &pb.FindByIdCategoryRequest{Id: categoryID})
	s.NoError(err)

	// 6. Delete Permanent
	// Need to trash again before permanent delete
	_, err = s.client.TrashedCategory(ctx, &pb.FindByIdCategoryRequest{Id: categoryID})
	s.NoError(err)

	delRes, err := s.client.DeleteCategoryPermanent(ctx, &pb.FindByIdCategoryRequest{Id: categoryID})
	s.NoError(err)
	s.Equal("success", delRes.Status)

}

func TestCategoryGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryGapiTestSuite))
}
