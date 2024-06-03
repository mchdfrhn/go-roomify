package delivery

import (
	"go-roomify/config"
	"go-roomify/delivery/controller"
	"go-roomify/middleware"
	"go-roomify/repository"
	"go-roomify/usecase"
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
	roomUc     usecase.RoomUsecase
	facilityUc usecase.FacilityUsecase

	userUc           usecase.UserProfileUsecase
	roleUc           usecase.RoleUsecase
	divisionUc       usecase.DivisionUsecase
	userCredentialUc usecase.UserCredentialUsecase

	reservationUc usecase.ReservationUsecase
	reportUc      usecase.ReportUsecase

	engine *gin.Engine
	host   string
}

func (self *Server) setupControllers() {
	// controller.NewCustomerController(s.customerUc, s.engine, s.authMiddleware).Route()
	// controller.NewServiceController(s.serviceUc, s.engine).Route()
	// controller.NewTransactionController(s.transactionUc, s.engine).Route()
	controller.NewRoomController(self.roomUc, self.engine, self.authMiddleware).Route()
	controller.NewFacilityController(self.facilityUc, self.engine, self.authMiddleware).Route()

	controller.NewUserController(self.userUc, self.engine, self.authMiddleware).Route()
	controller.NewRoleController(self.roleUc, self.engine, self.authMiddleware).Route()
	controller.NewDivisionController(self.divisionUc, self.engine, self.authMiddleware).Route()
	controller.NewUserCredentialController(self.userCredentialUc, self.engine, self.authMiddleware).Route()

	controller.NewReservationController(self.reservationUc, self.engine, self.authMiddleware).Route()
	controller.NewReportController(self.reportUc, self.engine, self.authMiddleware).Route()
}

func (s *Server) Run() {
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
		log.Fatal("DB Connect : ", err.Error())
	}
	// User
	jwt_token := common.NewJwtToken(cfg.TokenConfig)
	auth_middleware := middleware.NewAuthMiddleware(jwt_token)

	// // Customer
	// customerRepo := repository.NewCustomerRepository(db.Conn())
	// customerUc := usecase.NewCustomerUsecase(customerRepo)

	// // Service
	// serviceRepo := repository.NewServiceRepository(db.Conn())
	// serivceUc := usecase.NewServiceUsecase(serviceRepo)

	roomRepo := repository.NewRoomRepository(db.Conn())
	roomUc := usecase.NewRoomUsecase(roomRepo)

	facilityRepo := repository.NewFacilityRepository(db.Conn())
	facilityUc := usecase.NewFacilityUsecase(facilityRepo)

	// ---- USER -----------

	userRepo := repository.NewUserProfileRepository(db.Conn())
	userUc := usecase.NewUserUsecase(userRepo)

	roleRepo := repository.NewRoleRepository(db.Conn())
	roleUc := usecase.NewRoleUsecase(roleRepo)

	divisionRepo := repository.NewDivisionRepository(db.Conn())
	divisionUc := usecase.NewDivisionUsecase(divisionRepo)

	user_credential_repo := repository.NewUserCredentialRepository(db.Conn())
	user_credential_uc := usecase.NewUserCredentialUsecase(user_credential_repo, jwt_token)

	// // Transaction
	resevationRepo := repository.NewReservationRepository(db.Conn())
	resevationUc := usecase.NewReservationUsecase(resevationRepo)

	// Report
	//reportRepo := repository.NewReportRepository(db.Conn())
	reportUc := usecase.NewReportUsecase(resevationRepo)

	// Gin Engine
	engine := gin.Default()

	return &Server{
		authMiddleware: auth_middleware,
		// customerUc: customerUc,
		// serviceUc: serivceUc,
		// transactionUc: trxUc,
		facilityUc: facilityUc,
		roomUc:     roomUc,

		userUc:           userUc,
		divisionUc:       divisionUc,
		roleUc:           roleUc,
		userCredentialUc: user_credential_uc,

		reservationUc: resevationUc,
		reportUc:      reportUc,

		engine: engine,
		host:   ":8085",
	}
}
