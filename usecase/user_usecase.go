package usecase

import (
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/model/dto/request"
	"go-roomify/repository"
)

type UserProfileUsecase interface {
	GetList(page, size int) ([]model.UserProfile, dto.Paging, error)
	GetById(id string) (model.UserProfile, error)
	GetByUsername(username string) (model.UserProfile, error)
	Create(newUser request.UserProfileRequest) (request.UserProfileRequest, error)
	Update(newUser request.UserProfileRequest) (request.UserProfileRequest, error)
	Delete(id string) error
}

type userProfileUsecase struct {
	repo repository.UserProfileRepository
}

func (u *userProfileUsecase) GetList(page, size int) ([]model.UserProfile, dto.Paging, error) {
	return u.repo.GetList(page, size)
}

func (u *userProfileUsecase) GetById(id string) (model.UserProfile, error) {
	return u.repo.GetById(id)
}

func (u *userProfileUsecase) GetByUsername(username string) (model.UserProfile, error) {
	return u.repo.GetByUsername(username)
}

func (u *userProfileUsecase) Create(newUser request.UserProfileRequest) (request.UserProfileRequest, error) {
	
	return u.repo.Create(newUser)
}

func (u *userProfileUsecase) Update(newUser request.UserProfileRequest) (request.UserProfileRequest, error) {
	return u.repo.Update(newUser)
}

func (u *userProfileUsecase) Delete(id string) error {
	return u.repo.Delete(id)
}

func NewUserUsecase(repo repository.UserProfileRepository) UserProfileUsecase {
	return &userProfileUsecase{
		repo: repo,
	}
}
