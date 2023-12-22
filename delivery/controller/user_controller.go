package controller

import (
	"net/http"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userUseCase usecase.UserUseCase
}

func NewUserController(r *gin.Engine, userUC usecase.UserUseCase) {
	userController := &UserController{
		userUseCase: userUC,
	}

	r.POST("/users", middleware.JWTAuthMiddleware("owner", "developer"), userController.CreateUser)
	r.PUT("/users/:id", middleware.JWTAuthMiddleware("owner", "developer"), userController.UpdateUser)
	r.GET("/users/:username", middleware.JWTAuthMiddleware("owner", "developer"), userController.GetUserByUsername)
	r.GET("/users", middleware.JWTAuthMiddleware("owner", "developer"), userController.GetAllUsers)
	r.DELETE("/users/:username", middleware.JWTAuthMiddleware("owner", "developer"), userController.DeleteUser)
}
func (uc *UserController) CreateUser(c *gin.Context) {
	var user model.User
	username, err := utils.GetUsernameFromContext(c)
	if err := c.ShouldBindJSON(&user); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusBadRequest, "Bad request", nil)
		return
	}

	logrus.Infof("[%v] Created user", username)

	if user.Role != "owner" && user.Role != "developer" && user.Role != "admin" && user.Role != "employee" {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid role, must be owner, developer, admin or employee", nil)
		return

	}

	user.ID = uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to encrypt password", nil)
		return
	}

	user.Password = string(hashedPassword)

	err = uc.userUseCase.CreateUser(&user)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		if err.Error() == "username already exists" {
			utils.SendResponse(c, http.StatusConflict, "username already exists", nil)
			return
		}

		utils.SendResponse(c, http.StatusInternalServerError, "internal error", nil)
		return
	}
	logrus.Infof("[%v] Created user %v", username, user.Username)
	utils.SendResponse(c, http.StatusOK, "Success", user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("username")
	username, err := utils.GetUsernameFromContext(c)
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusBadRequest, "Bad request", nil)
		return
	}
	user.ID = userID

	token, err := utils.ExtractTokenFromAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid authorization header", nil)
		return
	}

	claims, err := utils.VerifyJWTToken(token)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	userName := claims["username"].(string)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to encrypt password", nil)
		return
	}

	user.Password = string(hashedPassword)
	user.IsActive = true
	if err := uc.userUseCase.UpdateUser(&user, userName); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	logrus.Infof("[%v] Updated user %v", username, user.Username)
	utils.SendResponse(c, http.StatusOK, "Success", user)
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	username, err := utils.GetUsernameFromContext(c)
	user, err := uc.userUseCase.GetUserByID(userID)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user", nil)
		return
	}
	if user == nil {
		logrus.Errorf("[%v] User not found", username)
		utils.SendResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}
	logrus.Infof("[%v] Get user %v", username, user.Username)
	utils.SendResponse(c, http.StatusOK, "Success", user)
}

func (uc *UserController) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := uc.userUseCase.GetUserByUsername(username)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user", nil)
		return
	}
	if user == nil {
		logrus.Errorf("[%v] User not found", username)
		utils.SendResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	logrus.Infof("[%v] Get user %v", username, user.Username)
	utils.SendResponse(c, http.StatusOK, "Success", user)
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	users, err := uc.userUseCase.GetAllUsers()
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	logrus.Infof("[%v] Get all users", username)
	utils.SendResponse(c, http.StatusOK, "Success", users)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	username := c.Param("username")

	if err := uc.userUseCase.DeleteUser(username); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete user", nil)
		return
	}
	logrus.Infof("[%v] Deleted user %v", username, username)
	utils.SendResponse(c, http.StatusOK, "Success", nil)
}
