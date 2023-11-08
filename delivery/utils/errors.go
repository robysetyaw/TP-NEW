package utils

import "errors"

var (
	ErrInvoiceNumberNotExist  = errors.New("Invoice number does not exist")
	ErrInvoiceAlreadyPaid     = errors.New("Invoice is already paid")
	ErrAmountGreaterThanTotal = errors.New("Amount is greater than total transaction")
	ErrMeatNameAlreadyExist   = errors.New("Meatname already exists")
	ErrMeatNotFound           = errors.New("Meat not found")
	ErrCustomerNotFound       = errors.New("Customer not found")
)
