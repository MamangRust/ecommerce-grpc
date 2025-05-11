package service

import (
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/errors/role_errors"
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

func (s *roleService) FindAll(request *requests.FindAllRole) ([]*response.RoleResponse, *int, *response.ErrorResponse) {
	page := request.Page
	pageSize := request.PageSize
	search := request.Search

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

	res, totalRecords, err := s.roleRepository.FindAllRoles(request)
	if err != nil {
		s.logger.Error("Failed to fetch role",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, role_errors.ErrFailedFindAll
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
		s.logger.Error("Failed to fetch role record by ID", zap.Error(err))

		return nil, role_errors.ErrRoleNotFoundRes
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

		return nil, role_errors.ErrRoleNotFoundRes
	}

	s.logger.Debug("Successfully fetched role by user ID", zap.Int("id", id))

	so := s.mapping.ToRolesResponse(res)

	return so, nil
}

func (s *roleService) FindByActive(request *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page := request.Page
	pageSize := request.PageSize
	search := request.Search

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

	res, totalRecords, err := s.roleRepository.FindByActiveRole(request)
	if err != nil {
		s.logger.Error("Failed to fetch active role",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, role_errors.ErrFailedFindActive
	}

	s.logger.Debug("Successfully fetched active role",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	so := s.mapping.ToRolesResponseDeleteAt(res)

	return so, totalRecords, nil
}

func (s *roleService) FindByTrashed(request *requests.FindAllRole) ([]*response.RoleResponseDeleteAt, *int, *response.ErrorResponse) {
	page := request.Page
	pageSize := request.PageSize
	search := request.Search

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

	res, totalRecords, err := s.roleRepository.FindByTrashedRole(request)

	if err != nil {
		s.logger.Error("Failed to fetch trashed role",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, nil, role_errors.ErrFailedFindTrashed
	}

	s.logger.Debug("Successfully fetched trashed role",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	so := s.mapping.ToRolesResponseDeleteAt(res)

	return so, totalRecords, nil
}

func (s *roleService) Create(request *requests.CreateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting CreateRole process",
		zap.String("roleName", request.Name),
	)

	res, err := s.roleRepository.CreateRole(request)

	if err != nil {
		s.logger.Error("Failed to create role",
			zap.String("roleName", request.Name),
			zap.Error(err),
		)

		return nil, role_errors.ErrFailedCreateRole
	}

	so := s.mapping.ToRoleResponse(res)

	s.logger.Debug("CreateRole process completed",
		zap.String("roleName", request.Name),
		zap.Int("roleID", res.ID),
	)

	return so, nil
}

func (s *roleService) Update(request *requests.UpdateRoleRequest) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting UpdateRole process",
		zap.Int("roleID", *request.ID),
		zap.String("newRoleName", request.Name),
	)

	res, err := s.roleRepository.UpdateRole(request)

	if err != nil {
		s.logger.Error("Failed to update role",
			zap.Int("roleID", *request.ID),
			zap.String("newRoleName", request.Name),
			zap.Error(err))

		return nil, role_errors.ErrFailedUpdateRole
	}

	so := s.mapping.ToRoleResponse(res)

	s.logger.Debug("UpdateRole process completed",
		zap.Int("roleID", *request.ID),
		zap.String("newRoleName", request.Name),
	)

	return so, nil
}

func (s *roleService) Trashed(id int) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting TrashedRole process",
		zap.Int("roleID", id),
	)

	res, err := s.roleRepository.TrashedRole(id)

	if err != nil {
		s.logger.Error("Failed to move role to trash",
			zap.Int("roleID", id),
			zap.Error(err))

		return nil, role_errors.ErrFailedTrashedRole
	}

	so := s.mapping.ToRoleResponse(res)

	s.logger.Debug("TrashedRole process completed",
		zap.Int("roleID", id),
	)

	return so, nil
}

func (s *roleService) Restore(id int) (*response.RoleResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting RestoreRole process",
		zap.Int("roleID", id),
	)

	res, err := s.roleRepository.RestoreRole(id)

	if err != nil {
		s.logger.Error("Failed to restore role", zap.Error(err))

		return nil, role_errors.ErrFailedRestoreRole
	}

	so := s.mapping.ToRoleResponse(res)

	s.logger.Debug("RestoreRole process completed",
		zap.Int("roleID", id),
	)

	return so, nil
}

func (s *roleService) DeletePermanent(id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Starting DeleteRolePermanent process",
		zap.Int("roleID", id),
	)

	_, err := s.roleRepository.DeleteRolePermanent(id)

	if err != nil {
		s.logger.Error("Failed to delete role permanently",
			zap.Int("roleID", id),
			zap.Error(err),
		)

		return false, role_errors.ErrFailedDeletePermanent
	}

	s.logger.Debug("DeleteRolePermanent process completed",
		zap.Int("roleID", id),
	)

	return true, nil
}

func (s *roleService) RestoreAll() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all roles")

	_, err := s.roleRepository.RestoreAllRole()

	if err != nil {
		s.logger.Error("Failed to restore all roles", zap.Error(err))
		return false, role_errors.ErrFailedRestoreAll
	}

	s.logger.Debug("Successfully restored all roles")
	return true, nil
}

func (s *roleService) DeleteAllPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all roles")

	_, err := s.roleRepository.DeleteAllRolePermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all roles", zap.Error(err))
		return false, role_errors.ErrFailedDeletePermanent
	}

	s.logger.Debug("Successfully deleted all roles permanently")
	return true, nil
}
