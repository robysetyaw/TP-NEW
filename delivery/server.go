package delivery

import (
	"trackprosto/config"
	"trackprosto/delivery/controller"
	"trackprosto/manager"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	return &Server{useCaseManager: usecase, engine: r}
}
