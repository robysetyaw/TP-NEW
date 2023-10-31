package controller

import (
	"net/http"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

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
	// r.GET("/transactions/:id", middleware.JWTAuthMiddleware(), controller.GetTransactionByID)
	r.GET("/transactions/:invoice_number", middleware.JWTAuthMiddleware(), controller.GetTransactionByInvoiceNumber)
	r.GET("/transactions", middleware.JWTAuthMiddleware(), controller.GetAllTransactions)
	r.DELETE("/transactions/:id", middleware.JWTAuthMiddleware(), controller.DeleteTransaction)

	return controller
}

func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	var request model.TransactionHeader
	if err := c.ShouldBindJSON(&request); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	transaction, err := tc.transactionUseCase.CreateTransaction(&request)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully"})
	utils.SendResponse(c, http.StatusOK, "Transaction created successfully", transaction)
}

func (tc *TransactionController) GetTransactionByID(c *gin.Context) {
	id := c.Param("id")

	transaction, err := tc.transactionUseCase.GetTransactionByID(id)
	if err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		utils.SendResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, transaction)
	utils.SendResponse(c, http.StatusOK, "Transaction found", transaction)
}

func (tc *TransactionController) GetAllTransactions(c *gin.Context) {
	transactions, err := tc.transactionUseCase.GetAllTransactions()
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, transactions)
	utils.SendResponse(c, http.StatusOK, "Transactions found", transactions)
}

func (tc *TransactionController) DeleteTransaction(c *gin.Context) {
	id := c.Param("id")

	err := tc.transactionUseCase.DeleteTransaction(id)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
	utils.SendResponse(c, http.StatusOK, "Transaction deleted successfully", nil)
}

func (tc *TransactionController) GetTransactionByInvoiceNumber(c *gin.Context) {
	invoice_number := c.Param("invoice_number")

	transaction, err := tc.transactionUseCase.GetTransactionByInvoiceNumber(invoice_number)
	if err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		utils.SendResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, transaction)
	utils.SendResponse(c, http.StatusOK, "Transaction found", transaction)
}
