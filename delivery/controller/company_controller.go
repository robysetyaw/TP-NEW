package controller

import (
	"net/http"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CompanyController struct {
	companyUseCase usecase.CompanyUseCase
}

func NewCompanyController(r *gin.Engine, companyUseCase usecase.CompanyUseCase) *CompanyController {
	controller := &CompanyController{
		companyUseCase: companyUseCase,
	}
	r.POST("/companies", middleware.JWTAuthMiddleware(), controller.CreateCompany)
	r.PUT("/companies/:id", middleware.JWTAuthMiddleware(), controller.UpdateCompany)
	r.GET("/companies/:id", middleware.JWTAuthMiddleware(), controller.GetCompanyById)
	r.GET("/companies", middleware.JWTAuthMiddleware(), controller.GetAllCompany)
	r.DELETE("/companies/:id", middleware.JWTAuthMiddleware(), controller.DeleteCompany)

	return controller
}

func (cc *CompanyController) CreateCompany(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if condition := err != nil; condition {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
	}
	logrus.Infof("User %s is creating a company", username)

	var company model.Company
	company.CreatedBy = username
	company.ID = uuid.New().String()
	company.IsActive = true

	if err := cc.companyUseCase.CreateCompany(&company); err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to create company", nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Success create company", company)
}

func (cc *CompanyController) UpdateCompany(c *gin.Context) {
	companyID := c.Param("id")

	var company model.Company
	if err := c.ShouldBindJSON(&company); err != nil {
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
	company.UpdatedBy = userName
	company.ID = companyID

	if err := cc.companyUseCase.UpdateCompany(&company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update update data company",
	})
}

func (cc *CompanyController) GetCompanyById(c *gin.Context) {
	companyId := c.Param("id")

	company, err := cc.companyUseCase.GetCompanyById(companyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get company"})
		return
	}

	if company == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func (cc *CompanyController) GetAllCompany(c *gin.Context) {
	companies, err := cc.companyUseCase.GetAllCompany()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get companies"})
		return
	}

	c.JSON(http.StatusOK, companies)
}

func (cc *CompanyController) DeleteCompany(c *gin.Context) {
	companyId := c.Param("id")

	if err := cc.companyUseCase.DeleteCompany(companyId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}
