package controller

import (
	//"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserCredentialController struct {
	uc usecase.UserCredentialUsecase
	rg *gin.RouterGroup
}

func (self *UserCredentialController) loginHandler(ctx *gin.Context) {
	var login_payload request.UserCredentialRequest

	if err := ctx.ShouldBindJSON(&login_payload); err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	r_user_cr, err := self.uc.LoginUser(login_payload)

	if err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	response.SendSingleResponseCreated(
		ctx,
		r_user_cr,
		fmt.Sprintf("Success Login"),
	)
}

func (self *UserCredentialController) Route(){
	router := self.rg.Group("/auth")
	router.POST("/login", self.loginHandler)	
}

func NewUserCredentialController(uc usecase.UserCredentialUsecase, router *gin.Engine) *UserCredentialController {
	return &UserCredentialController{
		uc: uc,
		rg : &router.RouterGroup,
	}
}
