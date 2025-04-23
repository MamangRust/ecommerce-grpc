package service

import (
	"database/sql"
	"ecommerce/internal/domain/requests"
	"ecommerce/internal/domain/response"
	response_service "ecommerce/internal/mapper/response/services"
	"ecommerce/internal/repository"
	"ecommerce/pkg/hash"
	"ecommerce/pkg/logger"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type userService struct {
	userRepository repository.UserRepository
	logger         logger.LoggerInterface
	mapping        response_service.UserResponseMapper
	hashing        hash.HashPassword
}

func NewUserService(
	userRepository repository.UserRepository,
	logger logger.LoggerInterface,
	mapper response_service.UserResponseMapper,
	hashing hash.HashPassword,
) *userService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
		mapping:        mapper,
		hashing:        hashing,
	}
}

func (s *userService) FindAll(req *requests.FindAllUsers) ([]*response.UserResponse, *int, *response.ErrorResponse) {
	s.logger.Debug("Fetching users",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	users, totalRecords, err := s.userRepository.FindAllUsers(req)

	if err != nil {
		s.logger.Error("Failed to fetch user",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return nil, nil, &response.ErrorResponse{
			Status:  "fail",
			Message: "Unable to retrieve user list at this time",
			Code:    http.StatusInternalServerError,
		}
	}

	userResponses := s.mapping.ToUsersResponse(users)

	s.logger.Debug("Successfully fetched user",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return userResponses, totalRecords, nil
}

func (s *userService) FindByID(id int) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching user by id", zap.Int("user_id", id))

	user, err := s.userRepository.FindById(id)

	if err != nil {
		s.logger.Error("failed to find user by ID", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "User not found",
		}
	}

	so := s.mapping.ToUserResponse(user)

	s.logger.Debug("Successfully fetched user", zap.Int("user_id", id))

	return so, nil
}

func (s *userService) FindByActive(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	s.logger.Debug("Fetching active user",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	users, totalRecords, err := s.userRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active users from database",
			zap.String("search", req.Search),
			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.Error(err))
		return nil, nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to retrieve active users",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.logger.Debug("Successfully fetched active user",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return so, totalRecords, nil
}

func (s *userService) FindByTrashed(req *requests.FindAllUsers) ([]*response.UserResponseDeleteAt, *int, *response.ErrorResponse) {
	s.logger.Debug("Fetching trashed user",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	users, totalRecords, err := s.userRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve archived users",
			zap.Error(err),
			zap.String("search", req.Search),
			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
		)

		return nil, nil, &response.ErrorResponse{
			Status:  "fail",
			Message: "Unable to retrieve archived user accounts",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.logger.Debug("Successfully fetched trashed user",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize),
	)

	return so, totalRecords, nil
}

func (s *userService) CreateUser(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new user", zap.String("email", request.Email), zap.Any("request", request))

	existingUser, err := s.userRepository.FindByEmail(request.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("Error while checking email availability",
			zap.String("email", request.Email),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "System error while validating email",
			Code:    http.StatusInternalServerError,
		}
	}

	if existingUser != nil {
		s.logger.Fatal("Email already registered",
			zap.String("email", request.Email))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "This email address is already registered",
			Code:    http.StatusConflict,
		}
	}

	hash, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		s.logger.Error("Password hashing failed",
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "System error during account creation",
			Code:    http.StatusInternalServerError,
		}
	}
	request.Password = hash

	res, err := s.userRepository.CreateUser(request)
	if err != nil {
		s.logger.Error("User creation failed",
			zap.String("email", request.Email),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create user account",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToUserResponse(res)

	s.logger.Debug("User created successfully",
		zap.String("email", so.Email),
		zap.Int("user_id", so.ID))

	return so, nil
}

func (s *userService) UpdateUser(request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating user",
		zap.Int("user_id", *request.UserID),
		zap.Any("request", request))

	existingUser, err := s.userRepository.FindById(*request.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("User not found for update",
				zap.Int("user_id", *request.UserID))
			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: "User account not found",
				Code:    http.StatusNotFound,
			}
		}
		s.logger.Error("Error finding user",
			zap.Int("user_id", *request.UserID),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "System error while retrieving user",
			Code:    http.StatusInternalServerError,
		}
	}

	if request.Email != existingUser.Email {
		duplicateUser, err := s.userRepository.FindByEmail(request.Email)

		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				s.logger.Error("Error checking email availability",
					zap.String("email", request.Email),
					zap.Error(err))
				return nil, &response.ErrorResponse{
					Status:  "error",
					Message: "System error while validating email",
					Code:    http.StatusInternalServerError,
				}
			}
		}

		if duplicateUser != nil && duplicateUser.ID != existingUser.ID {
			s.logger.Debug("Email already in use by another user",
				zap.String("email", request.Email),
				zap.Int("existing_user_id", existingUser.ID),
				zap.Int("duplicate_user_id", duplicateUser.ID))

			return nil, &response.ErrorResponse{
				Status:  "conflict",
				Message: "This email address is already registered by another user",
				Code:    http.StatusConflict,
			}
		}
	}

	hash, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		s.logger.Error("Password hashing failed",
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "System error during password update",
			Code:    http.StatusInternalServerError,
		}
	}
	existingUser.Password = hash

	res, err := s.userRepository.UpdateUser(request)

	if err != nil {
		s.logger.Error("User update failed",
			zap.Int("user_id", *request.UserID),
			zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update user account",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToUserResponse(res)
	s.logger.Debug("User updated successfully",
		zap.Int("user_id", so.ID))

	return so, nil
}

func (s *userService) TrashedUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing user", zap.Int("user_id", user_id))

	res, err := s.userRepository.TrashedUser(user_id)

	if err != nil {
		s.logger.Error("Failed to move user to trash",
			zap.Int("user_id", user_id),
			zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("user with ID %d not found", user_id),
				Code:    http.StatusNotFound,
			}
		}
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to move user to trash",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToUserResponseDeleteAt(res)

	s.logger.Debug("Successfully trashed user", zap.Int("user_id", user_id))

	return so, nil
}

