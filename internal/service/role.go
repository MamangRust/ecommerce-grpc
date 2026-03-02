package service

import (
	"context"
	role_cache "ecommerce/internal/cache/role"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/role_errors"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type roleService struct {
	roleRepository repository.RoleRepository
	logger         logger.LoggerInterface
	observability  observability.TraceLoggerObservability
	cache          role_cache.RoleMencache
}

type RoleServiceDeps struct {
	RoleRepo      repository.RoleRepository
	Observability observability.TraceLoggerObservability
	Logger        logger.LoggerInterface
	Cache         role_cache.RoleMencache
}

func NewRoleService(deps RoleServiceDeps) *roleService {
	return &roleService{
		roleRepository: deps.RoleRepo,
		logger:         deps.Logger,
		observability:  deps.Observability,
		cache:          deps.Cache,
	}
}

func (s *roleService) FindAll(ctx context.Context, req *requests.FindAllRole) ([]*db.GetRolesRow, *int, error) {
	const method = "FindAll"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedRoles(ctx, req); found {
		logSuccess("Successfully retrieved all role records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, err := s.roleRepository.FindAllRoles(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetRolesRow](
			s.logger,
			role_errors.ErrFailedFindAll,
			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedRoles(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched role",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return res, &totalCount, nil
}

func (s *roleService) FindById(ctx context.Context, id int) (*db.Role, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedRoleById(ctx, id); found {
		logSuccess("Successfully retrieved role from cache", zap.Int("id", id))

		return data, nil
	}

	res, err := s.roleRepository.FindById(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrRoleNotFoundRes,
			method,
			span,

			zap.Int("role_id", id),
		)
	}

	s.cache.SetCachedRoleById(ctx, id, res)

	logSuccess("Successfully fetched role", zap.Int("id", id))

	return res, nil
}

func (s *roleService) FindByUserId(ctx context.Context, id int) ([]*db.Role, error) {
	const method = "FindByUserId"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedRoleByUserId(ctx, id); found {
		logSuccess("Successfully fetched role by user ID from cache", zap.Int("id", id))
		return data, nil
	}

	res, err := s.roleRepository.FindByUserId(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.Role](
			s.logger,
			role_errors.ErrRoleNotFoundRes,
			method,
			span,

			zap.Int("user_id", id),
		)
	}

	s.cache.SetCachedRoleByUserId(ctx, id, res)

	logSuccess("Successfully fetched role by user ID", zap.Int("id", id))

	return res, nil
}

func (s *roleService) FindByActiveRole(ctx context.Context, req *requests.FindAllRole) ([]*db.GetActiveRolesRow, *int, error) {
	const method = "FindByActiveRole"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	res, err := s.roleRepository.FindByActiveRole(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetActiveRolesRow](
			s.logger,
			role_errors.ErrFailedFindActive,
			method,
			span,

			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
		)
	}
	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedRoleActive(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched active role",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return res, &totalCount, nil
}

func (s *roleService) FindByTrashedRole(ctx context.Context, req *requests.FindAllRole) ([]*db.GetTrashedRolesRow, *int, error) {
	const method = "FindByTrashedRole"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedRoleTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed role from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	res, err := s.roleRepository.FindByTrashedRole(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTrashedRolesRow](
			s.logger,
			role_errors.ErrFailedFindTrashed,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedRoleTrashed(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched trashed role",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return res, &totalCount, nil
}

func (s *roleService) CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*db.Role, error) {
	const method = "CreateRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("roleName", request.Name))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting CreateRole process",
		zap.String("roleName", request.Name),
	)

	role, err := s.roleRepository.CreateRole(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrFailedCreateRole,
			method,
			span,
			zap.String("role_name", request.Name),
		)
	}

	logSuccess("CreateRole process completed",
		zap.String("roleName", request.Name),
		zap.Int("roleID", int(role.RoleID)),
	)

	return role, nil
}

func (s *roleService) UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*db.Role, error) {
	const method = "UpdateRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("roleID", *request.ID),
		attribute.String("newRoleName", request.Name))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting UpdateRole process",
		zap.Int("roleID", *request.ID),
		zap.String("newRoleName", request.Name),
	)

	role, err := s.roleRepository.UpdateRole(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrFailedUpdateRole,
			method,
			span,
			zap.Int("role_id", *request.ID),
			zap.String("new_name", request.Name),
		)
	}

	logSuccess("UpdateRole process completed",
		zap.Int("roleID", *request.ID),
		zap.String("newRoleName", request.Name),
	)

	return role, nil
}

func (s *roleService) TrashedRole(ctx context.Context, id int) (*db.Role, error) {
	const method = "TrashedRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("roleID", id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting TrashedRole process",
		zap.Int("roleID", id),
	)

	role, err := s.roleRepository.TrashedRole(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrFailedTrashedRole,
			method,
			span,
			zap.Int("role_id", id),
		)
	}

	logSuccess("TrashedRole process completed",
		zap.Int("roleID", id),
	)

	return role, nil
}

func (s *roleService) RestoreRole(ctx context.Context, id int) (*db.Role, error) {
	const method = "RestoreRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("roleID", id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting RestoreRole process",
		zap.Int("roleID", id),
	)

	role, err := s.roleRepository.RestoreRole(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Role](
			s.logger,
			role_errors.ErrFailedRestoreRole,
			method,
			span,
			zap.Int("role_id", id),
		)
	}

	logSuccess("RestoreRole process completed",
		zap.Int("roleID", id),
	)

	return role, nil
}

func (s *roleService) DeleteRolePermanent(ctx context.Context, id int) (bool, error) {
	const method = "DeleteRolePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("roleID", id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting DeleteRolePermanent process",
		zap.Int("roleID", id),
	)

	_, err := s.roleRepository.DeleteRolePermanent(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			role_errors.ErrFailedDeletePermanent,
			method,
			span,
			zap.Int("role_id", id),
		)
	}

	logSuccess("DeleteRolePermanent process completed",
		zap.Int("roleID", id),
	)

	return true, nil
}

func (s *roleService) RestoreAllRole(ctx context.Context) (bool, error) {
	const method = "RestoreAllRole"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring all roles")

	_, err := s.roleRepository.RestoreAllRole(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			role_errors.ErrFailedRestoreAll,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all roles")
	return true, nil
}

func (s *roleService) DeleteAllRolePermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllRolePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Permanently deleting all roles")

	_, err := s.roleRepository.DeleteAllRolePermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			role_errors.ErrFailedDeletePermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all roles permanently")
	return true, nil
}
