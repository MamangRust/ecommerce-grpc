package service

import (
	"context"
	user_cache "ecommerce/internal/cache/user"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	db "ecommerce/pkg/database/schema"
	"ecommerce/pkg/errors/user_errors"
	"ecommerce/pkg/hash"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type userService struct {
	userRepository repository.UserRepository
	logger         logger.LoggerInterface
	observability  observability.TraceLoggerObservability
	hashing        hash.HashPassword
	cache          user_cache.UserMencache
}

type UserServiceDeps struct {
	UserRepo repository.UserRepository

	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability

	Hashing hash.HashPassword
	Cache   user_cache.UserMencache
}

func NewUserService(deps UserServiceDeps) *userService {
	return &userService{
		userRepository: deps.UserRepo,
		logger:         deps.Logger,
		observability:  deps.Observability,
		hashing:        deps.Hashing,
		cache:          deps.Cache,
	}
}

func (s *userService) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedUsersCache(ctx, req); found {
		logSuccess("Successfully retrieved all user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	users, err := s.userRepository.FindAllUsers(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetUsersRow](
			s.logger,
			user_errors.ErrFailedFindAll,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(users) > 0 {
		totalCount = int(users[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedUsersCache(ctx, req, users, &totalCount)

	logSuccess("Successfully fetched user",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return users, &totalCount, nil
}

func (s *userService) FindByID(ctx context.Context, id int) (*db.GetUserByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedUserCache(ctx, id); found {
		logSuccess("Successfully retrieved user record from cache", zap.Int("user.id", id))
		return data, nil
	}

	user, err := s.userRepository.FindById(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetUserByIDRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,

			zap.Int("user_id", id),
		)
	}

	s.cache.SetCachedUserCache(ctx, user)

	logSuccess("Successfully fetched user", zap.Int("user_id", id))

	return user, nil
}

func (s *userService) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUsersActiveRow, *int, error) {
	const method = "FindByActive"

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

	if data, total, found := s.cache.GetCachedUserActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active user records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	users, err := s.userRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetUsersActiveRow](
			s.logger,
			user_errors.ErrFailedFindActive,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(users) > 0 {
		totalCount = int(users[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedUserActiveCache(ctx, req, users, &totalCount)

	logSuccess("Successfully fetched active user",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return users, &totalCount, nil
}

func (s *userService) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*db.GetUserTrashedRow, *int, error) {
	const method = "FindByTrashed"

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

	if data, total, found := s.cache.GetCachedUserTrashedCache(ctx, req); found {
		return data, total, nil
	}

	users, err := s.userRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetUserTrashedRow](
			s.logger,
			user_errors.ErrFailedFindTrashed,
			method,
			span,

			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
		)
	}

	var totalCount int

	if len(users) > 0 {
		totalCount = int(users[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedUserTrashedCache(ctx, req, users, &totalCount)

	logSuccess("Successfully fetched trashed user",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return users, &totalCount, nil
}

func (s *userService) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error) {
	const method = "CreateUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Creating new user", zap.String("email", request.Email), zap.Any("request", request))

	existingUser, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err == nil && existingUser != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateUserRow](
			s.logger,
			user_errors.ErrUserEmailAlready,
			method,
			span,
			zap.String("email", request.Email),
		)
	}

	hash, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateUserRow](
			s.logger,
			user_errors.ErrUserPassword,
			method,
			span,
		)
	}

	request.Password = hash

	res, err := s.userRepository.CreateUser(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateUserRow](
			s.logger,
			user_errors.ErrFailedCreateUser,
			method,
			span,
		)
	}

	logSuccess("Successfully created new user", zap.String("email", res.Email), zap.Int("user_id", int(res.UserID)))

	return res, nil
}

func (s *userService) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*db.UpdateUserRow, error) {
	const method = "UpdateUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", *request.UserID))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Updating user", zap.Int("user_id", *request.UserID), zap.Any("request", request))

	existingUser, err := s.userRepository.FindByIdWithPassword(ctx, *request.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateUserRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,

			zap.Int("user_id", *request.UserID),
		)
	}

	if request.Email != "" && request.Email != existingUser.Email {
		duplicateUser, _ := s.userRepository.FindByEmail(ctx, request.Email)
		if duplicateUser != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateUserRow](
				s.logger,
				user_errors.ErrUserEmailAlready,
				method,
				span,
				zap.String("email", request.Email),
			)
		}
		existingUser.Email = request.Email
	}

	if request.Password != "" {
		hash, err := s.hashing.HashPassword(request.Password)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateUserRow](
				s.logger,
				user_errors.ErrUserPassword,
				method,
				span,
			)
		}
		existingUser.Password = hash
	}

	res, err := s.userRepository.UpdateUser(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateUserRow](
			s.logger,
			user_errors.ErrFailedUpdateUser,
			method,
			span,

			zap.Int("user_id", *request.UserID),
		)
	}

	logSuccess("Successfully updated user", zap.Int("user_id", int(res.UserID)))

	return res, nil
}

func (s *userService) TrashedUser(ctx context.Context, user_id int) (*db.TrashUserRow, error) {
	const method = "TrashedUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", user_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Trashing user", zap.Int("user_id", user_id))

	res, err := s.userRepository.TrashedUser(ctx, user_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.TrashUserRow](
			s.logger,
			user_errors.ErrFailedTrashedUser,
			method,
			span,

			zap.Int("user_id", user_id),
		)
	}

	logSuccess("Successfully trashed user", zap.Int("user_id", user_id))

	return res, nil
}

func (s *userService) RestoreUser(ctx context.Context, user_id int) (*db.RestoreUserRow, error) {
	const method = "RestoreUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", user_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring user", zap.Int("user_id", user_id))

	res, err := s.userRepository.RestoreUser(ctx, user_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.RestoreUserRow](
			s.logger,
			user_errors.ErrFailedRestoreUser,
			method,
			span,

			zap.Int("user_id", user_id),
		)
	}

	logSuccess("Successfully restored user", zap.Int("user_id", user_id))

	return res, nil
}

func (s *userService) DeleteUserPermanent(ctx context.Context, user_id int) (bool, error) {
	const method = "DeleteUserPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", user_id))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Deleting user permanently", zap.Int("user_id", user_id))

	_, err := s.userRepository.DeleteUserPermanent(ctx, user_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			user_errors.ErrFailedDeletePermanent,
			method,
			span,

			zap.Int("user_id", user_id),
		)
	}

	logSuccess("Successfully deleted user permanently", zap.Int("user_id", user_id))

	return true, nil
}

func (s *userService) RestoreAllUser(ctx context.Context) (bool, error) {
	const method = "RestoreAllUser"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Restoring all users")

	_, err := s.userRepository.RestoreAllUser(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			user_errors.ErrFailedRestoreAll,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all users")

	return true, nil
}

func (s *userService) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllUserPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	s.logger.Debug("Permanently deleting all users")

	_, err := s.userRepository.DeleteAllUserPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			user_errors.ErrFailedDeleteAll,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all users permanently")

	return true, nil
}
