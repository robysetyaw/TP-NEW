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
	username, err := utils.GetUsernameFromContext(ctx)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(ctx, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is creating a meat", username)
	var meat model.Meat
	if err := ctx.ShouldBindJSON(&meat); err != nil {
		logrus.Error(err)
		utils.SendResponse(ctx, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	meat.CreatedBy = username
	meat.ID = uuid.New().String()
	err = mc.meatUseCase.CreateMeat(&meat)
	if err != nil {
		if err == utils.ErrMeatNameAlreadyExist {
			logrus.Error(err)
			utils.SendResponse(ctx, http.StatusConflict, "meatname already exists", nil)
		} else {
			logrus.Error(err)
			utils.SendResponse(ctx, http.StatusInternalServerError, "Failed to create meat", nil)
		}
		return
	}
	logrus.Info("Meat created successfully, meatname ", meat.Name)
	utils.SendResponse(ctx, http.StatusCreated, "Success", meat)
}

func (mc *MeatController) GetAllMeats(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Info("[", username, "] get all meats")
	meats, err := mc.meatUseCase.GetAllMeats()
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get meats", nil)
		return
	}
	logrus.Info("Success get all meats", meats)
	utils.SendResponse(c, http.StatusOK, "Success", meats)
}

func (mc *MeatController) GetMeatByName(c *gin.Context) {
	name := c.Param("name")
	meat, err := mc.meatUseCase.GetMeatByName(name)
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get meat", nil)
		return
	}
	if meat == nil {
		utils.SendResponse(c, http.StatusNotFound, "Meat not found", nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Success", meat)
}

func (mc *MeatController) GetMeatById(c *gin.Context) {
	id := c.Param("id")
	meat, err := mc.meatUseCase.GetMeatByName(id)
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get meat", nil)
		return
	}
	if meat == nil {
		utils.SendResponse(c, http.StatusNotFound, "Meat not found", nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Success", meat)
}

func (uc *MeatController) DeleteMeat(c *gin.Context) {
	meatID := c.Param("id")

	if err := uc.meatUseCase.DeleteMeat(meatID); err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete meat", nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Success", nil)
}

func (uc *MeatController) UpdateMeat(ctx *gin.Context) {
	meatID := ctx.Param("id")

	var meat model.Meat
	if err := ctx.ShouldBindJSON(&meat); err != nil {
		logrus.Error(err)
		utils.SendResponse(ctx, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}
	userName, err := utils.GetUsernameFromContext(ctx)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(ctx, http.StatusUnauthorized, "Invalid token", nil)
		return

	}
	meat.UpdatedBy = userName
	meat.ID = meatID
	logrus.Infof("[%s] is updating meat [%s]", userName, meatID)
	if err := uc.meatUseCase.UpdateMeat(&meat); err != nil {
		logrus.Error(err)
		utils.SendResponse(ctx, http.StatusInternalServerError, "Failed to update meat", nil)
		return
	}

	logrus.Info("Meat updated successfully", meat)
	ctx.JSON(http.StatusOK, gin.H{"message": "Meat updated successfully"})
}
