package delivery

import (
	"go-roomify/config"
	"go-roomify/delivery/controller"
	"go-roomify/repository"
	"go-roomify/usecase"
	"go-roomify/middleware"
	"go-roomify/utils/common"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	authMiddleware middleware.AuthMiddleware
	// customerUc usecase.CustomerUsecase
	// serviceUc usecase.ServiceUsecase
	// transactionUc usecase.TransactionUsecase
	// userUc usecase.UserUsecase
	engine     *gin.Engine
	host       string
}

func (s *Server) setupControllers() {
	// controller.NewCustomerController(s.customerUc, s.engine, s.authMiddleware).Route()
	// controller.NewServiceController(s.serviceUc, s.engine).Route()
	// controller.NewTransactionController(s.transactionUc, s.engine).Route()
	// controller.NewUserController(s.userUc, s.engine).Route()
}

func (s *Server) Run(){
	s.setupControllers()
	if err := s.engine.Run(s.host); err != nil {
		log.Fatal("Server Error : ", err.Error())
	}
}

func NewServer() *Server {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatal("Config : ", err.Error())
	}

	db, err := config.NewDbConnection(cfg)
	
	if err != nil {
		log.Fatal("DB Connect : ",err.Error())
	}

	// // Customer
	// customerRepo := repository.NewCustomerRepository(db.Conn())
	// customerUc := usecase.NewCustomerUsecase(customerRepo)

	// // Service
	// serviceRepo := repository.NewServiceRepository(db.Conn())
	// serivceUc := usecase.NewServiceUsecase(serviceRepo)

	// // Transaction
	// trxRepo := repository.NewTransactionRepository(db.Conn())
	// trxUc := usecase.NewTransactionUsecase(trxRepo,customerUc,serivceUc)

	// // User
	// jwt_token := common.NewJwtToken(cfg.TokenConfig)
	// auth_middleware := middleware.NewAuthMiddleware(jwt_token)

	// userRepo := repository.NewUserRepository(db.Conn())
	// userUc := usecase.NewUserUsecase(userRepo, jwt_token)

	// Gin Engine
	engine := gin.Default()

	return &Server{
		authMiddleware: auth_middleware,
		// customerUc: customerUc,
		// serviceUc: serivceUc,
		// transactionUc: trxUc,
		// userUc: userUc,
		engine: engine,
		host: ":8085",
	}
}