package delivery

import (
	"io/ioutil"
	"os"
	"trackprosto/config"
	"trackprosto/delivery/controller"
	"trackprosto/manager"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	useCaseManager manager.UsecaseManager
	engine         *gin.Engine
}

func (s *Server) Run() {
	s.initController()
	err := s.engine.Run()
	if err != nil {
		panic(err)
	}
}

func (s *Server) initController() {
	controller.NewUserController(s.engine, s.useCaseManager.GetUserUsecase())
	controller.NewLoginController(s.engine, s.useCaseManager.GetLoginUsecase())
	controller.NewMeatController(s.engine, s.useCaseManager.GetMeatUsecase())
	controller.NewTransactionController(s.engine, s.useCaseManager.GetTransactionUseCase())
	controller.NewCreditPaymentController(s.engine, s.useCaseManager.GetCreditPaymentUseCase())
	controller.NewCustomerController(s.engine, s.useCaseManager.GetCustomerUsecase())
	controller.NewCompanyController(s.engine, s.useCaseManager.GetCompanyUsecase())
}

func NewServer() *Server {
	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true
	configCors.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	configCors.AllowHeaders = []string{"Origin", "v-Length", "Content-Type", "Authorization"}

	r.Use(cors.New(configCors))

	infra := manager.NewInfraManager(c)
	repo := manager.NewRepoManager(infra)
	usecase := manager.NewUsecaseManager(repo)

	// Inisialisasi logger
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)        // Atur output ke stdout atau file log
	gin.SetMode(gin.ReleaseMode)       // Mengubah mode Gin menjadi "release" di lingkungan produksi
	gin.DefaultWriter = ioutil.Discard // Menyembunyikan log bawaan Gin

	// Contoh pengaturan level logging
	logger.SetLevel(log.DebugLevel)

	return &Server{useCaseManager: usecase, engine: r}
}
