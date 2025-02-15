package repository

import (
	"context"
	"ecommerce/internal/domain/record"
	"ecommerce/internal/domain/requests"
	recordmapper "ecommerce/internal/mapper/record"
	db "ecommerce/pkg/database/schema"
	"fmt"
)

type roleRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.RoleRecordMapping
}

func NewRoleRepository(db *db.Queries, ctx context.Context, mapping recordmapper.RoleRecordMapping) *roleRepository {
	return &roleRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *roleRepository) FindAllRoles(page int, pageSize int, search string) ([]*record.RoleRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetRolesParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetRoles(r.ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find roles: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToRolesRecordAll(res), totalCount, nil
}

func (r *roleRepository) FindById(id int) (*record.RoleRecord, error) {
	res, err := r.db.GetRole(r.ctx, int32(id))

	if err != nil {
		return nil, fmt.Errorf("failed to find role by id: %w", err)
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleRepository) FindByName(name string) (*record.RoleRecord, error) {
	res, err := r.db.GetRoleByName(r.ctx, name)

	if err != nil {
		return nil, fmt.Errorf("failed to find role by name: %w", err)
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleRepository) FindByUserId(user_id int) ([]*record.RoleRecord, error) {
	res, err := r.db.GetUserRoles(r.ctx, int32(user_id))

	if err != nil {
		return nil, fmt.Errorf("failed to find role by user id: %w", err)
	}

	return r.mapping.ToRolesRecord(res), nil
}

func (r *roleRepository) FindByActiveRole(page int, pageSize int, search string) ([]*record.RoleRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetActiveRolesParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveRoles(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find active roles: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToRolesRecordActive(res), totalCount, nil
}

func (r *roleRepository) FindByTrashedRole(page int, pageSize int, search string) ([]*record.RoleRecord, int, error) {
	offset := (page - 1) * pageSize

	req := db.GetTrashedRolesParams{
		Column1: search,
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedRoles(r.ctx, req)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find trashed roles: %w", err)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToRolesRecordTrashed(res), totalCount, nil
}

func (r *roleRepository) CreateRole(req *requests.CreateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.CreateRole(r.ctx, req.Name)

	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleRepository) UpdateRole(req *requests.UpdateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.UpdateRole(r.ctx, db.UpdateRoleParams{
		RoleID:   int32(req.ID),
		RoleName: req.Name,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleRepository) TrashedRole(id int) (*record.RoleRecord, error) {
	err := r.db.TrashRole(r.ctx, int32(id))

	if err != nil {
		return nil, fmt.Errorf("failed to trash role: %w", err)
	}

	role, err := r.FindById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find role after restore: %w", err)
	}

	return role, nil
}

func (r *roleRepository) RestoreRole(id int) (*record.RoleRecord, error) {
	err := r.db.RestoreRole(r.ctx, int32(id))

	if err != nil {
		return nil, fmt.Errorf("failed to restore role: %w", err)
	}

	role, err := r.FindById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find role after restore: %w", err)
	}

	return role, nil
}

func (r *roleRepository) DeleteRolePermanent(role_id int) (bool, error) {
	err := r.db.DeletePermanentRole(r.ctx, int32(role_id))

	if err != nil {
		return false, fmt.Errorf("failed to delete role: %w", err)
	}

	return true, nil
}

func (r *roleRepository) RestoreAllRole() (bool, error) {
	err := r.db.RestoreAllRoles(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to restore all roles: %w", err)
	}

	return true, nil
}

func (r *roleRepository) DeleteAllRolePermanent() (bool, error) {
	err := r.db.DeleteAllPermanentRoles(r.ctx)

	if err != nil {
		return false, fmt.Errorf("failed to delete all roles permanently: %w", err)
	}

	return true, nil
}
