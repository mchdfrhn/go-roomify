package usecase

import (
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/model/dto/request"
	"go-roomify/repository"

	"github.com/google/uuid"
)

type DivisionUsecase interface {
	GetAllDivisions(page, size int) ([]model.Division, dto.Paging, error)
	GetDivisionById(id string) (model.Division, error)
	CreateDivision(division request.DivisionRequest) (model.Division, error)
	UpdateDivision(division model.Division) error
	DeleteDivision(id string) error
}

type divisionUsecase struct {
	repo repository.DivisionRepository
}

func (u *divisionUsecase) GetAllDivisions(page, size int) ([]model.Division, dto.Paging, error) {
	return u.repo.GetAllDivisi(page, size)
}

func (uc *divisionUsecase) GetDivisionById(id string) (model.Division, error) {
	division, err := uc.repo.GetDivisiById(id)
	if err != nil {
		return model.Division{}, err
	}
	return division, nil
}

func (u *divisionUsecase) CreateDivision(division request.DivisionRequest) (model.Division, error) {
	division.Id = uuid.NewString()

	r_division := model.Division{
		Id:   division.Id,
		Name: division.Name,
	}
	return r_division, u.repo.CreateDivisi(r_division)
}

func (u *divisionUsecase) UpdateDivision(division model.Division) error {
	return u.repo.UpdateDivisi(division)
}

func (u *divisionUsecase) DeleteDivision(id string) error {
	return u.repo.DeleteDivisi(id)
}

func NewDivisionUsecase(repo repository.DivisionRepository) DivisionUsecase {
	return &divisionUsecase{repo: repo}
}
