package usecase

import (
	"errors"
	//"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/repository"
	"go-roomify/utils/common"
	"go-roomify/utils/encryption"

	// "fmt"

	// "github.com/google/uuid"
)

type UserCredentialUsecase interface {
	LoginUser(paylod request.UserCredentialRequest) (response.UserCredentialResponse,error)
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

	isValid := encryption.ComparePassword(
		payload.Password,
		r_user_cr.Password)

	if !isValid {
		return response.UserCredentialResponse{}, errors.New("Password is invalid")
	}

	access_token, err := self.jwtToken.GenerateTokenJwt(r_user_cr)

	if err != nil {
		return response.UserCredentialResponse{},err
	}

	return response.UserCredentialResponse{
		AccessToken: access_token,
		UserId: r_user_cr.Id,
	}, nil
}

func NewUserCredentialUsecase(repo repository.UserCredentialRepository, jwt_token common.JwtToken) UserCredentialUsecase {
	return &userCredentialUsecase {
		repo : repo,
		jwtToken: jwt_token,
	}
}
