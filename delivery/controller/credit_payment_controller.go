package controller

import (
	"net/http"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	userName, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is creating a credit payment", userName)

	var payment model.CreditPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	payment.CreatedBy = userName
	transaction, err := cc.creditPaymentUseCase.CreateCreditPayment(&payment)
	if err != nil {
		if err == utils.ErrInvoiceNumberNotExist {
			utils.SendResponse(c, http.StatusNotFound, "Invoice number does not exist", nil)
		} else if err == utils.ErrInvoiceAlreadyPaid {
			utils.SendResponse(c, http.StatusConflict, "Invoice is already paid", nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		}
		logrus.Error(err)
		return
	}
	logrus.Info("Credit Payment successfully added")
	utils.SendResponse(c, http.StatusOK, "Credit Payment successfully added", transaction)
}

func (cc *CreditPaymentController) GetCreditPayments(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is geting a credit payment", username)
	payments, err := cc.creditPaymentUseCase.GetCreditPayments()
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	logrus.Info("Credit Payment found", payments)
	utils.SendResponse(c, http.StatusOK, "Success get credit payments", payments)
}

func (cc *CreditPaymentController) GetCreditPaymentByID(c *gin.Context) {
	id := c.Param("id")

	payment, err := cc.creditPaymentUseCase.GetCreditPaymentByID(id)
	if err != nil {
		if err == utils.ErrCreditPaymentNotFound {
			utils.SendResponse(c, http.StatusNotFound, "Credit payment not found", nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		}
		logrus.Error(err)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Success get credit payment by ID", payment)
}

func (cc *CreditPaymentController) UpdateCreditPayment(c *gin.Context) {
	id := c.Param("id")

	var payment model.CreditPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	payment.ID = id

	err := cc.creditPaymentUseCase.UpdateCreditPayment(&payment)
	if err != nil {
		if err == utils.ErrCreditPaymentNotFound {
			utils.SendResponse(c, http.StatusNotFound, "Credit payment not found", nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		}
		logrus.Error(err)
		return
	}

	logrus.Info("Credit payment updated successfully", payment)
	utils.SendResponse(c, http.StatusOK, "Credit payment updated successfully", nil)
}

func (cc *CreditPaymentController) GetCreditPaymentsByInvoiceNumber(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is geting a credit payment by invoice number %s", username, c.Param("invoice_number"))
	invoice_number := c.Param("invoice_number")

	payments, err := cc.creditPaymentUseCase.GetCreditPaymentsByInvoiceNumber(invoice_number)
	if err != nil {
		if err == utils.ErrCreditPaymentNotFound {
			utils.SendResponse(c, http.StatusNotFound, "Credit payment not found", nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		}
		logrus.Error(err)
		return
	}

	logrus.Info("Credit Payment found", payments)
	utils.SendResponse(c, http.StatusOK, "Success get credit payments", payments)
}
