package utils

import "errors"

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
	ErrInvalidPrice             = errors.New("Invalid price")
)
