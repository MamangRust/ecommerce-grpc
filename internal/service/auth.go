package service

import (
	"database/sql"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/auth"
	refreshtoken_errors "ecommerce/pkg/errors/refresh_token_errors"
	"ecommerce/pkg/errors/user_errors"
	userrole_errors "ecommerce/pkg/errors/user_role_errors"
	"ecommerce/pkg/hash"
	"ecommerce/pkg/logger"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type authService struct {
	auth         repository.UserRepository
	refreshToken repository.RefreshTokenRepository
	userRole     repository.UserRoleRepository
	role         repository.RoleRepository
	hash         hash.HashPassword
	token        auth.TokenManager
	logger       logger.LoggerInterface
	mapping      response_service.UserResponseMapper
}

func NewAuthService(auth repository.UserRepository, refreshToken repository.RefreshTokenRepository, role repository.RoleRepository, userRole repository.UserRoleRepository, hash hash.HashPassword, token auth.TokenManager, logger logger.LoggerInterface, mapping response_service.UserResponseMapper) *authService {
	return &authService{auth: auth, refreshToken: refreshToken, role: role, userRole: userRole, hash: hash, token: token, logger: logger, mapping: mapping}
}

func (s *authService) Register(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting user registration",
		zap.String("email", request.Email),
		zap.String("first_name", request.FirstName),
		zap.String("last_name", request.LastName),
	)

	existingUser, err := s.auth.FindByEmail(request.Email)
	if err == nil && existingUser != nil {
		s.logger.Debug("Email already exists",
			zap.String("email", request.Email),
		)
		return nil, user_errors.ErrUserEmailAlready
	}

	passwordHash, err := s.hash.HashPassword(request.Password)
	if err != nil {
		s.logger.Error("Failed to hash password",
			zap.Error(err),
		)
		return nil, user_errors.ErrUserPassword
	}
	request.Password = passwordHash

	const defaultRoleName = "Cashier"
	role, err := s.role.FindByName(defaultRoleName)
	if err != nil || role == nil {
		s.logger.Error("Failed to find default role",
			zap.String("role", defaultRoleName),
			zap.Error(err),
		)
		return nil, &response.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Status:  "server_error",
			Message: "Failed to assign user role",
		}
	}
	newUser, err := s.auth.CreateUser(request)
	if err != nil {
		s.logger.Error("Failed to create user",
			zap.String("email", request.Email),
			zap.Error(err),
		)
		return nil, user_errors.ErrFailedCreateUser
	}

	_, err = s.userRole.AssignRoleToUser(&requests.CreateUserRoleRequest{
		UserId: newUser.ID,
		RoleId: role.ID,
	})
	if err != nil {
		s.logger.Error("Failed to assign role to user",
			zap.Int("user_id", newUser.ID),
			zap.Int("role_id", role.ID),
			zap.Error(err),
		)
		return nil, userrole_errors.ErrFailedAssignRoleToUser
	}

	userResponse := s.mapping.ToUserResponse(newUser)

	s.logger.Debug("User registered successfully",
		zap.Int("user_id", newUser.ID),
		zap.String("email", request.Email),
	)

	return userResponse, nil
}

