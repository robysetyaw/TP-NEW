package controller

import (
	"net/http"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LoginController struct {
	loginUseCase usecase.LoginUseCase
}

func NewLoginController(r *gin.Engine, loginUC usecase.LoginUseCase) {
	loginController := &LoginController{
		loginUseCase: loginUC,
	}
	r.POST("/login", loginController.Login)
	r.POST("/send-log", middleware.JWTAuthMiddleware("employee", "admin", "owner", "developer"), loginController.SendLog)
}

func (uc *LoginController) Login(c *gin.Context) {
	var loginData model.LoginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusBadRequest, "Invalid login data", nil)
		return
	}

	if loginData.Username == "" || loginData.Password == "" {
		logrus.Errorf("Invalid username or password for [%s]", loginData.Username)
		utils.SendResponse(c, http.StatusBadRequest, "Invalid username or password", nil)
		return
	}
	logrus.Infof("[%s] is logging in", loginData.Username)
	token, err := uc.loginUseCase.Login(loginData.Username, loginData.Password)
	if err != nil {
		logrus.Error(err)
		utils.HandleError(c, err)
		return
	}
	logrus.Infof("[%s] logged in successfully", loginData.Username)
	utils.SendResponse(c, http.StatusOK, "Login success", token)
}

func (uc *LoginController) SendLog(c *gin.Context) {
	// Mendapatkan log dari body request
	var logRequest struct {
		Log string `json:"log"`
	}

	if err := c.ShouldBindJSON(&logRequest); err != nil {
		logrus.Error("Invalid log format")
		utils.SendResponse(c, http.StatusBadRequest, "Invalid log format", nil)
		return
	}
	username, role, err := utils.GetUserDetailsFromContext(c)
	if err != nil {
		logrus.Errorf("[%s] %v", username, err)
		utils.HandleError(c, err)
		return
	}
	// Memasukkan log ke dalam logrus
	logrus.Infof("[Frontend Log][%s][%s] %s", username, role, logRequest.Log)
	utils.SendResponse(c, http.StatusOK, "Log created", nil)
}
