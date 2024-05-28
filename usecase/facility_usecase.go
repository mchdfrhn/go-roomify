package usecase

import (
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/repository"
)

type FacilityUsecase interface {
	FindAllPagingFacility(page int, size int) ([]model.Facility, dto.Paging, error)
	FindFacilityById(id string) (model.Facility, error)
	FindAllFacility() ([]model.Facility, error)
	InputFacility(newFacility model.Facility) error
	UpdatedFacility(newFacility model.Facility) error
	DeletedFacility(id string) error
}

type facilityUsecase struct {
	repo repository.FacilityRepository
}

func NewFacilityUsecase(repo repository.FacilityRepository) FacilityUsecase {
	return &facilityUsecase{
		repo: repo,
	}
}

func (fu *facilityUsecase) FindAllPagingFacility(page int, size int) ([]model.Facility, dto.Paging, error) {
	return fu.repo.GetPagingFacility(page, size)
}

func (fu *facilityUsecase) FindFacilityById(id string) (model.Facility, error) {
	return fu.repo.GetFacilityById(id)
}

func (fu *facilityUsecase) FindAllFacility() ([]model.Facility, error) {
	return fu.repo.GetListFacility()
}

func (fu *facilityUsecase) InputFacility(newFacility model.Facility) error {
	return fu.repo.InsertFacility(newFacility)
}

func (fu *facilityUsecase) UpdatedFacility(newFacility model.Facility) error {
	return fu.repo.UpdateFacility(newFacility)
}

func (fu *facilityUsecase) DeletedFacility(id string) error {
	return fu.repo.DeleteFacility(id)
}
