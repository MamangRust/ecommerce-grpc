package service

import (
	"context"
	"database/sql"
	auth_cache "ecommerce/internal/cache/auth"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	"ecommerce/internal/errorhandler"
	"ecommerce/internal/repository"
	"ecommerce/pkg/auth"
	db "ecommerce/pkg/database/schema"
	refreshtoken_errors "ecommerce/pkg/errors/refresh_token_errors"
	"ecommerce/pkg/errors/role_errors"
	"ecommerce/pkg/errors/user_errors"
	userrole_errors "ecommerce/pkg/errors/user_role_errors"
	"ecommerce/pkg/hash"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/observability"
	"errors"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type authService struct {
	auth          repository.UserRepository
	refreshToken  repository.RefreshTokenRepository
	userRole      repository.UserRoleRepository
	role          repository.RoleRepository
	hash          hash.HashPassword
	token         auth.TokenManager
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
	cacheIdentity auth_cache.IdentityCache
	cacheLogin    auth_cache.LoginCache
}

type AuthServiceDeps struct {
	AuthRepo         repository.UserRepository
	RefreshTokenRepo repository.RefreshTokenRepository
	RoleRepo         repository.RoleRepository
	UserRoleRepo     repository.UserRoleRepository

	Hash   hash.HashPassword
	Token  auth.TokenManager
	Logger logger.LoggerInterface
	Tracer observability.TraceLoggerObservability

	CacheIdentity auth_cache.IdentityCache
	CacheLogin    auth_cache.LoginCache
}

func NewAuthService(deps AuthServiceDeps) *authService {
	return &authService{
		auth:          deps.AuthRepo,
		refreshToken:  deps.RefreshTokenRepo,
		role:          deps.RoleRepo,
		userRole:      deps.UserRoleRepo,
		hash:          deps.Hash,
		token:         deps.Token,
		logger:        deps.Logger,
		observability: deps.Tracer,
		cacheIdentity: deps.CacheIdentity,
		cacheLogin:    deps.CacheLogin,
	}
}

func (s *authService) Register(ctx context.Context, request *requests.CreateUserRequest) (*db.CreateUserRow, error) {
	const method = "Register"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting user registration",
		zap.String("email", request.Email),
		zap.String("first_name", request.FirstName),
		zap.String("last_name", request.LastName),
	)

	existingUser, err := s.auth.FindByEmail(ctx, request.Email)
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

	passwordHash, err := s.hash.HashPassword(request.Password)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateUserRow](
			s.logger,
			user_errors.ErrUserPassword,
			method,
			span,
		)
	}
	request.Password = passwordHash

	const defaultRoleName = "ROLE_ADMIN"
	role, err := s.role.FindByName(ctx, defaultRoleName)
	if err != nil || role == nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateUserRow](
			s.logger,
			role_errors.ErrRoleNotFoundRes,
			method,
			span,
			zap.String("role", defaultRoleName),
		)
	}

	newUser, err := s.auth.CreateUser(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateUserRow](
			s.logger,
			user_errors.ErrFailedCreateUser,
			method,
			span,
			zap.String("email", request.Email),
		)
	}

	_, err = s.userRole.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: int(newUser.UserID),
		RoleId: int(role.RoleID),
	})
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateUserRow](
			s.logger,
			userrole_errors.ErrFailedAssignRoleToUser,
			method,
			span,
			zap.Int("user_id", int(newUser.UserID)),
			zap.Int("role_id", int(role.RoleID)),
		)
	}

	logSuccess("User registered successfully",
		zap.Int("user_id", int(newUser.UserID)),
		zap.String("email", request.Email),
	)

	return newUser, nil
}

func (s *authService) Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, error) {
	const method = "Login"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Starting login process",
		zap.String("email", request.Email),
	)

	res, err := s.auth.FindByEmailWithPassword(ctx, request.Email)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.TokenResponse](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,
			zap.String("email", request.Email),
		)
	}

	err = s.hash.ComparePassword(res.Password, request.Password)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.TokenResponse](
			s.logger,
			user_errors.ErrUserPassword,
			method,
			span,
			zap.String("email", request.Email),
		)
	}

	token, err := s.createAccessToken(int(res.UserID))
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.TokenResponse](
			s.logger,
			refreshtoken_errors.ErrFailedCreateAccess,
			method,
			span,
			zap.Int("user_id", int(res.UserID)),
		)
	}

	refreshToken, err := s.createRefreshToken(ctx, int(res.UserID))
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.TokenResponse](
			s.logger,
			refreshtoken_errors.ErrFailedCreateRefresh,
			method,
			span,
			zap.Int("user_id", int(res.UserID)),
		)
	}

	logSuccess("User logged in successfully", zap.String("email", request.Email))

	return &response.TokenResponse{AccessToken: token, RefreshToken: refreshToken}, nil
}

