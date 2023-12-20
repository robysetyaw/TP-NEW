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
	r.POST("/credit_payment", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.CreateCreditPayment)
	r.GET("/credit_payments/:invoice_number", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.GetCreditPaymentsByInvoiceNumber)
	return controller
}

func (cc *CreditPaymentController) CreateCreditPayment(c *gin.Context) {

	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is creating a credit payment", username)

	var payment model.CreditPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	payment.CreatedBy = username
	transaction, err := cc.creditPaymentUseCase.CreateCreditPayment(&payment)
	if err != nil {
		utils.HandleError(c, err)
		logrus.Errorf("[%v]%v", username, err)
		return
	}
	logrus.Infof("[%v] Credit Payment successfully added with id = %v", username, transaction.Transaction.ID)
	utils.SendResponse(c, http.StatusOK, "Credit Payment successfully added", transaction)
}

func (cc *CreditPaymentController) GetCreditPayments(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is geting a credit payment", username)
	payments, err := cc.creditPaymentUseCase.GetCreditPayments()
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	logrus.Infof("[%v] Credit Payment found", username)
	utils.SendResponse(c, http.StatusOK, "Success get credit payments", payments)
}

func (cc *CreditPaymentController) GetCreditPaymentByID(c *gin.Context) {
	id := c.Param("id")
	username, err := utils.GetUsernameFromContext(c)
	payment, err := cc.creditPaymentUseCase.GetCreditPaymentByID(id)
	if err != nil {
		utils.HandleError(c, err)
		logrus.Errorf("[%v]%v", username, err)
		return
	}
	logrus.Infof("[%v] Credit Payment found, paymet id = %v", username, payment.ID)
	utils.SendResponse(c, http.StatusOK, "Success get credit payment by ID", payment)
}

func (cc *CreditPaymentController) UpdateCreditPayment(c *gin.Context) {
	id := c.Param("id")

	username, err := utils.GetUsernameFromContext(c)
	var payment model.CreditPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	payment.ID = id
	err = cc.creditPaymentUseCase.UpdateCreditPayment(&payment)
	if err != nil {
		utils.HandleError(c, err)
		logrus.Errorf("[%v]%v", username, err)
		return
	}

	logrus.Infof("[%v] Credit Payment updated successfully, payment id = %v", username, payment.ID)
	utils.SendResponse(c, http.StatusOK, "Credit payment updated successfully", nil)
}

func (cc *CreditPaymentController) GetCreditPaymentsByInvoiceNumber(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is geting a credit payment by invoice number %s", username, c.Param("invoice_number"))
	invoice_number := c.Param("invoice_number")

	payments, err := cc.creditPaymentUseCase.GetCreditPaymentsByInvoiceNumber(invoice_number)
	if err != nil {
		utils.HandleError(c, err)
		logrus.Errorf("[%v]%v", username, err)
		return
	}

	logrus.Infof("[%v] Credit Payment found by invoice number = %v", username, invoice_number)
	utils.SendResponse(c, http.StatusOK, "Success get credit payments", payments)
}
