package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvoiceNumberNotExist   = errors.New("Invoice number does not exist")
	ErrInvoiceAlreadyPaid      = errors.New("Invoice is already paid")
	ErrAmountGreaterThanTotal  = errors.New("Amount is greater than total transaction")
	ErrMeatNameAlreadyExist    = errors.New("Meatname already exists")
	ErrMeatNotFound            = errors.New("Meat not found")
	ErrCustomerNotFound        = errors.New("Customer not found")
	ErrCompanyNotFound         = errors.New("Company not found")
	ErrUserNotFound            = errors.New("User not found")
	ErrTransactionNotFound     = errors.New("Transaction not found")
	ErrTransactionAlreadyPaid  = errors.New("Transaction is already paid")
	ErrInvalidToken            = errors.New("Invalid token")
	ErrInvalidUsername         = errors.New("Invalid username")
	ErrInvalidPassword         = errors.New("Invalid password")
	ErrInvalidUsernamePassword = errors.New("Invalid username or password")
	ErrCompanyNameAlreadyExist = errors.New("Company name already exists")
	ErrInvalidMeatName         = errors.New("Invalid meat name")
	ErrInvalidAmount           = errors.New("Invalid amount")
	ErrInvalidInvoiceNumber    = errors.New("Invalid invoice number")
	ErrCreditPaymentNotFound   = errors.New("Credit payment not found")
	ErrInsufficientMeatStock   = errors.New("Insufficient meat stock")
	ErrMeatStockNotEnough      = errors.New("Meat stock not enough")
	ErrInvalidPrice            = errors.New("Invalid Meat price")
	ErrInvalidQty              = errors.New("Invalid quantity")
)

func HandleError(c *gin.Context, err error) {
	switch err {
	case ErrInvoiceNumberNotExist:
		SendResponse(c, http.StatusNotFound, err.Error(), nil)
	case ErrInvoiceAlreadyPaid:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrAmountGreaterThanTotal:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrMeatNameAlreadyExist:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrMeatNotFound:
		SendResponse(c, http.StatusNotFound, err.Error(), nil)
	case ErrCustomerNotFound:
		SendResponse(c, http.StatusNotFound, err.Error(), nil)
	case ErrCompanyNotFound:
		SendResponse(c, http.StatusNotFound, err.Error(), nil)
	case ErrUserNotFound:
		SendResponse(c, http.StatusNotFound, err.Error(), nil)
	case ErrTransactionNotFound:
		SendResponse(c, http.StatusNotFound, err.Error(), nil)
	case ErrTransactionAlreadyPaid:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrInvalidToken:
		SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	case ErrInvalidUsername:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrInvalidPassword:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrInvalidUsernamePassword:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrCompanyNameAlreadyExist:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrInvalidMeatName:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrInvalidAmount:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrInvalidInvoiceNumber:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrCreditPaymentNotFound:
		SendResponse(c, http.StatusNotFound, err.Error(), nil)
	case ErrInsufficientMeatStock:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrMeatStockNotEnough:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	case ErrInvalidPrice:
		SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	default:
		logrus.Error(err)
		SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
}
