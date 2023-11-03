package controller

import (
	"net/http"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MeatController struct {
	meatUseCase usecase.MeatUseCase
}

func NewMeatController(r *gin.Engine, meatUC usecase.MeatUseCase) {
	meatController := &MeatController{
		meatUseCase: meatUC,
	}

	r.POST("/meats", middleware.JWTAuthMiddleware(), meatController.CreateMeat)
	r.GET("/meats", middleware.JWTAuthMiddleware(), meatController.GetAllMeats)
	r.GET("/meats/:name", middleware.JWTAuthMiddleware(), meatController.GetMeatByName)
	r.DELETE("/meats/:id", middleware.JWTAuthMiddleware(), meatController.DeleteMeat)
	r.PUT("/meats/:id", middleware.JWTAuthMiddleware(), meatController.UpdateMeat)
}

func (mc *MeatController) CreateMeat(ctx *gin.Context) {
	var meat model.Meat
	if err := ctx.ShouldBindJSON(&meat); err != nil {
		// ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		utils.SendResponse(ctx, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	token, err := utils.ExtractTokenFromAuthHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		// ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		utils.SendResponse(ctx, http.StatusUnauthorized, "Invalid authorization header", nil)
		return
	}

	claims, err := utils.VerifyJWTToken(token)
	if err != nil {
		// ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		utils.SendResponse(ctx, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	userName := claims["username"].(string)
	meat.CreatedBy = userName
	meat.ID = uuid.New().String()
	err = mc.meatUseCase.CreateMeat(&meat)
	if err != nil {
		if err == utils.ErrMeatNameAlreadyExist {
			utils.SendResponse(ctx, http.StatusConflict, "meatname already exists", nil)
		} else {
			utils.SendResponse(ctx, http.StatusInternalServerError, "Failed to create meat", nil)
		}
		// ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meat"})

		return
	}
	// ctx.JSON(http.StatusCreated, meat)
	utils.SendResponse(ctx, http.StatusCreated, "Success", meat)
}

func (mc *MeatController) GetAllMeats(c *gin.Context) {
	meats, err := mc.meatUseCase.GetAllMeats()
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get meats"})
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get meats", nil)
		return
	}

	// c.JSON(http.StatusOK, meats)
	utils.SendResponse(c, http.StatusOK, "Success", meats)
}

func (mc *MeatController) GetMeatByName(c *gin.Context) {
	name := c.Param("name")
	meat, err := mc.meatUseCase.GetMeatByName(name)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get meat"})
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get meat", nil)
		return
	}
	if meat == nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Meat not found"})
		utils.SendResponse(c, http.StatusNotFound, "Meat not found", nil)
		return
	}

	// c.JSON(http.StatusOK, meat)
	utils.SendResponse(c, http.StatusOK, "Success", meat)
}

func (mc *MeatController) GetMeatById(c *gin.Context) {
	id := c.Param("id")
	meat, err := mc.meatUseCase.GetMeatByName(id)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get meat"})
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get meat", nil)
		return
	}
	if meat == nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Meat not found"})
		utils.SendResponse(c, http.StatusNotFound, "Meat not found", nil)
		return
	}

	// c.JSON(http.StatusOK, meat)
	utils.SendResponse(c, http.StatusOK, "Success", meat)
}

func (uc *MeatController) DeleteMeat(c *gin.Context) {
	meatID := c.Param("id")

	if err := uc.meatUseCase.DeleteMeat(meatID); err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meat"})
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete meat", nil)
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "Meat deleted successfully"})
	utils.SendResponse(c, http.StatusOK, "Success", nil)
}

func (uc *MeatController) UpdateMeat(ctx *gin.Context) {
	meatID := ctx.Param("id")

	var meat model.Meat
	if err := ctx.ShouldBindJSON(&meat); err != nil {
		// ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	token, err := utils.ExtractTokenFromAuthHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		// ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		return
	}

	claims, err := utils.VerifyJWTToken(token)
	if err != nil {
		// ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userName := claims["username"].(string)
	meat.UpdatedBy = userName
	meat.ID = meatID

	if err := uc.meatUseCase.UpdateMeat(&meat); err != nil {
		// ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ctx.JSON(http.StatusOK, meat)
	ctx.JSON(http.StatusOK, gin.H{"message": "Meat updated successfully"})
}
