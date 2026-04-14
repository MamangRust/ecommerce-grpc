package gapi_test

import (
	"context"
	"ecommerce/internal/cache"
	product_cache "ecommerce/internal/cache/product"
	"ecommerce/internal/handler/gapi"
	"ecommerce/internal/pb"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"ecommerce/tests"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductGapiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	dbPool      *pgxpool.Pool
	redisClient *redis.Client
	grpcServer  *grpc.Server
	client      pb.ProductServiceClient
	conn        *grpc.ClientConn
	
	// IDs for cleanup and reference
	merchantID int32
	categoryID int32
}

func (s *ProductGapiTestSuite) SetupSuite() {
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
	log, _ := logger.NewLogger("test-product-gapi", lp)
	obs, _ := observability.NewObservability("test-product-gapi", log)
	
	cacheMetrics, _ := observability.NewCacheMetrics("test-product-gapi")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	prodCache := product_cache.NewProductMencache(cacheStore)

	productService := service.NewProductService(service.ProductServiceDeps{
		CategoryRepository: repos.Category,
		MerchantRepository: repos.Merchant,
		ProductRepository:  repos.Product,
		Logger:            log,
		Observability:     obs,
		Cache:             prodCache,
	})

	// Start gRPC Server
	productHandler := gapi.NewProductHandleGrpc(productService)
	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, productHandler)
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
	s.client = pb.NewProductServiceClient(conn)

	// Setup prerequisites directly via DB
	ctx := context.Background()
	
	// Create User & Merchant
	var userID int32
	err = pool.QueryRow(ctx, "INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING user_id",
		"Gapi", "User", "gapi@example.com", "password").Scan(&userID)
	s.Require().NoError(err)

	err = pool.QueryRow(ctx, "INSERT INTO merchants (user_id, name, description) VALUES ($1, $2, $3) RETURNING merchant_id",
		userID, "Gapi Merchant", "Gapi Description").Scan(&s.merchantID)
	s.Require().NoError(err)

	// Create Category
	err = pool.QueryRow(ctx, "INSERT INTO categories (name, description, slug_category) VALUES ($1, $2, $3) RETURNING category_id",
		"Gapi Category", "Gapi Category Description", "gapi-category").Scan(&s.categoryID)
	s.Require().NoError(err)
}

func (s *ProductGapiTestSuite) TearDownSuite() {
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

func (s *ProductGapiTestSuite) TestProductLifecycle() {
	ctx := context.Background()
	timestamp := time.Now().UnixNano()
	slug := fmt.Sprintf("gapi-product-%d", timestamp)

	// 1. Create
	createReq := &pb.CreateProductRequest{
		MerchantId:   s.merchantID,
		CategoryId:   s.categoryID,
		Name:         "Gapi Product",
		Description:  "A product created via gRPC",
		Price:        300,
		CountInStock: 20,
		Brand:        "GapiBrand",
		Weight:       150,
		Rating:       5,
		SlugProduct:  slug,
		ImageProduct: "gapi_prod.jpg",
	}

	createRes, err := s.client.Create(ctx, createReq)
	s.NoError(err)
	s.Equal("success", createRes.Status)
	s.Equal(createReq.Name, createRes.Data.Name)
	productID := createRes.Data.Id
	fmt.Printf("CREATED PRODUCT ID: %d\n", productID)

	// 2. Find By ID
	found, err := s.client.FindById(ctx, &pb.FindByIdProductRequest{Id: productID})
	s.NoError(err)
	s.Equal(productID, found.Data.Id)
	s.Equal(createReq.Name, found.Data.Name)

	// 3. Find All
	all, err := s.client.FindAll(ctx, &pb.FindAllProductRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.GreaterOrEqual(all.GetPagination().GetTotalRecords(), int32(1))

	// 4. Find By Category
	byCat, err := s.client.FindByCategory(ctx, &pb.FindAllProductCategoryRequest{
		CategoryName: "Gapi Category",
		Page:         1,
		PageSize:     10,
	})
	s.NoError(err)
	s.GreaterOrEqual(byCat.GetPagination().GetTotalRecords(), int32(1))

	// 5. Update
	updateReq := &pb.UpdateProductRequest{
		ProductId:    productID,
		MerchantId:   s.merchantID,
		CategoryId:   s.categoryID,
		Name:         "Gapi Product Updated",
		Description:  "Updated via gRPC",
		Price:        350,
		CountInStock: 15,
		Brand:        "GapiBrand v2",
		Weight:       160,
		Rating:       4,
		SlugProduct:  slug + "-updated",
		ImageProduct: "gapi_prod_updated.jpg",
	}
	updateRes, err := s.client.Update(ctx, updateReq)
	s.NoError(err)
	s.Equal(updateReq.Name, updateRes.Data.Name)

	// 6. Trash
	trashRes, err := s.client.TrashedProduct(ctx, &pb.FindByIdProductRequest{Id: productID})
	s.NoError(err)
	s.NotNil(trashRes.Data.DeletedAt)

	// 7. Find By Trashed
	trashed, err := s.client.FindByTrashed(ctx, &pb.FindAllProductRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.GreaterOrEqual(trashed.GetPagination().GetTotalRecords(), int32(1))

	// 8. Restore
	restoreRes, err := s.client.RestoreProduct(ctx, &pb.FindByIdProductRequest{Id: productID})
	s.NoError(err)
	s.Nil(restoreRes.Data.DeletedAt)

	// 9. Delete Permanent
	// Trash again before permanent delete
	_, err = s.client.TrashedProduct(ctx, &pb.FindByIdProductRequest{Id: productID})
	s.NoError(err)

	delRes, err := s.client.DeleteProductPermanent(ctx, &pb.FindByIdProductRequest{Id: productID})
	s.NoError(err)
	s.Equal("success", delRes.Status)

	// 10. Verify deletion
	_, err = s.client.FindById(ctx, &pb.FindByIdProductRequest{Id: productID})
	s.Error(err)
}

func TestProductGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductGapiTestSuite))
}
