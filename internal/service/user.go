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

func (s *userService) FindAll(page int, pageSize int, search string) ([]*response.UserResponse, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching users",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.userRepository.FindAllUsers(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch user",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to fetch users",
		}
	}

	userResponses := s.mapping.ToUsersResponse(users)

	s.logger.Debug("Successfully fetched user",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return userResponses, int(totalRecords), nil
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

func (s *userService) FindByActive(page int, pageSize int, search string) ([]*response.UserResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching active user",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.userRepository.FindByActive(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to fetch active user",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find active users",
		}
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.logger.Debug("Successfully fetched active user",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *userService) FindByTrashed(page int, pageSize int, search string) ([]*response.UserResponseDeleteAt, int, *response.ErrorResponse) {
	s.logger.Debug("Fetching trashed user",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	users, totalRecords, err := s.userRepository.FindByTrashed(search, page, pageSize)

	if err != nil {
		s.logger.Error("Failed to find trashed users", zap.Error(err))

		return nil, 0, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to find trashed users",
		}
	}

	so := s.mapping.ToUsersResponseDeleteAt(users)

	s.logger.Debug("Successfully fetched trashed user",
		zap.Int("totalRecords", totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *userService) CreateUser(request *requests.CreateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Debug("Creating new user", zap.String("email", request.Email), zap.Any("request", request))

	existingUser, err := s.userRepository.FindByEmail(request.Email)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Debug("Email is available, proceeding to create user", zap.String("email", request.Email))
		} else {
			s.logger.Error("Error checking existing email", zap.String("email", request.Email), zap.Error(err))
			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: "Error checking existing email",
			}
		}
	} else if existingUser != nil {
		s.logger.Error("Email is already in use", zap.String("email", request.Email))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Email is already in use",
		}
	}

	// Hash password
	hash, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to hash password",
		}
	}

	request.Password = hash

	res, err := s.userRepository.CreateUser(request)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to create user",
		}
	}

	// Mapping hasil ke response
	so := s.mapping.ToUserResponse(res)

	s.logger.Debug("Successfully created new user", zap.String("email", so.Email), zap.Int("user", so.ID))

	return so, nil
}

func (s *userService) UpdateUser(request *requests.UpdateUserRequest) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Debug("Updating user", zap.Int("user_id", request.UserID), zap.Any("request", request))

	existingUser, err := s.userRepository.FindById(request.UserID)

	if err != nil {
		s.logger.Error("Failed to find user by ID", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "User not found",
		}
	}

	if request.Email != "" && request.Email != existingUser.Email {
		duplicateUser, _ := s.userRepository.FindByEmail(request.Email)

		if duplicateUser != nil {
			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: "Email is already in use",
			}
		}

		existingUser.Email = request.Email
	}

	if request.Password != "" {
		hash, err := s.hashing.HashPassword(request.Password)
		if err != nil {
			s.logger.Error("Failed to hash password", zap.Error(err))
			return nil, &response.ErrorResponse{
				Status:  "error",
				Message: "Failed to hash password",
			}
		}
		existingUser.Password = hash
	}

	res, err := s.userRepository.UpdateUser(request)
	if err != nil {
		s.logger.Error("Failed to update user", zap.Error(err))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to update user",
		}
	}

	so := s.mapping.ToUserResponse(res)

	s.logger.Debug("Successfully updated user", zap.Int("user_id", so.ID))

	return so, nil
}

func (s *userService) TrashedUser(user_id int) (*response.UserResponseDeleteAt, *response.ErrorResponse) {
	s.logger.Debug("Trashing user", zap.Int("user_id", user_id))

	res, err := s.userRepository.TrashedUser(user_id)

	if err != nil {
		s.logger.Error("Failed to trash user", zap.Error(err), zap.Int("user_id", user_id))
		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to trash user",
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
		s.logger.Error("Failed to restore user", zap.Error(err), zap.Int("user_id", user_id))

		return nil, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore user",
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
		s.logger.Error("Failed to delete user permanently", zap.Error(err), zap.Int("user_id", user_id))

		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to delete user permanently",
		}
	}

	s.logger.Debug("Successfully deleted user permanently", zap.Int("user_id", user_id))

	return true, nil
}

func (s *userService) RestoreAllUser() (bool, *response.ErrorResponse) {
	s.logger.Debug("Restoring all users")

	_, err := s.userRepository.RestoreAllUser()

	if err != nil {
		s.logger.Error("Failed to restore all users", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to restore all users: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully restored all users")

	return true, nil
}

func (s *userService) DeleteAllUserPermanent() (bool, *response.ErrorResponse) {
	s.logger.Debug("Permanently deleting all users")

	_, err := s.userRepository.DeleteAllUserPermanent()

	if err != nil {
		s.logger.Error("Failed to permanently delete all users", zap.Error(err))
		return false, &response.ErrorResponse{
			Status:  "error",
			Message: "Failed to permanently delete all users: " + err.Error(),
		}
	}

	s.logger.Debug("Successfully deleted all users permanently")

	return true, nil
}
