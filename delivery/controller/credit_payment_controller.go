package controller

import (
	"net/http"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
)

type CreditPaymentController struct {
	creditPaymentUseCase usecase.CreditPaymentUseCase
}

func NewCreditPaymentController(r *gin.Engine, creditPaymentUseCase usecase.CreditPaymentUseCase) *CreditPaymentController {
	controller := &CreditPaymentController{
		creditPaymentUseCase: creditPaymentUseCase,
	}
	r.POST("/credit_payment", middleware.JWTAuthMiddleware(), controller.CreateCreditPayment)
	return controller
}

func (cc *CreditPaymentController) CreateCreditPayment(c *gin.Context) {
	var payment model.CreditPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := cc.creditPaymentUseCase.CreateCreditPayment(&payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "Credit payment created successfully"})
	utils.SendResponse(c,http.StatusOK,"Credit Payment succesfully add",transaction)	
}

func (cc *CreditPaymentController) GetCreditPayments(c *gin.Context) {
	payments, err := cc.creditPaymentUseCase.GetCreditPayments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (cc *CreditPaymentController) GetCreditPaymentByID(c *gin.Context) {
	id := c.Param("id")

	payment, err := cc.creditPaymentUseCase.GetCreditPaymentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (cc *CreditPaymentController) UpdateCreditPayment(c *gin.Context) {
	id := c.Param("id")

	var payment model.CreditPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment.ID = id

	err := cc.creditPaymentUseCase.UpdateCreditPayment(&payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credit payment updated successfully"})
}
