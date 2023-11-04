package controller

import (
	"net/http"
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
}

func (uc *LoginController) Login(c *gin.Context) {
	var loginData model.LoginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid login data", nil)
		return
	}

	if loginData.Username == "" || loginData.Password == "" {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid username or password", nil)
		return
	}
	logrus.Infof("User [%s] is logging in", loginData.Username)
	token, err := uc.loginUseCase.Login(loginData.Username, loginData.Password)
	if err != nil {
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid username or password", nil)
		return
	}
	logrus.Infof("User [%s] logged in successfully", loginData.Username)
	utils.SendResponse(c, http.StatusOK, "Login success", token)
}
