package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/logger"

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

func (s *roleService) FindAll(page int, pageSize int, search string) ([]*response.RoleResponse, int, *response.ErrorResponse) {
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

	res, totalRecords, err := s.roleRepository.FindAllRoles(page, pageSize, search)
	if err != nil {
		s.logger.Error("Failed to fetch role",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch role records",
		}
	}

	s.logger.Debug("Successfully fetched role",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	so := s.mapping.ToRolesResponse(res)

	return so, totalRecords, nil
}

func (s *roleService) FindById(id int) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching role by ID", zap.Int("id", id))

	res, err := s.roleRepository.FindById(id)

	if err != nil {
		s.logger.Error("Failed to fetch role record by ID", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch role record by ID",
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
		s.logger.Error("Failed to fetch role record by ID", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch role record by ID",
		}
	}

	s.logger.Debug("Successfully fetched role by user ID", zap.Int("id", id))

	so := s.mapping.ToRolesResponse(res)

	return so, nil
}

func (s *roleService) FindByActiveRole(page int, pageSize int, search string) ([]*response.RoleResponseDeleteAt, int, *response.ErrorResponse) {
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

	res, totalRecords, err := s.roleRepository.FindByActiveRole(page, pageSize, search)
	if err != nil {
		s.logger.Error("Failed to fetch active role",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch role records",
		}
	}

	s.logger.Debug("Successfully fetched active role",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	so := s.mapping.ToRolesResponseDeleteAt(res)

	return so, totalRecords, nil
}

func (s *roleService) FindByTrashedRole(page int, pageSize int, search string) ([]*response.RoleResponseDeleteAt, int, *response.ErrorResponse) {
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

	res, totalRecords, err := s.roleRepository.FindByTrashedRole(page, pageSize, search)

	if err != nil {
		s.logger.Error("Failed to fetch trashed role",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch role records",
		}
	}

	s.logger.Debug("Successfully fetched trashed role",
		zap.Int("totalRecords", totalRecords),
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
		s.logger.Error("Failed to create role",
			zap.String("roleName", request.Name),
			zap.Error(err),
		)

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create role record",
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
		zap.Int("roleID", request.ID),
		zap.String("newRoleName", request.Name),
	)

	res, err := s.roleRepository.UpdateRole(request)

	if err != nil {
		s.logger.Error("Failed to update role",
			zap.Int("roleID", request.ID),
			zap.String("newRoleName", request.Name),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update role record",
		}
	}

	so := s.mapping.ToRoleResponse(res)

	s.logger.Debug("UpdateRole process completed",
		zap.Int("roleID", request.ID),
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
			zap.Int("roleID", id),
			zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trashed role record",
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
		s.logger.Error("Failed to restore role", zap.Error(err))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore role record",
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
		s.logger.Error("Failed to delete role permanently",
			zap.Int("roleID", id),
			zap.Error(err),
		)

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete role record",
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
		s.logger.Error("Failed to restore all roles", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all roles: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully restored all roles")
	return true, nil
}

func (s *roleService) DeleteAllRolePermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all roles")

	_, err := s.roleRepository.DeleteAllRolePermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all roles", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all roles: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully deleted all roles permanently")
	return true, nil
}
