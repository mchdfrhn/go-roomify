package usecase

import (
	"errors"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/repository"
	"go-roomify/utils/common"
	"go-roomify/utils/encryption"

	// "fmt"

	"github.com/google/uuid"
)

type UserCredentialUsecase interface {
	LoginUser(paylod request.UserCredentialRequest) (response.UserCredentialResponse, error)
	GetListUser(page int, size int) ([]model.UserCredential, dto.Paging, error)
	GetUserById(id string) (model.UserCredential, error)
	CreateNewUser(new_user model.UserCredential) (model.UserCredential, error)
	UpdatePassword(new_user request.UserUpdatePasswordRequest) (error)
	DeleteUser(id string) (error)
}

type userCredentialUsecase struct {
	repo repository.UserCredentialRepository
	jwtToken common.JwtToken
}

func (self *userCredentialUsecase) LoginUser(payload request.UserCredentialRequest) (response.UserCredentialResponse, error) {
	r_user_cr, err := self.repo.GetByUsername(payload.Username)

	if err != nil {
		return response.UserCredentialResponse{}, err
	}

	// fmt.Println("payload: ", payload.Password, r_user_cr.Password)

	isValid := encryption.Sha512ComparePassword(payload.Password, r_user_cr.Password)

	if !isValid {
		return response.UserCredentialResponse{}, errors.New("Password is invalid")
	}

	access_token, err := self.jwtToken.GenerateTokenJwt(r_user_cr)

	if err != nil {
		return response.UserCredentialResponse{}, err
	}

	return response.UserCredentialResponse{
		AccessToken: access_token,
		UserId: r_user_cr.Id,
	}, nil
}

func (self *userCredentialUsecase) CreateNewUser(new_user model.UserCredential) (model.UserCredential, error) {
	new_user.Id = uuid.NewString()
	new_user.Password = encryption.Sha512HashPassword(new_user.Password)

	r_user, err := self.repo.AddNew(new_user)

	if err != nil {
		return model.UserCredential{}, err
	}

	return r_user, nil
}

func (self *userCredentialUsecase) GetListUser(page int, size int) ([]model.UserCredential, dto.Paging, error) {
	return self.repo.GetList(page, size)
}

func (self *userCredentialUsecase) GetUserById(id string) (model.UserCredential, error) {
	return self.repo.GetById(id)
}

func (self *userCredentialUsecase) UpdatePassword(new_user request.UserUpdatePasswordRequest) (error) {
	new_user.OldPassword = encryption.Sha512HashPassword(new_user.OldPassword)
	new_user.NewPassword = encryption.Sha512HashPassword(new_user.NewPassword)

	r_user, err := self.GetUserById(new_user.Id)

	if err != nil {
		return err
	}

	if r_user.Password != new_user.OldPassword {
		return errors.New("Password Doesn't Match")
	}

	return self.repo.UpdatePassword(new_user)
}

func (self *userCredentialUsecase) DeleteUser(id string) (error) {
	return self.repo.Delete(id)
}

func NewUserCredentialUsecase(repo repository.UserCredentialRepository, jwt_token common.JwtToken) UserCredentialUsecase {
	return &userCredentialUsecase {
		repo : repo,
		jwtToken: jwt_token,
	}
}