func (s *authService) RefreshToken(ctx context.Context, token string) (*response.TokenResponse, error) {
	const method = "RefreshToken"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("token", token))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Refreshing token", zap.String("token", token))

	userIdStr, err := s.token.ValidateToken(token)
	if err != nil {
		status = "error"
		if errors.Is(err, auth.ErrTokenExpired) {
			if err := s.refreshToken.DeleteRefreshToken(ctx, token); err != nil {
				return errorhandler.HandleError[*response.TokenResponse](
					s.logger,
					refreshtoken_errors.ErrFailedDeleteRefreshToken,
					method,
					span,
					zap.String("token", token),
				)
			}
			return errorhandler.HandleError[*response.TokenResponse](
				s.logger,
				refreshtoken_errors.ErrFailedExpire,
				method,
				span,
				zap.String("token", token),
			)
		}
		return errorhandler.HandleError[*response.TokenResponse](
			s.logger,
			refreshtoken_errors.ErrRefreshTokenNotFound,
			method,
			span,
			zap.String("token", token),
		)
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.TokenResponse](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,
			zap.String("user_id_str", userIdStr),
		)
	}

	accessToken, err := s.createAccessToken(userId)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.TokenResponse](
			s.logger,
			refreshtoken_errors.ErrFailedCreateAccess,
			method,
			span,
			zap.Int("user_id", userId),
		)
	}

	refreshToken, err := s.createRefreshToken(ctx, userId)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*response.TokenResponse](
			s.logger,
			refreshtoken_errors.ErrFailedCreateRefreshToken,
			method,
			span,
			zap.Int("user_id", userId),
		)
	}

	expiryTime := time.Now().Add(24 * time.Hour)
	updateRequest := &requests.UpdateRefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: expiryTime.Format("2006-01-02 15:04:05"),
	}

	if _, err = s.refreshToken.UpdateRefreshToken(ctx, updateRequest); err != nil {
		status = "error"
		return errorhandler.HandleError[*response.TokenResponse](
			s.logger,
			refreshtoken_errors.ErrFailedUpdateRefreshToken,
			method,
			span,
			zap.Int("user_id", userId),
		)
	}

	logSuccess("Refresh token refreshed successfully")

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) GetMe(ctx context.Context, token string) (*db.GetUserByIDRow, error) {
	const method = "GetMe"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("token", token))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Fetching user details", zap.String("token", token))

	userIdStr, err := s.token.ValidateToken(token)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetUserByIDRow](
			s.logger,
			refreshtoken_errors.ErrFailedInValidToken,
			method,
			span,
			zap.String("token", token),
		)
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetUserByIDRow](
			s.logger,
			refreshtoken_errors.ErrFailedInValidUserId,
			method,
			span,
			zap.String("user_id_str", userIdStr),
		)
	}

	user, err := s.auth.FindById(ctx, userId)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetUserByIDRow](
			s.logger,
			user_errors.ErrUserNotFoundRes,
			method,
			span,
			zap.Int("user_id", userId),
		)
	}

	logSuccess("User details fetched successfully", zap.Int("userID", userId))

	return user, nil
}

func (s *authService) createAccessToken(id int) (string, error) {
	s.logger.Debug("Creating access token",
		zap.Int("userID", id),
	)

	res, err := s.token.GenerateToken(id, "access")

	if err != nil {
		s.logger.Error("Failed to create access token",
			zap.Int("userID", id),
			zap.Error(err))
		return "", err
	}

	s.logger.Debug("Access token created successfully",
		zap.Int("userID", id),
	)

	return res, nil
}

func (s *authService) createRefreshToken(ctx context.Context, id int) (string, error) {
	s.logger.Debug("Creating refresh token",
		zap.Int("userID", id),
	)

	res, err := s.token.GenerateToken(id, "refresh")

	if err != nil {
		s.logger.Error("Failed to create refresh token",
			zap.Int("userID", id),
		)

		return "", err
	}

	if err := s.refreshToken.DeleteRefreshTokenByUserId(ctx, id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("Failed to delete existing refresh token", zap.Error(err))
		return "", err
	}

	_, err = s.refreshToken.CreateRefreshToken(ctx, &requests.CreateRefreshToken{Token: res, UserId: id, ExpiresAt: time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")})
	if err != nil {
		s.logger.Error("Failed to create refresh token", zap.Error(err))

		return "", err
	}

	s.logger.Debug("Refresh token created successfully",
		zap.Int("userID", id),
	)

	return res, nil
}

func maskToken(token string) string {
	if len(token) < 8 {
		return "******"
	}
	return token[:4] + "****" + token[len(token)-4:]
}
