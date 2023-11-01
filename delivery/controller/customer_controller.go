package controller

import (
	"net/http"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"

	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	// r.GET("/customers/transaction/:username", middleware.JWTAuthMiddleware(), controller.GetAllCustomerTransactions)
	return controller
}

func (cc *CustomerController) CreateCustomer(c *gin.Context) {
	var customer model.CustomerModel
	if err := c.ShouldBindJSON(&customer); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token, err := utils.ExtractTokenFromAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid authorization header", nil)
		return
	}

	claims, err := utils.VerifyJWTToken(token)
	if err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	userName := claims["username"].(string)
	customer.CreatedBy = userName
	customer.Id = uuid.New().String()
	customer.Debt = 0

	customers, err := cc.customerUsecase.CreateCustomer(&customer)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "success insert data customer",
	// })
	utils.SendResponse(c, http.StatusOK, "success insert data customer", customers)
}

func (cc *CustomerController) UpdateCustomer(c *gin.Context) {
	customerID := c.Param("id")
	var customer model.CustomerModel
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.ExtractTokenFromAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
		return
	}

	claims, err := utils.VerifyJWTToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	userName := claims["username"].(string)
	customer.UpdatedBy = userName
	customer.Id = customerID

	if err := cc.customerUsecase.UpdateCustomer(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success update data company",
	})
}

func (cc *CustomerController) GetAllCustomer(c *gin.Context) {
	customers, err := cc.customerUsecase.GetAllCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, customers)
}

func (cc *CustomerController) GetCustomerByID(c *gin.Context) {
	username := c.Param("username")

	expenditure, err := cc.customerUsecase.GetCustomerById(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, expenditure)
}

func (cc *CustomerController) DeleteCustomer(c *gin.Context) {
	customerId := c.Param("id")

	if err := cc.customerUsecase.DeleteCustomer(customerId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted successfully"})
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
