package utils

import "errors"

var (
    ErrInvoiceNumberNotExist = errors.New("Invoice number does not exist")
    ErrInvoiceAlreadyPaid    = errors.New("Invoice is already paid")
)