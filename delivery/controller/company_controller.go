package controller

import (
	"net/http"
	"trackprosto/delivery/middleware"
	"trackprosto/delivery/utils"
	model "trackprosto/models"
	"trackprosto/models/dto"
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
	r.POST("/companies", middleware.JWTAuthMiddleware("admin","owner","developer"), controller.CreateCompany)
	r.PUT("/companies/:id", middleware.JWTAuthMiddleware("admin","owner","developer"), controller.UpdateCompany)
	r.GET("/companies/:id", middleware.JWTAuthMiddleware("employee","admin","owner","developer"), controller.GetCompanyById)
	r.GET("/companies", middleware.JWTAuthMiddleware("employee","admin","owner","developer"), controller.GetAllCompany)
	r.DELETE("/companies/:id", middleware.JWTAuthMiddleware("admin","owner","developer"), controller.DeleteCompany)

	return controller
}

func (cc *CompanyController) CreateCompany(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if condition := err != nil; condition {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid token", nil)
	}
	logrus.Infof("[%s] is creating a company", username)

	var company model.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusBadRequest, "Bad request", nil)
		return
	}
	company.CreatedBy = username
	company.ID = uuid.New().String()
	company.IsActive = true

	if err := cc.companyUseCase.CreateCompany(&company); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
	}

	logrus.Infof("[%v] Created company %v", username, company.CompanyName)
	utils.SendResponse(c, http.StatusOK, "Success create company", company)
}

func (cc *CompanyController) UpdateCompany(c *gin.Context) {
	companyID := c.Param("id")
	username, err := utils.GetUsernameFromContext(c)
	var CompanyRequest dto.CompanyRequest
	if err := c.ShouldBindJSON(&CompanyRequest); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusBadRequest, "Bad request", nil)
		return
	}
	userName, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, "Invalid token", nil)
		return
	}

	CompanyRequest.UpdatedBy = userName
	CompanyRequest.ID = companyID

	companyResponse, err := cc.companyUseCase.UpdateCompany(&CompanyRequest)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
		return
	}

	logrus.Infof("[%v] Updated company %v", username, companyID)
	utils.SendResponse(c, http.StatusOK, "Success update company", companyResponse)
}

func (cc *CompanyController) GetCompanyById(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, "Invalid token", nil)
		return
	}
	companyId := c.Param("id")
	logrus.Infof("[%s] is geting a company", username)
	company, err := cc.companyUseCase.GetCompanyById(companyId)

	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
		return
	}
	logrus.Infof("[%v] Company found", username)
	utils.SendResponse(c, http.StatusOK, "Success get company", company)
}

func (cc *CompanyController) GetAllCompany(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is geting a company", username)
	companies, err := cc.companyUseCase.GetAllCompany()
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
		return
	}

	logrus.Infof("[%v] Company found", username)
	utils.SendResponse(c, http.StatusOK, "Success get all company", companies)
}

func (cc *CompanyController) DeleteCompany(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.SendResponse(c, http.StatusInternalServerError, "Invalid token", nil)
		return
	}
	companyId := c.Param("id")
	logrus.Infof("[%s] is deleting a company", username)
	if err := cc.companyUseCase.DeleteCompany(companyId); err != nil {
		logrus.Errorf("[%v]%v", username, err)
		utils.HandleError(c, err)
	}
	logrus.Infof("[%v] Deleted company %v", username, companyId)
	utils.SendResponse(c, http.StatusOK, "Success delete company", nil)
}
