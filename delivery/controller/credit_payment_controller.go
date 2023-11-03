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
	r.GET("/credit_payments/:invoice_number", middleware.JWTAuthMiddleware(), controller.GetCreditPaymentsByInvoiceNumber)
	return controller
}

func (cc *CreditPaymentController) CreateCreditPayment(c *gin.Context) {
	var payment model.CreditPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	transaction, err := cc.creditPaymentUseCase.CreateCreditPayment(&payment)
	if err != nil {
		// Handle error from use case
		if err == utils.ErrInvoiceNumberNotExist {
			utils.SendResponse(c, http.StatusNotFound, "Invoice number does not exist", nil)
		} else if err == utils.ErrInvoiceAlreadyPaid {
			utils.SendResponse(c, http.StatusConflict, "Invoice is already paid", nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}
	// c.JSON(http.StatusOK, gin.H{"message": "Credit payment created successfully"})
	utils.SendResponse(c, http.StatusOK, "Credit Payment successfully added", transaction)
}

func (cc *CreditPaymentController) GetCreditPayments(c *gin.Context) {
	payments, err := cc.creditPaymentUseCase.GetCreditPayments()
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, payments)
	utils.SendResponse(c, http.StatusOK, "Success get credit payments", payments)
}

func (cc *CreditPaymentController) GetCreditPaymentByID(c *gin.Context) {
	id := c.Param("id")

	payment, err := cc.creditPaymentUseCase.GetCreditPaymentByID(id)
	if err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, payment)
	utils.SendResponse(c, http.StatusOK, "Success get credit payment by ID", payment)
}

func (cc *CreditPaymentController) UpdateCreditPayment(c *gin.Context) {
	id := c.Param("id")

	var payment model.CreditPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	payment.ID = id

	err := cc.creditPaymentUseCase.UpdateCreditPayment(&payment)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "Credit payment updated successfully"})
	utils.SendResponse(c, http.StatusOK, "Credit payment updated successfully", nil)
}

func (cc *CreditPaymentController) GetCreditPaymentsByInvoiceNumber(c *gin.Context) {
	invoice_number := c.Param("invoice_number")

	payments, err := cc.creditPaymentUseCase.GetCreditPaymentsByInvoiceNumber(invoice_number)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, payments)
	utils.SendResponse(c, http.StatusOK, "Success get credit payments", payments)
}
