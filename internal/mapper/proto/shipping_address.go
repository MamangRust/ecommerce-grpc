package protomapper

import (
	"ecommerce/internal/domain/response"
	"ecommerce/internal/pb"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type shippingAddressProtoMapper struct {
}

func NewShippingAddressProtoMapper() *shippingAddressProtoMapper {
	return &shippingAddressProtoMapper{}
}

func (s *shippingAddressProtoMapper) ToProtoResponseShippingAddress(status string, message string, pbResponse *response.ShippingAddressResponse) *pb.ApiResponseShipping {
	return &pb.ApiResponseShipping{
		Status:  status,
		Message: message,
		Data:    s.mapResponseShippingAddress(pbResponse),
	}
}

func (s *shippingAddressProtoMapper) ToProtoResponseShippingAddressDeleteAt(status string, message string, pbResponse *response.ShippingAddressResponseDeleteAt) *pb.ApiResponseShippingDeleteAt {
	return &pb.ApiResponseShippingDeleteAt{
		Status:  status,
		Message: message,
		Data:    s.mapResponseShippingAddressDeleteAt(pbResponse),
	}
}

func (s *shippingAddressProtoMapper) ToProtoResponsesShippingAddress(status string, message string, pbResponse []*response.ShippingAddressResponse) *pb.ApiResponsesShipping {
	return &pb.ApiResponsesShipping{
		Status:  status,
		Message: message,
		Data:    s.mapResponsesShippingAddress(pbResponse),
	}
}

func (s *shippingAddressProtoMapper) ToProtoResponseShippingAddressDelete(status string, message string) *pb.ApiResponseShippingDelete {
	return &pb.ApiResponseShippingDelete{
		Status:  status,
		Message: message,
	}
}

func (c *shippingAddressProtoMapper) ToProtoResponseShippingAddressAll(status string, message string) *pb.ApiResponseShippingAll {
	return &pb.ApiResponseShippingAll{
		Status:  status,
		Message: message,
	}
}

func (s *shippingAddressProtoMapper) ToProtoResponsePaginationShippingAddressDeleteAt(pagination *pb.PaginationMeta, status string, message string, addresses []*response.ShippingAddressResponseDeleteAt) *pb.ApiResponsePaginationShippingDeleteAt {
	return &pb.ApiResponsePaginationShippingDeleteAt{
		Status:     status,
		Message:    message,
		Data:       s.mapResponsesShippingAddressDeleteAt(addresses),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *shippingAddressProtoMapper) ToProtoResponsePaginationShippingAddress(pagination *pb.PaginationMeta, status string, message string, addresses []*response.ShippingAddressResponse) *pb.ApiResponsePaginationShipping {
	return &pb.ApiResponsePaginationShipping{
		Status:     status,
		Message:    message,
		Data:       s.mapResponsesShippingAddress(addresses),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (s *shippingAddressProtoMapper) mapResponseShippingAddress(address *response.ShippingAddressResponse) *pb.ShippingResponse {
	return &pb.ShippingResponse{
		Id:        int32(address.ID),
		OrderId:   int32(address.OrderID),
		Alamat:    address.Alamat,
		Provinsi:  address.Provinsi,
		Negara:    address.Negara,
		Kota:      address.Kota,
		CreatedAt: address.CreatedAt,
		UpdatedAt: address.UpdatedAt,
	}
}

func (s *shippingAddressProtoMapper) mapResponsesShippingAddress(addresses []*response.ShippingAddressResponse) []*pb.ShippingResponse {
	var mappedAddresses []*pb.ShippingResponse

	for _, address := range addresses {
		mappedAddresses = append(mappedAddresses, s.mapResponseShippingAddress(address))
	}

	return mappedAddresses
}

func (s *shippingAddressProtoMapper) mapResponseShippingAddressDeleteAt(address *response.ShippingAddressResponseDeleteAt) *pb.ShippingResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue

	if address.DeletedAt != nil {
		deletedAt = wrapperspb.String(*address.DeletedAt)
	}

	return &pb.ShippingResponseDeleteAt{
		Id:        int32(address.ID),
		OrderId:   int32(address.OrderID),
		Alamat:    address.Alamat,
		Provinsi:  address.Provinsi,
		Negara:    address.Negara,
		Kota:      address.Kota,
		CreatedAt: address.CreatedAt,
		UpdatedAt: address.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (s *shippingAddressProtoMapper) mapResponsesShippingAddressDeleteAt(addresses []*response.ShippingAddressResponseDeleteAt) []*pb.ShippingResponseDeleteAt {
	var mappedAddresses []*pb.ShippingResponseDeleteAt

	for _, address := range addresses {
		mappedAddresses = append(mappedAddresses, s.mapResponseShippingAddressDeleteAt(address))
	}

	return mappedAddresses
}
