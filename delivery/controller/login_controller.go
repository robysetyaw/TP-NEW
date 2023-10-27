package controller

import (
	"net/http"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
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
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		utils.SendResponse(c, http.StatusBadRequest, "Invalid login data", nil)
		return
	}

	if loginData.Username == "" || loginData.Password == "" {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		utils.SendResponse(c, http.StatusBadRequest, "Invalid username or password", nil)
		return
	}

	token, err := uc.loginUseCase.Login(loginData.Username, loginData.Password)
	if err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid username or password", nil)
		return
	}

	// c.JSON(http.StatusOK, gin.H{"token ": token})
	utils.SendResponse(c, http.StatusOK, "Login success", token)
}
