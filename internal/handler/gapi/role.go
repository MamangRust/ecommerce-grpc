package gapi

import (
	"context"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	protomapper "ecommerce/internal/mapper/proto"
	"ecommerce/internal/pb"
	"ecommerce/internal/service"
	"ecommerce/pkg/errors/role_errors"
	"math"

	"google.golang.org/protobuf/types/known/emptypb"
)

type roleHandleGrpc struct {
	pb.UnimplementedRoleServiceServer
	roleService service.RoleService
	mapping     protomapper.RoleProtoMapper
}

func NewRoleHandleGrpc(role service.RoleService, mapping protomapper.RoleProtoMapper) *roleHandleGrpc {
	return &roleHandleGrpc{
		roleService: role,
		mapping:     mapping,
	}
}

func (s *roleHandleGrpc) FindAllRole(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRole, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	role, totalRecords, err := s.roleService.FindAll(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationRole(paginationMeta, "success", "Successfully fetched role records", role)

	return so, nil
}

func (s *roleHandleGrpc) FindByIdRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	id := int(req.GetRoleId())

	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	roleResponse := s.mapping.ToProtoResponseRole("success", "Successfully fetched role", role)

	return roleResponse, nil
}

func (s *roleHandleGrpc) FindByUserId(ctx context.Context, req *pb.FindByIdUserRoleRequest) (*pb.ApiResponsesRole, error) {
	id := int(req.GetUserId())

	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleService.FindByUserId(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	roleResponse := s.mapping.ToProtoResponsesRole("success", "Successfully fetched role by user ID", role)

	return roleResponse, nil
}

func (s *roleHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleService.FindByActiveRole(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched active roles", roles)

	return so, nil
}

func (s *roleHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllRoleRequest) (*pb.ApiResponsePaginationRoleDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllRole{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleService.FindByTrashedRole(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationRoleDeleteAt(paginationMeta, "success", "Successfully fetched trashed roles", roles)

	return so, nil
}

func (s *roleHandleGrpc) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.ApiResponseRole, error) {
	name := req.GetName()

	request := &requests.CreateRoleRequest{
		Name: name,
	}

	if err := request.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateCreateRole
	}

	role, err := s.roleService.CreateRole(&requests.CreateRoleRequest{
		Name: name,
	})

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully created role", role)

	return so, nil
}

func (s *roleHandleGrpc) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.ApiResponseRole, error) {
	id := int(req.GetId())

	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	name := req.GetName()

	request := &requests.UpdateRoleRequest{
		ID:   &id,
		Name: name,
	}

	if err := request.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateUpdateRole
	}

	role, err := s.roleService.UpdateRole(request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully updated role", role)

	return so, nil
}

func (s *roleHandleGrpc) TrashedRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	id := int(req.GetRoleId())

	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleService.TrashedRole(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully trashed role", role)

	return so, nil
}

func (s *roleHandleGrpc) RestoreRole(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRole, error) {
	id := int(req.GetRoleId())

	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleService.RestoreRole(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRole("success", "Successfully restored role", role)

	return so, nil
}

func (s *roleHandleGrpc) DeleteRolePermanent(ctx context.Context, req *pb.FindByIdRoleRequest) (*pb.ApiResponseRoleDelete, error) {
	id := int(req.GetRoleId())

	if id == 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	_, err := s.roleService.DeleteRolePermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRoleDelete("success", "Successfully deleted role permanently")

	return so, nil
}

func (s *roleHandleGrpc) RestoreAllRole(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	_, err := s.roleService.RestoreAllRole()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRoleAll("success", "Successfully restored all roles")

	return so, nil
}

func (s *roleHandleGrpc) DeleteAllRolePermanent(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseRoleAll, error) {
	_, err := s.roleService.DeleteAllRolePermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseRoleAll("success", "Successfully deleted all roles")

	return so, nil
}
