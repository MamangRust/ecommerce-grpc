package api

import (
	response_api "ecommerce/internal/mapper/response/api"
	"ecommerce/internal/pb"
	"ecommerce/pkg/auth"
	"ecommerce/pkg/logger"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type Deps struct {
	Conn    *grpc.ClientConn
	Token   auth.TokenManager
	E       *echo.Echo
	Logger  logger.LoggerInterface
	Mapping response_api.ResponseApiMapper
}

func NewHandler(deps Deps) {

	clientAuth := pb.NewAuthServiceClient(deps.Conn)
	clientRole := pb.NewRoleServiceClient(deps.Conn)
	clientUser := pb.NewUserServiceClient(deps.Conn)
	clientCategory := pb.NewCategoryServiceClient(deps.Conn)
	clientMerchant := pb.NewMerchantServiceClient(deps.Conn)
	clientOrderItem := pb.NewOrderItemServiceClient(deps.Conn)
	clientOrder := pb.NewOrderServiceClient(deps.Conn)
	clientProduct := pb.NewProductServiceClient(deps.Conn)
	clientTransaction := pb.NewTransactionServiceClient(deps.Conn)
	clientCart := pb.NewCartServiceClient(deps.Conn)
	clientReview := pb.NewReviewServiceClient(deps.Conn)
	clientSlider := pb.NewSliderServiceClient(deps.Conn)
	clientShipping := pb.NewShippingServiceClient(deps.Conn)

	NewHandlerAuth(deps.E, clientAuth, deps.Logger, deps.Mapping.AuthResponseMapper)
	NewHandlerRole(deps.E, clientRole, deps.Logger, deps.Mapping.RoleResponseMapper)
	NewHandlerUser(deps.E, clientUser, deps.Logger, deps.Mapping.UserResponseMapper)
	NewHandlerCategory(deps.E, clientCategory, deps.Logger, deps.Mapping.CategoryResponseMapper)
	NewHandlerMerchant(deps.E, clientMerchant, deps.Logger, deps.Mapping.MerchantResponseMapper)
	NewHandlerOrderItem(deps.E, clientOrderItem, deps.Logger, deps.Mapping.OrderItemResponseMapper)
	NewHandlerOrder(deps.E, clientOrder, deps.Logger, deps.Mapping.OrderResponseMapper)
	NewHandlerProduct(deps.E, clientProduct, deps.Logger, deps.Mapping.ProductResponseMapper)
	NewHandlerTransaction(deps.E, clientTransaction, deps.Logger, deps.Mapping.TransactionResponseMapper)
	NewHandlerCart(deps.E, clientCart, deps.Logger, deps.Mapping.CartResponseMapper)
	NewHandlerReview(deps.E, clientReview, deps.Logger, deps.Mapping.ReviewMapper)
	NewHandlerSlider(deps.E, clientSlider, deps.Logger, deps.Mapping.SliderMapper)
	NewHandlerShippingAddress(deps.E, clientShipping, deps.Logger, deps.Mapping.ShippingAddressResponseMapper)
}
