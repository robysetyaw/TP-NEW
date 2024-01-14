package controller

import (
	"net/http"
	"strconv"
	"time"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"

	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CustomerController struct {
	customerUsecase usecase.CustomerUsecase
}

func NewCustomerController(r *gin.Engine, customerUsecase usecase.CustomerUsecase) *CustomerController {
	controller := &CustomerController{
		customerUsecase: customerUsecase,
	}
	r.POST("/customers", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.CreateCustomer)
	r.GET("/customers", middleware.JWTAuthMiddleware("employee", "admin", "owner", "developer"), controller.GetAllCustomer)
	r.GET("/customers/:id", middleware.JWTAuthMiddleware("employee", "admin", "owner", "developer"), controller.GetCustomerByID)
	r.PUT("/customers/:id", middleware.JWTAuthMiddleware("admin", "owner", "developer"), controller.UpdateCustomer)
	r.DELETE("/customers/:id", middleware.JWTAuthMiddleware("owner", "developer"), controller.DeleteCustomer)
	r.GET("customers/company/:company_id", middleware.JWTAuthMiddleware("employee", "admin", "owner", "developer"), controller.GetAllCustomerByCompanyId)
	r.GET("/customers/transaction/:id", middleware.JWTAuthMiddleware("employee", "admin", "owner", "developer"), controller.GetAllTransactionsByCustomerId)
	return controller
}

func (cc *CustomerController) CreateCustomer(c *gin.Context) {
	var customer model.CustomerModel
	if err := c.ShouldBindJSON(&customer); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	username, err := utils.GetUsernameFromContext(c)
	logrus.Infof("[%v] created customer %v ", username, customer.FullName)
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return

	}
	customer.CreatedBy = username
	customer.Id = uuid.New().String()
	customer.Debt = 0

	customers, err := cc.customerUsecase.CreateCustomer(&customer)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
		return
	}

	logrus.Infof("[%s] created succes create customer %s ", username, customer.FullName)
	utils.SendResponse(c, http.StatusOK, "success insert data customer", customers)
}

func (cc *CustomerController) UpdateCustomer(c *gin.Context) {
	customerID := c.Param("id")
	var customer model.CustomerModel
	if err := c.ShouldBindJSON(&customer); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	customer.UpdatedAt = time.Now()
	customer.UpdatedBy = username
	customer.Id = customerID

	if err := cc.customerUsecase.UpdateCustomer(&customer); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
	}

	logrus.Infof("[%s] updated succes update customer %s ", username, customer.FullName)
	utils.SendResponse(c, http.StatusOK, "success update data customer", customer)
}

func (cc *CustomerController) GetAllCustomerByCompanyId(c *gin.Context) {
	company_id := c.Param("company_id")
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	logrus.Infof("[%s] get all customer ", username)
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
	customers, totalPages, err := cc.customerUsecase.GetAllCustomerByCompanyId(page, itemsPerPage, company_id)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
	}

	logrus.Infof("[%s] get all customer ", username)
	utils.SendResponse(c, http.StatusOK, "success get all customer", map[string]interface{}{
		"customers":  customers,
		"totalPages": totalPages,
	})
}

func (cc *CustomerController) GetAllCustomer(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	logrus.Infof("[%s] get all customer ", username)
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
	customers, totalPages, err := cc.customerUsecase.GetAllCustomers(page, itemsPerPage)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	logrus.Infof("[%s] get all customer ", username)
	utils.SendResponse(c, http.StatusOK, "success get all customer by company_id", map[string]interface{}{
		"customers":  customers,
		"totalPages": totalPages,
	})
}

func (cc *CustomerController) GetCustomerByID(c *gin.Context) {

	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)

	}
	logrus.Infof("[%s] is geting a customer", username)
	customer_id := c.Param("id")
	customers, err := cc.customerUsecase.GetCustomerById(customer_id)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
		return

	}
	utils.SendResponse(c, http.StatusOK, "Success", customers)
}

func (cc *CustomerController) DeleteCustomer(c *gin.Context) {
	customerId := c.Param("id")
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	if err := cc.customerUsecase.DeleteCustomer(customerId); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Success", nil)
}

func (cc *CustomerController) GetAllTransactionsByCustomerId(c *gin.Context) {
	customerId := c.Param("id")
	paymentStatus := c.Query("payment_status")

	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	logrus.Infof("[%s] get all transaction by customer id", username)

	if paymentStatus != "" && paymentStatus != "paid" && paymentStatus != "unpaid" {
		logrus.Errorf("[%s] Invalid payment status, please use paid or unpaid or empty ", username)
		utils.SendResponse(c, http.StatusBadRequest, "Invalid payment status, please use paid or unpaid or empty ", nil)
		return
	}

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

	customerTransactions, totalPages,err := cc.customerUsecase.GetAllTransactionsByCustomerId(customerId, paymentStatus, page, itemsPerPage)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
		return
	}

	paginationData := map[string]interface{}{
		"page":         page,
		"itemsPerPage": itemsPerPage,
		"totalPages":   totalPages,
	}

	logrus.Infof("[%s] get all transaction by customer id", username)
	logrus.Infof("[%v] Transactions found with pagination data = %v and transactions = %v", username, paginationData, customerTransactions)
	utils.SendResponse(c, http.StatusOK, "Success", map[string]interface{}{"transactions": customerTransactions, "pagination": paginationData})

}