func (s *userService) RestoreUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Restoring user", zap.Int("user_id", user_id))

	res, err := s.userRepository.RestoreUser(user_id)

	if err != nil {
		s.logger.Error("Failed to restore user from trash",
			zap.Int("user_id", user_id),
			zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("user with ID %d not found in trash", user_id),
				Code:    http.StatusNotFound,
			}
		}
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore user from trash",
			Code:    http.StatusInternalServerError,
		}
	}

	so := s.mapping.ToUserResponseDeleteAt(res)

	s.logger.Debug("Successfully restored user", zap.Int("user_id", user_id))

	return so, nil
}

func (s *userService) DeleteUserPermanent(user_id int) (bool, *response.ErrorResponse) {
	s.logger.Debug("Deleting user permanently", zap.Int("user_id", user_id))

	_, err := s.userRepository.DeleteUserPermanent(user_id)

	if err != nil {
		s.logger.Error("Failed to permanently delete user",
			zap.Int("user_id", user_id),
			zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return false, &response.ErrorResponse{
				Status:  "not_found",
				Message: fmt.Sprintf("user with ID %d not found", user_id),
				Code:    http.StatusNotFound,
			}
		}
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete user",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully deleted user permanently", zap.Int("user_id", user_id))

	return true, nil
}

func (s *userService) RestoreAllUser() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all users")

	_, err := s.userRepository.RestoreAllUser()

	if err != nil {
		s.logger.Error("Failed to restore all trashed users",
			zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all users",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully restored all users")

	return true, nil
}

func (s *userService) DeleteAllUserPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all users")

	_, err := s.userRepository.DeleteAllUserPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all trashed users",
			zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all users",
			Code:    http.StatusInternalServerError,
		}
	}

	s.logger.Debug("Successfully deleted all users permanently")

	return true, nil
}
