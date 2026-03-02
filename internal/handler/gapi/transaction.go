package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors"
	"ecommerce/pkg/errors/transaction_errors"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type transactionHandleGrpc struct {
	pb.UnimplementedTransactionServiceServer
	transactionService service.TransactionService
}

func NewTransactionHandleGrpc(
	transactionService service.TransactionService,
) *transactionHandleGrpc {
	return &transactionHandleGrpc{
		transactionService: transactionService,
	}
}

func (s *transactionHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransaction, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := s.transactionService.FindAllTransactions(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTransactions := make([]*pb.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		paymentStatus := ""
		if transaction.PaymentStatus != "" {
			paymentStatus = transaction.PaymentStatus
		}

		protoTransactions[i] = &pb.TransactionResponse{
			Id:            int32(transaction.TransactionID),
			OrderId:       int32(transaction.OrderID),
			MerchantId:    int32(transaction.MerchantID),
			PaymentMethod: transaction.PaymentMethod,
			Amount:        int32(transaction.Amount),
			PaymentStatus: paymentStatus,
			CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transaction records",
		Data:       protoTransactions,
		Pagination: paginationMeta,
	}, nil
}

func (s *transactionHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindAllTransactionMerchantRequest) (*pb.ApiResponsePaginationTransaction, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchant_id := int(request.GetMerchantId())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransactionByMerchant{
		MerchantID: merchant_id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	transactions, totalRecords, err := s.transactionService.FindByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTransactions := make([]*pb.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		paymentStatus := ""
		if transaction.PaymentStatus != "" {
			paymentStatus = transaction.PaymentStatus
		}

		protoTransactions[i] = &pb.TransactionResponse{
			Id:            int32(transaction.TransactionID),
			OrderId:       int32(transaction.OrderID),
			MerchantId:    int32(transaction.MerchantID),
			PaymentMethod: transaction.PaymentMethod,
			Amount:        int32(transaction.Amount),
			PaymentStatus: paymentStatus,
			CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched merchant transaction records",
		Data:       protoTransactions,
		Pagination: paginationMeta,
	}, nil
}

func (s *transactionHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionService.FindById(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	paymentStatus := ""
	if transaction.PaymentStatus != "" {
		paymentStatus = transaction.PaymentStatus
	}

	protoTransaction := &pb.TransactionResponse{
		Id:            int32(transaction.TransactionID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		PaymentStatus: paymentStatus,
		CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction",
		Data:    protoTransaction,
	}, nil
}

func (s *transactionHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := s.transactionService.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTransactions := make([]*pb.TransactionResponseDeleteAt, len(transactions))
	for i, transaction := range transactions {
		paymentStatus := ""
		if transaction.PaymentStatus != "" {
			paymentStatus = transaction.PaymentStatus
		}

		protoTransactions[i] = &pb.TransactionResponseDeleteAt{
			Id:            int32(transaction.TransactionID),
			OrderId:       int32(transaction.OrderID),
			MerchantId:    int32(transaction.MerchantID),
			PaymentMethod: transaction.PaymentMethod,
			Amount:        int32(transaction.Amount),
			PaymentStatus: paymentStatus,
			CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active transaction records",
		Data:       protoTransactions,
		Pagination: paginationMeta,
	}, nil
}

func (s *transactionHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := s.transactionService.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoTransactions := make([]*pb.TransactionResponseDeleteAt, len(transactions))
	for i, transaction := range transactions {
		paymentStatus := ""
		if transaction.PaymentStatus != "" {
			paymentStatus = transaction.PaymentStatus
		}

		var deletedAt string
		if transaction.DeletedAt.Valid {
			deletedAt = transaction.DeletedAt.Time.Format("2006-01-02")
		}

		protoTransactions[i] = &pb.TransactionResponseDeleteAt{
			Id:            int32(transaction.TransactionID),
			OrderId:       int32(transaction.OrderID),
			MerchantId:    int32(transaction.MerchantID),
			PaymentMethod: transaction.PaymentMethod,
			Amount:        int32(transaction.Amount),
			PaymentStatus: paymentStatus,
			CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed transaction records",
		Data:       protoTransactions,
		Pagination: paginationMeta,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthStatusSuccess(ctx context.Context, request *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionService.FindMonthlyAmountSuccess(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionMonthlyAmountSuccess, len(res))
	for i, item := range res {
		protoData[i] = &pb.TransactionMonthlyAmountSuccess{
			Year:         item.Year,
			Month:        item.Month,
			TotalSuccess: int32(item.TotalSuccess),
			TotalAmount:  int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Monthly success data retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindYearStatusSuccess(ctx context.Context, request *pb.FindYearlyTransactionStatus) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	year := int(request.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	res, err := s.transactionService.FindYearlyAmountSuccess(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionYearlyAmountSuccess, len(res))
	for i, item := range res {
		protoData[i] = &pb.TransactionYearlyAmountSuccess{
			Year:         item.Year,
			TotalSuccess: int32(item.TotalSuccess),
			TotalAmount:  int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Yearly success data retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthStatusFailed(ctx context.Context, request *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionService.FindMonthlyAmountFailed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionMonthlyAmountFailed, len(res))
	for i, item := range res {
		protoData[i] = &pb.TransactionMonthlyAmountFailed{
			Year:        item.Year,
			Month:       item.Month,
			TotalFailed: int32(item.TotalFailed),
			TotalAmount: int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Monthly failed data retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindYearStatusFailed(ctx context.Context, request *pb.FindYearlyTransactionStatus) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	year := int(request.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	res, err := s.transactionService.FindYearlyAmountFailed(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionYearlyAmountFailed, len(res))
	for i, item := range res {
		protoData[i] = &pb.TransactionYearlyAmountFailed{
			Year:        item.Year,
			TotalFailed: int32(item.TotalFailed),
			TotalAmount: int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Yearly failed data retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthStatusSuccessByMerchant(ctx context.Context, request *pb.FindMonthlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.MonthAmountTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	res, err := s.transactionService.FindMonthlyAmountSuccessByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionMonthlyAmountSuccess, len(res))
	for i, item := range res {
		protoData[i] = &pb.TransactionMonthlyAmountSuccess{
			Year:         item.Year,
			Month:        item.Month,
			TotalSuccess: int32(item.TotalSuccess),
			TotalAmount:  int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Merchant monthly success data retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindYearStatusSuccessByMerchant(ctx context.Context, request *pb.FindYearlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearAmountTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.transactionService.FindYearlyAmountSuccessByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionYearlyAmountSuccess, len(res))
	for i, item := range res {
		protoData[i] = &pb.TransactionYearlyAmountSuccess{
			Year:         item.Year,
			TotalSuccess: int32(item.TotalSuccess),
			TotalAmount:  int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Merchant yearly success data retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthStatusFailedByMerchant(ctx context.Context, request *pb.FindMonthlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.MonthAmountTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	res, err := s.transactionService.FindMonthlyAmountFailedByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionMonthlyAmountFailed, len(res))
	for i, item := range res {
		protoData[i] = &pb.TransactionMonthlyAmountFailed{
			Year:        item.Year,
			Month:       item.Month,
			TotalFailed: int32(item.TotalFailed),
			TotalAmount: int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Merchant monthly failed data retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindYearStatusFailedByMerchant(ctx context.Context, request *pb.FindYearlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearAmountTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.transactionService.FindYearlyAmountFailedByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionYearlyAmountFailed, len(res))
	for i, item := range res {
		protoData[i] = &pb.TransactionYearlyAmountFailed{
			Year:        item.Year,
			TotalFailed: int32(item.TotalFailed),
			TotalAmount: int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Merchant yearly failed data retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthMethodSuccess(ctx context.Context, req *pb.MonthTransactionMethod) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	methods, err := s.transactionService.FindMonthlyTransactionMethodSuccess(ctx, &requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionMonthlyMethod, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TransactionMonthlyMethod{
			Month:             item.Month,
			PaymentMethod:     item.PaymentMethod,
			TotalTransactions: int32(item.TotalTransactions),
			TotalAmount:       int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindYearMethodSuccess(ctx context.Context, req *pb.YearTransactionMethod) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := s.transactionService.FindYearlyTransactionMethodSuccess(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionYearlyMethod, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TransactionYearlyMethod{
			Year:              item.Year,
			PaymentMethod:     item.PaymentMethod,
			TotalTransactions: int32(item.TotalTransactions),
			TotalAmount:       int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthMethodByMerchantSuccess(ctx context.Context, req *pb.MonthTransactionMethodByMerchant) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	if month <= 0 || month > 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
		Month:      month,
	}

	methods, err := s.transactionService.FindMonthlyTransactionMethodByMerchantSuccess(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionMonthlyMethod, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TransactionMonthlyMethod{
			Month:             item.Month,
			PaymentMethod:     item.PaymentMethod,
			TotalTransactions: int32(item.TotalTransactions),
			TotalAmount:       int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindYearMethodByMerchantSuccess(ctx context.Context, req *pb.YearTransactionMethodByMerchant) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.transactionService.FindYearlyTransactionMethodByMerchantSuccess(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionYearlyMethod, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TransactionYearlyMethod{
			Year:              item.Year,
			PaymentMethod:     item.PaymentMethod,
			TotalTransactions: int32(item.TotalTransactions),
			TotalAmount:       int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthMethodFailed(ctx context.Context, req *pb.MonthTransactionMethod) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	methods, err := s.transactionService.FindMonthlyTransactionMethodFailed(ctx, &requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionMonthlyMethod, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TransactionMonthlyMethod{
			Month:             item.Month,
			PaymentMethod:     item.PaymentMethod,
			TotalTransactions: int32(item.TotalTransactions),
			TotalAmount:       int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindYearMethodFailed(ctx context.Context, req *pb.YearTransactionMethod) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := s.transactionService.FindYearlyTransactionMethodFailed(ctx, year)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionYearlyMethod, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TransactionYearlyMethod{
			Year:              item.Year,
			PaymentMethod:     item.PaymentMethod,
			TotalTransactions: int32(item.TotalTransactions),
			TotalAmount:       int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindMonthMethodByMerchantFailed(ctx context.Context, req *pb.MonthTransactionMethodByMerchant) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	if month <= 0 || month > 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
		Month:      month,
	}

	methods, err := s.transactionService.FindMonthlyTransactionMethodByMerchantFailed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionMonthlyMethod, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TransactionMonthlyMethod{
			Month:             item.Month,
			PaymentMethod:     item.PaymentMethod,
			TotalTransactions: int32(item.TotalTransactions),
			TotalAmount:       int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) FindYearMethodByMerchantFailed(ctx context.Context, req *pb.YearTransactionMethodByMerchant) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.transactionService.FindYearlyTransactionMethodByMerchantFailed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoData := make([]*pb.TransactionYearlyMethod, len(methods))
	for i, item := range methods {
		protoData[i] = &pb.TransactionYearlyMethod{
			Year:              item.Year,
			PaymentMethod:     item.PaymentMethod,
			TotalTransactions: int32(item.TotalTransactions),
			TotalAmount:       int32(item.TotalAmount),
		}
	}

	return &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    protoData,
	}, nil
}

func (s *transactionHandleGrpc) Create(ctx context.Context, request *pb.CreateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	req := &requests.CreateTransactionRequest{
		OrderID:       int(request.GetOrderId()),
		MerchantID:    int(request.GetMerchantId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateCreateTransaction
	}

	transaction, err := s.transactionService.CreateTransaction(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	paymentStatus := ""
	if transaction.PaymentStatus == "" {
		paymentStatus = transaction.PaymentStatus
	}

	protoTransaction := &pb.TransactionResponse{
		Id:            int32(transaction.TransactionID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		PaymentStatus: paymentStatus,
		CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully created transaction",
		Data:    protoTransaction,
	}, nil
}

func (s *transactionHandleGrpc) Update(ctx context.Context, request *pb.UpdateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateTransactionRequest{
		TransactionID: &id,
		OrderID:       int(request.GetOrderId()),
		MerchantID:    int(request.GetMerchantId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateUpdateTransaction
	}

	transaction, err := s.transactionService.UpdateTransaction(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	paymentStatus := ""
	if transaction.PaymentStatus == "" {
		paymentStatus = transaction.PaymentStatus
	}

	protoTransaction := &pb.TransactionResponse{
		Id:            int32(transaction.TransactionID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		PaymentStatus: paymentStatus,
		CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully updated transaction",
		Data:    protoTransaction,
	}, nil
}

func (s *transactionHandleGrpc) TrashedTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionService.TrashedTransaction(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	paymentStatus := ""
	if transaction.PaymentStatus == "" {
		paymentStatus = transaction.PaymentStatus
	}

	var deletedAt string
	if transaction.DeletedAt.Valid {
		deletedAt = transaction.DeletedAt.Time.Format("2006-01-02")
	}

	protoTransaction := &pb.TransactionResponseDeleteAt{
		Id:            int32(transaction.TransactionID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		PaymentStatus: paymentStatus,
		CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transaction",
		Data:    protoTransaction,
	}, nil
}

func (s *transactionHandleGrpc) RestoreTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionService.RestoreTransaction(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	paymentStatus := ""
	if transaction.PaymentStatus == "" {
		paymentStatus = transaction.PaymentStatus
	}

	var deletedAt string
	if transaction.DeletedAt.Valid {
		deletedAt = transaction.DeletedAt.Time.Format("2006-01-02")
	}

	protoTransaction := &pb.TransactionResponseDeleteAt{
		Id:            int32(transaction.TransactionID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		PaymentStatus: paymentStatus,
		CreatedAt:     transaction.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:     transaction.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:     &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully restored transaction",
		Data:    protoTransaction,
	}, nil
}

func (s *transactionHandleGrpc) DeleteTransactionPermanent(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	_, err := s.transactionService.DeleteTransactionPermanently(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted transaction permanently",
	}, nil
}

func (s *transactionHandleGrpc) RestoreAllTransaction(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	_, err := s.transactionService.RestoreAllTransactions(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully restored all transactions",
	}, nil
}

func (s *transactionHandleGrpc) DeleteAllTransactionPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	_, err := s.transactionService.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully deleted all transactions permanently",
	}, nil
}
