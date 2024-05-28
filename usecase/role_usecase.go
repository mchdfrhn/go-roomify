package usecase

import (
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/repository"

	"github.com/gin-gonic/gin"
)

type RoleUsecase interface { // layer untuk komunikasi | jembatan antar layer
	RegisterRole(newRole model.Role) error
	FindAllRole() ([]model.Role, error)
	UpdateRole(newRole model.Role) error
	FindRoleById(id string) (model.Role, error)
	DeleteRoleById(cxt *gin.Context, id string) error
	FindAllPaging(page int, size int) ([]model.Role, dto.Paging, error)
}

type roleUsecase struct {
	repo repository.RoleRepository
}

func (r *roleUsecase) FindAllPaging(page int, size int) ([]model.Role, dto.Paging, error) {
	return r.repo.GetListPaging(page, size)
}

func (r *roleUsecase) DeleteRoleById(ctx *gin.Context, id string) error {
	return r.repo.DeleteRoleById(ctx, id)
}

func (r *roleUsecase) FindRoleById(id string) (model.Role, error) {
	return r.repo.GetRoleById(id)
}

func (r *roleUsecase) UpdateRole(newRole model.Role) error {
	return r.repo.UpdateRole(newRole)
}

func (r *roleUsecase) RegisterRole(newRole model.Role) error {
	return r.repo.InsertRole(newRole)
}

func (r *roleUsecase) FindAllRole() ([]model.Role, error) {
	return r.repo.GetListRole()
}

func NewRoleUsecase(repo repository.RoleRepository) RoleUsecase {
	return &roleUsecase{
		repo: repo,
	}
}
