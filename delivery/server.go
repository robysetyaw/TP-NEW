package delivery

import (
	// "io"
	// "io/ioutil"
	// "net"
	"os"
	// "path"
	// "time"
	"trackprosto/config"
	"trackprosto/delivery/controller"
	"trackprosto/manager"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	// "gopkg.in/olivere/elastic.v5"
	// "github.com/bshuster-repo/logrus-logstash-hook"
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
	controller.NewDailyExpenditureController(s.engine, s.useCaseManager.GetDailyExpenditureUseCase())
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
	logger := logrus.New()
	logger.SetOutput(os.Stdout) // Atur output ke stdout atau file log
	// logPath := "log/"
	// if _, err := os.Stat(logPath); os.IsNotExist(err) {
	// 	// Buat direktori log jika belum ada
	// 	err := os.Mkdir(logPath, os.ModePerm)
	// 	if err != nil {
	// 		logrus.Errorf("Failed to create log directory: %v", err)
	// 		logger.SetOutput(ioutil.Discard) // Jangan menulis log jika gagal membuat direktori
	// 	}
	// }

	// logFile := time.Now().Format("2006-01-02") + ".log"
	// logFullPath := path.Join(logPath, logFile)

	// file, err := os.OpenFile(logFullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	// 	logger.SetOutput(io.MultiWriter(file, os.Stdout))
	// } else {
	// 	logrus.Errorf("Failed to open log file: %v", err)
	// 	// logger.SetOutput(ioutil.Discard) // Jangan menulis log jika gagal membuka file
	// }

	// Atur level logging
	logger.SetLevel(logrus.DebugLevel)

	// Konfigurasi hook Elasticsearch
	// conn, err := net.Dial("tcp", "no-company-search-3575562732.us-east-1.bonsaisearch.net:443") // Ganti dengan URL dan port Logstash Anda
	// if err == nil {
	//     esHook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{
	// 		"app": "trackprosto",
	// 		"env": "production",
	// 		"ver": "1.0.0",
	// 	}))
	//     logger.Hooks.Add(esHook)
	// } else {
	//     logrus.Errorf("Failed to create Logstash hook: %v", err)
	// }

	// Set Logrus logger sebagai default logger
	logrus.SetOutput(logger.Out)
	logrus.SetLevel(logger.Level)

	gin.SetMode(gin.ReleaseMode) // Mengubah mode Gin menjadi "release" di lingkungan produksi
	// gin.DefaultWriter = ioutil.Discard // Menyembunyikan log bawaan Gin

	return &Server{useCaseManager: usecase, engine: r}
}
