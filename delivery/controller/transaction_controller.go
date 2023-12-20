package controller

import (
	"net/http"
	"strconv"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionUseCase usecase.TransactionUseCase
}

func NewTransactionController(r *gin.Engine, transactionUseCase usecase.TransactionUseCase) *TransactionController {
	controller := &TransactionController{
		transactionUseCase: transactionUseCase,
	}

	r.POST("/transactions", middleware.JWTAuthMiddleware("admin", "owner", "developer"), middleware.JSONMiddleware(), controller.CreateTransaction)
	r.GET("/transactions/:invoice_number", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.GetTransactionByInvoiceNumber)
	r.GET("/transactions", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.GetAllTransactions)
	r.DELETE("/transactions/:id", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.DeleteTransaction)

	return controller
}

func (tc *TransactionController) CreateTransaction(c *gin.Context) {

	username, err := utils.GetUsernameFromContext(c)
	if condition := err != nil; condition {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is creating a transaction", username)

	var request model.TransactionHeader
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if request.TxType == "" {
		utils.SendResponse(c, http.StatusBadRequest, "TxType is required", nil)
		return
	} else if request.TxType != "in" && request.TxType != "out" {
		utils.SendResponse(c, http.StatusBadRequest, "TxType must be in or out", nil)
		return
	} else if request.PaymentAmount <= 0 {
		utils.SendResponse(c, http.StatusBadRequest, "PaymentAmount must be greater than 0", nil)
		return
	}
	request.CreatedBy = username
	transaction, err := tc.transactionUseCase.CreateTransaction(&request)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
		return
	}

	logrus.Infof("[%v] Transaction successfully added with id = %v, invoice number = %v", username, transaction.ID, transaction.InvoiceNumber)
	utils.SendResponse(c, http.StatusOK, "Transaction created successfully", transaction)
}

func (tc *TransactionController) GetTransactionByID(c *gin.Context) {
	id := c.Param("id")
	username, err := utils.GetUsernameFromContext(c)
	transaction, err := tc.transactionUseCase.GetTransactionByID(id)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	logrus.Infof("[%v] Transaction found, id = %v", username)
	utils.SendResponse(c, http.StatusOK, "Transaction found", transaction)
}

func (tc *TransactionController) GetAllTransactions(c *gin.Context) {

	username, err := utils.GetUsernameFromContext(c)
	if condition := err != nil; condition {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is getting all transaction", username)

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

	transactions, totalPages, err := tc.transactionUseCase.GetAllTransactions(page, itemsPerPage)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	paginationData := map[string]interface{}{
		"page":         page,
		"itemsPerPage": itemsPerPage,
		"totalPages":   totalPages,
	}

	logrus.Info("Success get all transactions", paginationData, transactions)
	logrus.Infof("[%v] Transactions found with pagination data = %v and transactions = %v", username, paginationData, transactions)
	utils.SendResponse(c, http.StatusOK, "Transactions found", map[string]interface{}{"transactions": transactions, "pagination": paginationData})
}

func (tc *TransactionController) DeleteTransaction(c *gin.Context) {

	username, err := utils.GetUsernameFromContext(c)
	if condition := err != nil; condition {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
	}
	logrus.Infof("[%s] is deleting a transaction", username)
	id := c.Param("id")
	err = tc.transactionUseCase.DeleteTransaction(id)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	logrus.Infof("[%v] Transaction deleted successfully, id = %v", username, id)
	utils.SendResponse(c, http.StatusOK, "Transaction deleted successfully", nil)
}

func (tc *TransactionController) GetTransactionByInvoiceNumber(c *gin.Context) {

	username, err := utils.GetUsernameFromContext(c)
	if condition := err != nil; condition {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return

	}
	logrus.Infof("[%s] is geting a transaction by invoice number", username)

	invoice_number := c.Param("invoice_number")

	transaction, err := tc.transactionUseCase.GetTransactionByInvoiceNumber(invoice_number)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	logrus.Infof("[%v] Transaction found, invoice number = %v", username, invoice_number)
	utils.SendResponse(c, http.StatusOK, "Transaction found", transaction)
}
