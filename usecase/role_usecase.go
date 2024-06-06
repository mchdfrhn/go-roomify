package usecase

import (
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/repository"
)

type RoleUsecase interface { // layer untuk komunikasi | jembatan antar layer
	RegisterRole(newRole model.Role) error
	UpdateRole(newRole model.Role) error
	FindRoleById(id string) (model.Role, error)
	DeleteRoleById(id string) error
	FindAllPaging(page int, size int) ([]model.Role, dto.Paging, error)
}

type roleUsecase struct {
	repo repository.RoleRepository
}

func (r *roleUsecase) FindAllPaging(page int, size int) ([]model.Role, dto.Paging, error) {
	return r.repo.GetListPaging(page, size)
}

func (r *roleUsecase) DeleteRoleById(id string) error {
	return r.repo.DeleteRoleById(id)
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

func NewRoleUsecase(repo repository.RoleRepository) RoleUsecase {
	return &roleUsecase{
		repo: repo,
	}
}
