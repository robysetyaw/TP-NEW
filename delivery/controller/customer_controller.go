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
	r.POST("/customers", middleware.JWTAuthMiddleware(), controller.CreateCustomer)
	r.GET("/customers", middleware.JWTAuthMiddleware(), controller.GetAllCustomer)
	r.GET("/customers/:id", middleware.JWTAuthMiddleware(), controller.GetCustomerByID)
	r.PUT("/customers/:id", middleware.JWTAuthMiddleware(), controller.UpdateCustomer)
	r.DELETE("/customers/:id", middleware.JWTAuthMiddleware(), controller.DeleteCustomer)
	r.GET("customers/company/:company_id", middleware.JWTAuthMiddleware(), controller.GetAllCustomerByCompanyId)
	// r.GET("/customers/transaction/:username", middleware.JWTAuthMiddleware(), controller.GetAllCustomerTransactions)
	return controller
}

func (cc *CustomerController) CreateCustomer(c *gin.Context) {
	var customer model.CustomerModel
	if err := c.ShouldBindJSON(&customer); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userName, err := utils.GetUsernameFromContext(c)
	logrus.Infof("[%v] created customer %v ", userName, customer.FullName)
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return

	}
	customer.CreatedBy = userName
	customer.Id = uuid.New().String()
	customer.Debt = 0

	customers, err := cc.customerUsecase.CreateCustomer(&customer)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	logrus.Infof("[%s] created succes create customer %s ", userName, customer.FullName)
	utils.SendResponse(c, http.StatusOK, "success insert data customer", customers)
}

func (cc *CustomerController) UpdateCustomer(c *gin.Context) {
	customerID := c.Param("id")
	var customer model.CustomerModel
	if err := c.ShouldBindJSON(&customer); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	userName, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	customer.UpdatedAt = time.Now()
	customer.UpdatedBy = userName
	customer.Id = customerID

	if err := cc.customerUsecase.UpdateCustomer(&customer); err != nil {
		if err == utils.ErrCustomerNotFound {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusNotFound, "Customer not found", nil)
			return
		} else {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	logrus.Infof("[%s] updated succes update customer %s ", userName, customer.FullName)
	utils.SendResponse(c, http.StatusOK, "success update data customer", customer)
}

func (cc *CustomerController) GetAllCustomerByCompanyId(c *gin.Context) {
	company_id := c.Param("company_id")
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	logrus.Infof("[%s] get all customer ", username)
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
	customers, totalPages, err := cc.customerUsecase.GetAllCustomerByCompanyId(page, itemsPerPage, company_id)
	if err != nil {
		if err == utils.ErrCustomerNotFound {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusNotFound, "Customer not found", nil)
			return
		} else if err == utils.ErrCompanyNotFound {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusNotFound, "Company not found", nil)
			return
		} else {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
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
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	logrus.Infof("[%s] get all customer ", username)
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
	customers, totalPages, err := cc.customerUsecase.GetAllCustomers(page, itemsPerPage)
	if err != nil {
		logrus.Error(err)
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
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)

	}
	logrus.Infof("[%s] is geting a customer", username)
	customer_id := c.Param("id")
	customers, err := cc.customerUsecase.GetCustomerById(customer_id)
	if err != nil {
		if err == utils.ErrCustomerNotFound {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusNotFound, "Customer not found", nil)
			return
		} else {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

	}
	utils.SendResponse(c, http.StatusOK, "Success", customers)
}

func (cc *CustomerController) DeleteCustomer(c *gin.Context) {
	customerId := c.Param("id")

	if err := cc.customerUsecase.DeleteCustomer(customerId); err != nil {
		if err == utils.ErrCustomerNotFound {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusNotFound, "Customer not found", nil)
			return
		} else {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	utils.SendResponse(c, http.StatusOK, "Success", nil)
}

// func (cc *CustomerController) GetAllCustomerTransactions(c *gin.Context) {
// 	userName := c.Param("username")
// 	customerTransactions, err := cc.customerUsecase.GetAllCustomerTransactions(userName)
// 	if err != nil {
// 		appError := apperror.AppError{}
// 		if errors.As(err, &appError) {
// 			fmt.Printf(" cc.customerUsecase.GetAllCustomers() : %v ", err.Error())
// 			c.JSON(http.StatusBadGateway, gin.H{
// 				"errorMessage": appError.Error(),
// 			})
// 		} else {
// 			fmt.Printf("ServiceHandler.InsertService() 2 : %v ", err.Error())
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"errorMessage": "failed to get customers",
// 			})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, customerTransactions)
// }
