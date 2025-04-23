package response_api

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"
)

type shippingAddressResponseMapper struct {
}

func NewShippingAddressResponseMapper() *shippingAddressResponseMapper {
	return &shippingAddressResponseMapper{}
}

func (s *shippingAddressResponseMapper) ToResponseShippingAddress(pbResponse *pb.ShippingResponse) *response.ShippingAddressResponse {
	return &response.ShippingAddressResponse{
		ID:        int(pbResponse.Id),
		OrderID:   int(pbResponse.OrderId),
		Alamat:    pbResponse.Alamat,
		Provinsi:  pbResponse.Provinsi,
		Negara:    pbResponse.Negara,
		Kota:      pbResponse.Kota,
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
	}
}

func (s *shippingAddressResponseMapper) ToResponsesShippingAddress(pbResponses []*pb.ShippingResponse) []*response.ShippingAddressResponse {
	var addresses []*response.ShippingAddressResponse
	for _, address := range pbResponses {
		addresses = append(addresses, s.ToResponseShippingAddress(address))
	}
	return addresses
}

func (s *shippingAddressResponseMapper) ToResponseShippingAddressDeleteAt(pbResponse *pb.ShippingResponseDeleteAt) *response.ShippingAddressResponseDeleteAt {
	var deletedAt string
	if pbResponse.DeletedAt != nil {
		deletedAt = pbResponse.DeletedAt.Value
	}

	return &response.ShippingAddressResponseDeleteAt{
		ID:        int(pbResponse.Id),
		OrderID:   int(pbResponse.OrderId),
		Alamat:    pbResponse.Alamat,
		Provinsi:  pbResponse.Provinsi,
		Negara:    pbResponse.Negara,
		Kota:      pbResponse.Kota,
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (s *shippingAddressResponseMapper) ToResponsesShippingAddressDeleteAt(pbResponses []*pb.ShippingResponseDeleteAt) []*response.ShippingAddressResponseDeleteAt {
	var addresses []*response.ShippingAddressResponseDeleteAt
	for _, address := range pbResponses {
		addresses = append(addresses, s.ToResponseShippingAddressDeleteAt(address))
	}
	return addresses
}

func (s *shippingAddressResponseMapper) ToApiResponseShippingAddress(pbResponse *pb.ApiResponseShipping) *response.ApiResponseShippingAddress {
	return &response.ApiResponseShippingAddress{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.ToResponseShippingAddress(pbResponse.Data),
	}
}

func (s *shippingAddressResponseMapper) ToApiResponseShippingAddressDeleteAt(pbResponse *pb.ApiResponseShippingDeleteAt) *response.ApiResponseShippingAddressDeleteAt {
	return &response.ApiResponseShippingAddressDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.ToResponseShippingAddressDeleteAt(pbResponse.Data),
	}
}

func (s *shippingAddressResponseMapper) ToApiResponsesShippingAddress(pbResponse *pb.ApiResponsesShipping) *response.ApiResponsesShippingAddress {
	return &response.ApiResponsesShippingAddress{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.ToResponsesShippingAddress(pbResponse.Data),
	}
}

func (s *shippingAddressResponseMapper) ToApiResponseShippingAddressDelete(pbResponse *pb.ApiResponseShippingDelete) *response.ApiResponseShippingAddressDelete {
	return &response.ApiResponseShippingAddressDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *shippingAddressResponseMapper) ToApiResponseShippingAddressAll(pbResponse *pb.ApiResponseShippingAll) *response.ApiResponseShippingAddressAll {
	return &response.ApiResponseShippingAddressAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *shippingAddressResponseMapper) ToApiResponsePaginationShippingAddressDeleteAt(pbResponse *pb.ApiResponsePaginationShippingDeleteAt) *response.ApiResponsePaginationShippingAddressDeleteAt {
	return &response.ApiResponsePaginationShippingAddressDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.ToResponsesShippingAddressDeleteAt(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *shippingAddressResponseMapper) ToApiResponsePaginationShippingAddress(pbResponse *pb.ApiResponsePaginationShipping) *response.ApiResponsePaginationShippingAddress {
	return &response.ApiResponsePaginationShippingAddress{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.ToResponsesShippingAddress(pbResponse.Data),
		Pagination: *mapPaginationMeta(pbResponse.Pagination),
	}
}
