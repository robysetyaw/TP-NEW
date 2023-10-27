package controller

import (
	"net/http"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userUseCase usecase.UserUseCase
}

func NewUserController(r *gin.Engine, userUC usecase.UserUseCase) {
	userController := &UserController{
		userUseCase: userUC,
	}

	r.POST("/users", userController.CreateUser)
	r.PUT("/users/:id", userController.UpdateUser)
	// r.GET("/users/:id", userController.GetUserByID)
	r.GET("/users/:username", middleware.JWTAuthMiddleware(), userController.GetUserByUsername)
	r.GET("/users", middleware.JWTAuthMiddleware(), userController.GetAllUsers)
	r.DELETE("/users/:username", middleware.JWTAuthMiddleware(), userController.DeleteUser)
}
func (uc *UserController) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusBadRequest, "Bad request", nil)
		return
	}

	user.ID = uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to encrypt password", nil)
		return
	}

	user.Password = string(hashedPassword)

	// if err := uc.userUseCase.CreateUser(&user);
	err = uc.userUseCase.CreateUser(&user)
	if err != nil {
		if err.Error() == "username already exists" {
			utils.SendResponse(c, http.StatusConflict, "username already exists", nil)
			return
		}
		utils.SendResponse(c, http.StatusInternalServerError, "internal error", nil)
		return
	}

	// c.JSON(http.StatusOK, user)
	utils.SendResponse(c, http.StatusOK, "Success", user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("username")

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusBadRequest, "Bad request", nil)
		return
	}
	user.ID = userID

	token, err := utils.ExtractTokenFromAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid authorization header", nil)
		return
	}

	claims, err := utils.VerifyJWTToken(token)
	if err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	userName := claims["username"].(string)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to encrypt password", nil)
		return
	}

	user.Password = string(hashedPassword)
	user.IsActive = true
	if err := uc.userUseCase.UpdateUser(&user, userName); err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, user)
	utils.SendResponse(c, http.StatusOK, "Success", user)
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := uc.userUseCase.GetUserByID(userID)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user", nil)
		return
	}
	if user == nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		utils.SendResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	// c.JSON(http.StatusOK, user)
	utils.SendResponse(c, http.StatusOK, "Success", user)
}

func (uc *UserController) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := uc.userUseCase.GetUserByUsername(username)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user", nil)
		return
	}
	if user == nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		utils.SendResponse(c, http.StatusNotFound, "User not found", nil)
		return
	}

	// c.JSON(http.StatusOK, user)
	utils.SendResponse(c, http.StatusOK, "Success", user)
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userUseCase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	// c.JSON(http.StatusOK, users)
	utils.SendResponse(c, http.StatusOK, "Success", users)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	username := c.Param("username")

	if err := uc.userUseCase.DeleteUser(username); err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete user", nil)
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	utils.SendResponse(c, http.StatusOK, "Success", nil)
}
