package controller

import (
	"net/http"
	"time"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type DailyExpenditureController struct {
	dailyExpenditureUseCase usecase.DailyExpenditureUseCase
}

func NewDailyExpenditureController(r *gin.Engine, deUC usecase.DailyExpenditureUseCase) *DailyExpenditureController {
	controller := &DailyExpenditureController{
		dailyExpenditureUseCase: deUC,
	}

	r.POST("/daily-expenditures", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.CreateDailyExpenditure)
	r.PUT("/daily-expenditures/:id", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.UpdateDailyExpenditure)
	r.GET("/daily-expenditures/:id", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.GetDailyExpenditureByID)
	r.GET("/daily-expenditures", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.GetAllDailyExpenditures)
	r.DELETE("/daily-expenditures/:id", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.DeleteDailyExpenditure)

	return controller
}

func (dec *DailyExpenditureController) GetTotalExpenditureByDateRange(c *gin.Context) {
	var payload struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", payload.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
		return
	}

	endDate, err := time.Parse("2006-01-02", payload.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
		return
	}

	total, err := dec.dailyExpenditureUseCase.GetTotalExpenditureByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total expenditure"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_expenditure": total})
}

func (dec *DailyExpenditureController) CreateDailyExpenditure(c *gin.Context) {
	var expenditure model.DailyExpenditure
	if err := c.ShouldBindJSON(&expenditure); err != nil {
		utils.HandleError(c, err)
		return
	}

	token, err := utils.ExtractTokenFromAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid authorization header", nil)
		return
	}

	claims, err := utils.VerifyJWTToken(token)
	if err != nil {
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	userID := claims["user_id"].(string)
	userName := claims["username"].(string)
	logrus.Infof("[%s] is creating daily expenditure", userName)
	expenditure.UserID = userID
	expenditure.CreatedBy = userName
	expenditure.ID = uuid.New().String()

	if err := dec.dailyExpenditureUseCase.CreateDailyExpenditure(&expenditure); err != nil {
		utils.HandleError(c, err)
		return
	}
	logrus.Infof("[%s] created daily expenditure", userName)
	utils.SendResponse(c, http.StatusOK, "Success", nil)
}

func (dec *DailyExpenditureController) UpdateDailyExpenditure(c *gin.Context) {
	expenditureID := c.Param("id")

	var expenditure model.DailyExpenditure
	if err := c.ShouldBindJSON(&expenditure); err != nil {
		utils.HandleError(c, err)
		return
	}

	expenditure.ID = expenditureID

	token, err := utils.ExtractTokenFromAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid authorization header", nil)
		return
	}

	claims, err := utils.VerifyJWTToken(token)
	if err != nil {
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is updating daily expenditure", claims["username"].(string))
	expenditure.UpdatedBy = claims["username"].(string)
	expenditure.IsActive = true

	if err := dec.dailyExpenditureUseCase.UpdateDailyExpenditure(&expenditure); err != nil {
		utils.HandleError(c, err)
		return
	}
	logrus.Infof("[%s] updated daily expenditure", claims["username"].(string))
	utils.SendResponse(c, http.StatusOK, "Success", nil)
}

func (dec *DailyExpenditureController) GetDailyExpenditureByID(c *gin.Context) {
	expenditureID := c.Param("id")
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	logrus.Infof("[%s] is geting daily expenditure", username)

	expenditure, err := dec.dailyExpenditureUseCase.GetDailyExpenditureByID(expenditureID)
	if err != nil {
		logrus.Infof("[%s] geted daily expenditure", username)
		utils.HandleError(c, err)
		return
	}

	if expenditure == nil {
		logrus.Infof("[%s] geted daily expenditure", username)
		utils.SendResponse(c, http.StatusNotFound, "Daily expenditure not found", nil)
		return
	}
	logrus.Infof("[%s] geted daily expenditure", username)
	utils.SendResponse(c, http.StatusOK, "Success", expenditure)
}

func (dec *DailyExpenditureController) GetAllDailyExpenditures(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	logrus.Infof("[%s] is geting all daily expenditure", username)
	expenditures, err := dec.dailyExpenditureUseCase.GetAllDailyExpenditures()
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	logrus.Infof("[%s] succes geting all daily expenditure", username)
	utils.SendResponse(c, http.StatusOK, "Success", expenditures)
}

func (dec *DailyExpenditureController) DeleteDailyExpenditure(c *gin.Context) {
	expenditureID := c.Param("id")

	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		utils.HandleError(c, err)
		return
	}
	logrus.Infof("[%s] is deleting daily expenditure [%s]", username, expenditureID)

	if err := dec.dailyExpenditureUseCase.DeleteDailyExpenditure(expenditureID); err != nil {
		utils.HandleError(c, err)
		return
	}
	logrus.Infof("[%s] succes delete daily expenditure [%s]", username, expenditureID)
	utils.SendResponse(c, http.StatusOK, "Success", nil)
}
