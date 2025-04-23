package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"fmt"
)

type userRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.UserRecordMapping
}

func NewUserRepository(db *db.Queries, ctx context.Context, mapping recordmapper.UserRecordMapping) *userRepository {
	return &userRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *userRepository) FindAllUsers(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsers(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find all users: %w", err)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToUsersRecordPagination(res), &totalCount, nil
}

func (r *userRepository) FindById(user_id int) (*record.UserRecord, error) {
	res, err := r.db.GetUserByID(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find id users: %w", err)
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) FindByActive(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsersActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find active users: %w", err)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToUsersRecordActivePagination(res), &totalCount, nil
}

func (r *userRepository) FindByTrashed(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUserTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUserTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to find users: %w", err)
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToUsersRecordTrashedPagination(res), &totalCount, nil
}

func (r *userRepository) FindByEmail(email string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByEmail(r.ctx, email)

	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) CreateUser(request *requests.CreateUserRequest) (*record.UserRecord, error) {
	req := db.CreateUserParams{
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	user, err := r.db.CreateUser(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed create user: %w", err)
	}

	return r.mapping.ToUserRecord(user), nil
}

func (r *userRepository) UpdateUser(request *requests.UpdateUserRequest) (*record.UserRecord, error) {
	req := db.UpdateUserParams{
		UserID:    int32(*request.UserID),
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	res, err := r.db.UpdateUser(r.ctx, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) TrashedUser(user_id int) (*record.UserRecord, error) {
	res, err := r.db.TrashUser(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash user: %w", err)
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) RestoreUser(user_id int) (*record.UserRecord, error) {
	res, err := r.db.RestoreUser(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore topup: %w", err)
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) DeleteUserPermanent(user_id int) (bool, error) {
	err := r.db.DeleteUserPermanently(r.ctx, int32(user_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}

	return true, nil
}

func (r *userRepository) RestoreAllUser() (bool, error) {
	err := r.db.RestoreAllUsers(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all users: %w", err)
	}
	return true, nil
}

func (r *userRepository) DeleteAllUserPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentUsers(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all users permanently: %w", err)
	}
	return true, nil
}
