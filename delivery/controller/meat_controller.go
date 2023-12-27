package controller

import (
	"net/http"
	"strconv"
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

	r.POST("/meats", middleware.JWTAuthMiddleware("admin", "owner", "developer"), meatController.CreateMeat)
	r.GET("/meats", middleware.JWTAuthMiddleware("employee", "admin", "owner", "developer"), meatController.GetAllMeats)
	r.GET("/meats/:name", middleware.JWTAuthMiddleware("employee", "admin", "owner", "developer"), meatController.GetMeatByName)
	r.DELETE("/meats/:id", middleware.JWTAuthMiddleware("admin", "owner", "developer"), meatController.DeleteMeat)
	r.PUT("/meats/:id", middleware.JWTAuthMiddleware("admin", "owner", "developer"), meatController.UpdateMeat)
}

func (mc *MeatController) CreateMeat(ctx *gin.Context) {
	username, err := utils.GetUsernameFromContext(ctx)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(ctx, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is creating a meat", username)
	var meat model.Meat
	if err := ctx.ShouldBindJSON(&meat); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(ctx, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	meat.CreatedBy = username
	meat.ID = uuid.New().String()
	err = mc.meatUseCase.CreateMeat(&meat)
	if err != nil {
		utils.HandleError(ctx, err)
		logrus.Errorf("[%v]%v", username, err)
		return
	}
	logrus.Info("Meat created successfully, meatname ", meat.Name)
	utils.SendResponse(ctx, http.StatusCreated, "Success", meat)
}

func (mc *MeatController) GetAllMeats(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Info("[", username, "] get all meats")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusBadRequest, "Invalid page number", nil)
		return
	}

	itemsPerPage, err := strconv.Atoi(c.DefaultQuery("itemsPerPage", "10"))
	if err != nil || itemsPerPage <= 0 {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusBadRequest, "Invalid itemsPerPage", nil)
		return
	}
	meats, totalPages, err := mc.meatUseCase.GetAllMeats(page, itemsPerPage)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get meats", nil)
		return
	}
	paginationData := map[string]interface{}{
		"page":         page,
		"itemsPerPage": itemsPerPage,
		"totalPages":   totalPages,
	}
	logrus.Infof("[%v] Get all meats %v, %v", username, paginationData, meats)
	utils.SendResponse(c, http.StatusOK, "Success", map[string]interface{}{"pagination": paginationData, "meats": meats})
}

func (mc *MeatController) GetMeatByName(c *gin.Context) {
	name := c.Param("name")
	username, err := utils.GetUsernameFromContext(c)
	meat, err := mc.meatUseCase.GetMeatByName(name)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
		return
	}

	logrus.Infof("[%v] Get meat %v", username, name)
	utils.SendResponse(c, http.StatusOK, "Success", meat)
}

func (mc *MeatController) GetMeatById(c *gin.Context) {
	id := c.Param("id")
	username, err := utils.GetUsernameFromContext(c)
	meat, err := mc.meatUseCase.GetMeatByName(id)
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get meat", nil)
		return
	}
	if meat == nil {
		utils.SendResponse(c, http.StatusNotFound, "Meat not found", nil)
		return
	}
	logrus.Infof("[%v] Get meat %v", username, id)
	utils.SendResponse(c, http.StatusOK, "Success", meat)
}

func (uc *MeatController) DeleteMeat(c *gin.Context) {
	meatID := c.Param("id")
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	if err := uc.meatUseCase.DeleteMeat(meatID); err != nil {
		if err == utils.ErrMeatNotFound {
			logrus.Errorf("[%v]%v", username, err)
			utils.SendResponse(c, http.StatusNotFound, "Meat not found", nil)
			return
		} else {
			logrus.Errorf("[%v]%v", username, err)
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete meat", nil)
			return

		}
	}
	logrus.Infof("[%v] Deleted meat %v", username, meatID)
	utils.SendResponse(c, http.StatusOK, "Success", nil)
}

func (uc *MeatController) UpdateMeat(ctx *gin.Context) {
	meatID := ctx.Param("id")
	username, err := utils.GetUsernameFromContext(ctx)
	var meat model.Meat
	if err := ctx.ShouldBindJSON(&meat); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(ctx, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}
	userName, err := utils.GetUsernameFromContext(ctx)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(ctx, http.StatusUnauthorized, "Invalid token", nil)
		return

	}
	meat.UpdatedBy = userName
	meat.ID = meatID
	logrus.Infof("[%s] is updating meat [%s]", userName, meatID)
	if err := uc.meatUseCase.UpdateMeat(&meat); err != nil {
		if err == utils.ErrMeatNotFound {
			logrus.Errorf("[%v]%v", username, err)
			utils.SendResponse(ctx, http.StatusNotFound, "Meat not found", nil)
			return
		} else {
			logrus.Errorf("[%v]%v", username, err)
			utils.SendResponse(ctx, http.StatusInternalServerError, "Failed to update meat", nil)
			return
		}
	}

	logrus.Infof("[%s] Meat updated successfully %v", userName, meat)
	utils.SendResponse(ctx, http.StatusOK, "Success", meat)
}