func (s *authService) Login(request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse) {
	s.logger.Debug("Starting login process",
		zap.String("email", request.Email),
	)

	res, err := s.auth.FindByEmail(request.Email)
	if err != nil {
		s.logger.Error("Failed to get user",
			zap.String("email", request.Email),
			zap.Error(err),
		)
		return nil, user_errors.ErrUserNotFoundRes
	}
	if err := s.hash.ComparePassword(res.Password, request.Password); err != nil {
		s.logger.Debug("Invalid password attempt",
			zap.String("email", request.Email),
			zap.Error(err),
		)
		return nil, user_errors.ErrUserPassword
	}

	token, err := s.createAccessToken(res.ID)
	if err != nil {
		s.logger.Error("Failed to generate JWT token",
			zap.Int("user_id", res.ID),
			zap.Error(err),
		)
		return nil, refreshtoken_errors.ErrFailedCreateAccess
	}

	refreshToken, err := s.createRefreshToken(res.ID)
	if err != nil {
		s.logger.Error("Failed to generate refresh token",
			zap.Int("user_id", res.ID),
			zap.Error(err),
		)
		return nil, refreshtoken_errors.ErrFailedCreateRefresh
	}

	s.logger.Debug("User logged in successfully",
		zap.Int("user_id", res.ID),
		zap.String("email", request.Email),
	)

	return &response.TokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(token string) (*response.TokenResponse, *response.ErrorResponse) {
	s.logger.Debug("Refreshing token",
		zap.String("token", token),
	)

	userIdStr, err := s.token.ValidateToken(token)

	if err != nil {
		if errors.Is(err, auth.ErrTokenExpired) {
			if err := s.refreshToken.DeleteRefreshToken(token); err != nil {
				s.logger.Error("Failed to delete expired refresh token", zap.Error(err))

				return nil, refreshtoken_errors.ErrFailedDeleteRefreshToken
			}

			s.logger.Error("Refresh token has expired", zap.Error(err))

			return nil, refreshtoken_errors.ErrFailedExpire
		}
		s.logger.Error("Invalid refresh token", zap.Error(err))
		return nil, refreshtoken_errors.ErrRefreshTokenNotFound
	}

	userId, err := strconv.Atoi(userIdStr)

	if err != nil {
		s.logger.Error("Invalid user ID format in token", zap.Error(err))
		return nil, user_errors.ErrUserIDInValid
	}

	accessToken, err := s.createAccessToken(userId)
	if err != nil {
		s.logger.Error("Failed to generate new access token", zap.Error(err))

		return nil, refreshtoken_errors.ErrFailedCreateAccess
	}

	refreshToken, err := s.createRefreshToken(userId)
	if err != nil {
		s.logger.Error("Failed to generate new refresh token", zap.Error(err))

		return nil, refreshtoken_errors.ErrFailedCreateRefreshToken
	}

	expiryTime := time.Now().Add(24 * time.Hour)

	updateRequest := &requests.UpdateRefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: expiryTime.Format("2006-01-02 15:04:05"),
	}

	if _, err = s.refreshToken.UpdateRefreshToken(updateRequest); err != nil {
		s.logger.Error("Failed to update refresh token in storage", zap.Error(err))

		return nil, refreshtoken_errors.ErrFailedUpdateRefreshToken
	}

	s.logger.Debug("Refresh token refreshed successfully")

	return &response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) GetMe(token string) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching user details",
		zap.String("token", token),
	)

	userIdStr, err := s.token.ValidateToken(token)

	if err != nil {
		s.logger.Error("Invalid access token", zap.Error(err))
		return nil, refreshtoken_errors.ErrFailedInValidToken
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		s.logger.Error("Invalid user ID format in token", zap.Error(err))
		return nil, refreshtoken_errors.ErrFailedInValidUserId
	}

	user, err := s.auth.FindById(userId)

	if err != nil {
		s.logger.Error("Failed to find user by ID", zap.Error(err))
		return nil, user_errors.ErrUserNotFoundRes
	}

	so := s.mapping.ToUserResponse(user)

	s.logger.Debug("User details fetched successfully",
		zap.Int("userID", userId),
	)

	return so, nil
}

func maskToken(token string) string {
	if len(token) < 8 {
		return "******"
	}
	return token[:4] + "****" + token[len(token)-4:]
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

func (s *authService) createRefreshToken(id int) (string, error) {
	s.logger.Debug("Creating refresh token",
		zap.Int("userID", id),
	)

	res, err := s.token.GenerateToken(id, "refresh")

	if err != nil {
		s.logger.Error("Failed to create refresh token",
			zap.Int("userID", id),
			zap.Error(err),
		)

		return "", err
	}

	if err := s.refreshToken.DeleteRefreshTokenByUserId(id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("Failed to delete existing refresh token", zap.Error(err))
		return "", err
	}

	_, err = s.refreshToken.CreateRefreshToken(&requests.CreateRefreshToken{Token: res, UserId: id, ExpiresAt: time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")})
	if err != nil {
		s.logger.Error("Failed to create refresh token", zap.Error(err))

		return "", err
	}

	s.logger.Debug("Refresh token created successfully",
		zap.Int("userID", id),
	)

	return res, nil
}
