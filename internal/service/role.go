package service

import (
	"database/sql"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type roleService struct {
	roleRepository repository.RoleRepository
	logger         logger.LoggerInterface
	mapping        response_service.RoleResponseMapper
}

func NewRoleService(roleRepository repository.RoleRepository, logger logger.LoggerInterface, mapping response_service.RoleResponseMapper) *roleService {
	return &roleService{
		roleRepository: roleRepository,
		logger:         logger,
		mapping:        mapping,
	}
}

func (s *roleService) FindAll(req *requests.FindAllRole) ([]*response.RoleResponse, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching role",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.roleRepository.FindAllRoles(req)

	if err != nil {
		s.logger.Error("Failed to retrieve role list from database",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve role list",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched role",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	so := s.mapping.ToRolesResponse(res)

	return so, totalRecords, nil
}

func (s *roleService) FindById(id int) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching role by ID", zap.Int("id", id))

	res, err := s.roleRepository.FindById(id)

	if err != nil {
		s.logger.Error("Failed to retrieve role details",
			zap.Int("role_id", id),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve role details",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched role", zap.Int("id", id))

	so := s.mapping.ToRoleResponse(res)

	return so, nil
}

func (s *roleService) FindByUserId(id int) ([]*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching role by user ID", zap.Int("id", id))

	res, err := s.roleRepository.FindByUserId(id)

	if err != nil {
		s.logger.Error("Failed to retrieve role by user ID",
			zap.Int("user_id", id),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve user role",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched role by user ID", zap.Int("id", id))

	so := s.mapping.ToRolesResponse(res)

	return so, nil
}

func (s *roleService) FindByActiveRole(req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching active role",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.roleRepository.FindByActiveRole(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active roles",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))

		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve active roles",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched active role",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	so := s.mapping.ToRolesResponseDeleteAt(res)

	return so, totalRecords, nil
}

func (s *roleService) FindByTrashedRole(req *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching trashed role",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	res, totalRecords, err := s.roleRepository.FindByTrashedRole(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed roles from database",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search))
		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve trashed roles",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully fetched trashed role",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	so := s.mapping.ToRolesResponseDeleteAt(res)

	return so, totalRecords, nil
}

func (s *roleService) CreateRole(request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting CreateRole process",
		zap.String("roleName", request.Name),
	)

	res, err := s.roleRepository.CreateRole(request)

	if err != nil {
		s.logger.Error("Failed to create new role record",
			zap.String("role_name", request.Name),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create new role",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToRoleResponse(res)

	s.logger.Debug("CreateRole process completed",
		zap.String("roleName", request.Name),
		zap.Int("roleID", res.ID),
	)

	return so, nil
}

func (s *roleService) UpdateRole(request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting UpdateRole process",
		zap.Int("roleID", *request.ID),
		zap.String("newRoleName", request.Name),
	)

	res, err := s.roleRepository.UpdateRole(request)

	if err != nil {
		s.logger.Error("Failed to update role record",
			zap.Int("role_id", *request.ID),
			zap.String("new_name", request.Name),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update role",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToRoleResponse(res)

	s.logger.Debug("UpdateRole process completed",
		zap.Int("roleID", *request.ID),
		zap.String("newRoleName", request.Name),
	)

	return so, nil
}

func (s *roleService) TrashedRole(id int) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting TrashedRole process",
		zap.Int("roleID", id),
	)

	res, err := s.roleRepository.TrashedRole(id)

	if err != nil {
		s.logger.Error("Failed to move role to trash",
			zap.Int("role_id", id),
			zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("Role with ID %d not found", id),
				Code:    http.StatusNotFound,
			}
		}

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to move role to trash",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToRoleResponse(res)

	s.logger.Debug("TrashedRole process completed",
		zap.Int("roleID", id),
	)

	return so, nil
}

func (s *roleService) RestoreRole(id int) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting RestoreRole process",
		zap.Int("roleID", id),
	)

	res, err := s.roleRepository.RestoreRole(id)

	if err != nil {
		s.logger.Error("Failed to restore role from trash",
			zap.Int("role_id", id),
			zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("Role with ID %d not found in trash", id),
				Code:    http.StatusNotFound,
			}
		}
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore role from trash",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToRoleResponse(res)

	s.logger.Debug("RestoreRole process completed",
		zap.Int("roleID", id),
	)

	return so, nil
}

func (s *roleService) DeleteRolePermanent(id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Starting DeleteRolePermanent process",
		zap.Int("roleID", id),
	)

	_, err := s.roleRepository.DeleteRolePermanent(id)

	if err != nil {
		s.logger.Error("Failed to permanently delete role",
			zap.Int("role_id", id),
			zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			return false, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("Role with ID %d not found", id),
				Code:    http.StatusNotFound,
			}
		}
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete role",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("DeleteRolePermanent process completed",
		zap.Int("roleID", id),
	)

	return true, nil
}

func (s *roleService) RestoreAllRole() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all roles")

	_, err := s.roleRepository.RestoreAllRole()

	if err != nil {
		s.logger.Error("Failed to restore all trashed roles",
			zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all roles",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully restored all roles")
	return true, nil
}

func (s *roleService) DeleteAllRolePermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all roles")

	_, err := s.roleRepository.DeleteAllRolePermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed roles",
			zap.Error(err))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all roles",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully deleted all roles permanently")
	return true, nil
}
