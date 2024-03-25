package services

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/hutamatr/GoBlogify/helpers"
	"github.com/hutamatr/GoBlogify/model/domain"
	"github.com/hutamatr/GoBlogify/model/web"
	repositoriesRole "github.com/hutamatr/GoBlogify/repositories/role"
)

type RoleServiceImpl struct {
	repository repositoriesRole.RoleRepository
	db         *sql.DB
	validator  *validator.Validate
}

func NewRoleService(roleRepository repositoriesRole.RoleRepository, db *sql.DB, validator *validator.Validate) RoleService {
	return &RoleServiceImpl{
		repository: roleRepository,
		db:         db,
		validator:  validator,
	}
}

func (service *RoleServiceImpl) Create(ctx context.Context, request web.RoleCreateRequest) web.RoleResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	roleRequest := domain.Role{
		Name: request.Name,
	}

	createdRole := service.repository.Save(ctx, tx, roleRequest)

	return web.ToRoleResponse(createdRole)
}

func (service *RoleServiceImpl) FindAll(ctx context.Context) []web.RoleResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	roles := service.repository.FindAll(ctx, tx)

	var rolesData []web.RoleResponse

	for _, role := range roles {
		rolesData = append(rolesData, web.ToRoleResponse(role))
	}

	return rolesData
}

func (service *RoleServiceImpl) FindById(ctx context.Context, roleId int) web.RoleResponse {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	role := service.repository.FindById(ctx, tx, roleId)

	return web.ToRoleResponse(role)
}

func (service *RoleServiceImpl) Update(ctx context.Context, request web.RoleUpdateRequest) web.RoleResponse {
	err := service.validator.Struct(request)
	helpers.PanicError(err)

	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	role := service.repository.FindById(ctx, tx, request.Id)

	role.Name = request.Name

	updatedRole := service.repository.Update(ctx, tx, role)

	return web.ToRoleResponse(updatedRole)
}

func (service *RoleServiceImpl) Delete(ctx context.Context, roleId int) {
	tx, err := service.db.Begin()
	helpers.PanicError(err)
	defer helpers.TxRollbackCommit(tx)

	service.repository.Delete(ctx, tx, roleId)
}