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

	r.POST("/transactions", middleware.JWTAuthMiddleware(), middleware.JSONMiddleware(), controller.CreateTransaction)
	r.GET("/transactions/:invoice_number", middleware.JWTAuthMiddleware(), controller.GetTransactionByInvoiceNumber)
	r.GET("/transactions", middleware.JWTAuthMiddleware(), controller.GetAllTransactions)
	r.DELETE("/transactions/:id", middleware.JWTAuthMiddleware(), controller.DeleteTransaction)

	return controller
}

func (tc *TransactionController) CreateTransaction(c *gin.Context) {

	userName, err := utils.GetUsernameFromContext(c)
	if condition := err != nil; condition {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is creating a transaction", userName)

	var request model.TransactionHeader
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	if request.TxType == "" {
		utils.SendResponse(c, http.StatusBadRequest, "TxType is required", nil)
		return
	} else if request.TxType != "in" && request.TxType != "out" {
		utils.SendResponse(c, http.StatusBadRequest, "TxType must be in or out", nil)
		return

	}
	request.CreatedBy = userName
	transaction, err := tc.transactionUseCase.CreateTransaction(&request)
	if err != nil {
		logrus.Error(err)
		if err == utils.ErrAmountGreaterThanTotal {
			utils.SendResponse(c, http.StatusBadRequest, utils.ErrAmountGreaterThanTotal.Error(), nil)
		} else if err == utils.ErrCustomerNotFound {
			utils.SendResponse(c, http.StatusNotFound, utils.ErrCustomerNotFound.Error(), nil)

		} else if err == utils.ErrCompanyNotFound {
			utils.SendResponse(c, http.StatusNotFound, utils.ErrCompanyNotFound.Error(), nil)
		} else if err == utils.ErrMeatNotFound {
			utils.SendResponse(c, http.StatusNotFound, utils.ErrMeatNotFound.Error(), nil)
		} else {
			utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	logrus.Info("Transaction created successfully, ", transaction)
	utils.SendResponse(c, http.StatusOK, "Transaction created successfully", transaction)
}

func (tc *TransactionController) GetTransactionByID(c *gin.Context) {
	id := c.Param("id")

	transaction, err := tc.transactionUseCase.GetTransactionByID(id)
	if err != nil {
		utils.SendResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Transaction found", transaction)
}

func (tc *TransactionController) GetAllTransactions(c *gin.Context) {

	userName, err := utils.GetUsernameFromContext(c)
	if condition := err != nil; condition {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is getting all transaction", userName)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusBadRequest, "Invalid page number", nil)
		return
	}

	itemsPerPage, err := strconv.Atoi(c.DefaultQuery("itemsPerPage", "10"))
	if err != nil || itemsPerPage <= 0 {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusBadRequest, "Invalid itemsPerPage", nil)
		return
	}

	transactions, totalPages, err := tc.transactionUseCase.GetAllTransactions(page, itemsPerPage)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	paginationData := map[string]interface{}{
		"page":         page,
		"itemsPerPage": itemsPerPage,
		"totalPages":   totalPages,
	}

	logrus.Info("Success get all transactions", paginationData, transactions)
	utils.SendResponse(c, http.StatusOK, "Transactions found", map[string]interface{}{"transactions": transactions, "pagination": paginationData})
}

func (tc *TransactionController) DeleteTransaction(c *gin.Context) {

	userName, err := utils.GetUsernameFromContext(c)
	if condition := err != nil; condition {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
	}
	logrus.Infof("[%s] is deleting a transaction", userName)
	id := c.Param("id")
	err = tc.transactionUseCase.DeleteTransaction(id)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	logrus.Info("Transaction deleted successfully, id = ", id)
	utils.SendResponse(c, http.StatusOK, "Transaction deleted successfully", nil)
}

func (tc *TransactionController) GetTransactionByInvoiceNumber(c *gin.Context) {

	userName, err := utils.GetUsernameFromContext(c)
	if condition := err != nil; condition {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return

	}
	logrus.Infof("[%s] is geting a transaction by invoice number", userName)

	invoice_number := c.Param("invoice_number")

	transaction, err := tc.transactionUseCase.GetTransactionByInvoiceNumber(invoice_number)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	logrus.Info("Transaction found", transaction)
	utils.SendResponse(c, http.StatusOK, "Transaction found", transaction)
}
