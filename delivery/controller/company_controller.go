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
	logrus.Infof("[%s] is creating a company", username)

	var company model.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Bad request", nil)
		return
	}
	company.CreatedBy = username
	company.ID = uuid.New().String()
	company.IsActive = true

	if err := cc.companyUseCase.CreateCompany(&company); err != nil {
		if err == utils.ErrCompanyNameAlreadyExist {
			logrus.WithField("error", err).Error("Company name already exists")
			utils.SendResponse(c, http.StatusConflict, "Company name already exists", nil)
		} else {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to create company", nil)
		}
	}

	utils.SendResponse(c, http.StatusOK, "Success create company", company)
}

func (cc *CompanyController) UpdateCompany(c *gin.Context) {
	companyID := c.Param("id")

	var company model.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusBadRequest, "Bad request", nil)
		return
	}

	userName, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, "Invalid token", nil)
		return
	}
	company.UpdatedBy = userName
	company.ID = companyID

	if err := cc.companyUseCase.UpdateCompany(&company); err != nil {
		if err == utils.ErrCompanyNotFound {
			logrus.WithField("error", err).Error("Company not found")
			utils.SendResponse(c, http.StatusNotFound, "Company not found", nil)
		} else {
			logrus.Error(err)
			utils.SendResponse(c, http.StatusInternalServerError, "Failed to update company", nil)
		}
	}

	utils.SendResponse(c, http.StatusOK, "Success update company", company)
}

func (cc *CompanyController) GetCompanyById(c *gin.Context) {
	username,err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, "Invalid token", nil)
		return
	}
	companyId := c.Param("id")
	logrus.Infof("[%s] is geting a company", username)
	company, err := cc.companyUseCase.GetCompanyById(companyId)
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get company", nil)
		return
	}

	if company == nil {
		utils.SendResponse(c, http.StatusNotFound, "Company not found", nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Success get all company", company)
}

func (cc *CompanyController) GetAllCompany(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, "Invalid token", nil)
		return
	}
	logrus.Infof("[%s] is geting a company", username)
	companies, err := cc.companyUseCase.GetAllCompany()
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to get company", nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Success get company", companies)
}

func (cc *CompanyController) DeleteCompany(c *gin.Context) {
	username, err := utils.GetUsernameFromContext(c)
	if err != nil {
		logrus.Error(err)
		utils.SendResponse(c, http.StatusInternalServerError, "Invalid token", nil)
		return
	}
	companyId := c.Param("id")
	logrus.Infof("[%s] is deleting a company", username)
	if err := cc.companyUseCase.DeleteCompany(companyId); err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete company", nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Success delete company", nil)
}
