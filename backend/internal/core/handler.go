package core

import (
	db "herp/db/sqlc"
	"herp/internal/auth"
	"herp/internal/config"
	"herp/internal/utils"
	"herp/pkg/monitoring/logging"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Handler struct {
	service CoreInterface
	config  *config.Config
	logger  *logging.Logger
}

func NewHandler(service CoreInterface, c *config.Config, l *logging.Logger) *Handler {
	return &Handler{
		service: service,
		config:  c,
		logger:  l,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup, authSvc *auth.Service) {
	core := r.Group("/core")
	core.Use(auth.AuthMiiddleware(authSvc))

	// Business endpoints
	business := core.Group("/business")
	{
		business.POST("", auth.PermissionMiddleware(authSvc, "core:create_business"), h.createBusiness)
		business.GET("/:id", auth.PermissionMiddleware(authSvc, "core:view_business"), h.getBusiness)
		business.PUT("/:id", auth.PermissionMiddleware(authSvc, "core:update_business"), h.updateBusiness)
		business.DELETE("/:id", auth.PermissionMiddleware(authSvc, "core:delete_business"), h.deleteBusiness)
		business.GET("/all", auth.PermissionMiddleware(authSvc, "core:view_business"), h.listBusinesses)
	}

	branch := core.Group("/branch")
	{
		branch.POST("", auth.PermissionMiddleware(authSvc, "core:create_branch"), h.createBranch)
		branch.GET("/:id", auth.PermissionMiddleware(authSvc, "core:view_branch"), h.getBranch)
		branch.PUT("/:id", auth.PermissionMiddleware(authSvc, "core:update_branch"), h.updateBranch)
		branch.DELETE("/:id", auth.PermissionMiddleware(authSvc, "core:delete_branch"), h.deleteBranch)
		branch.GET("", auth.PermissionMiddleware(authSvc, "core:view_branch"), h.listBranches)
	}
}

type CreateBusinessParams struct {
	Name              string   `json:"name" example:"Palmwineexpress hotels" binding:"required"`
	Email             *string  `json:"email" binding:"omitempty" example:"admin@palmwinexpress.com"`
	Website           *string  `json:"website" binding:"omitempty" example:"https://palmwinexpress.com"`
	TaxID             *string  `json:"tax_id" binding:"omitempty" example:"123456789"`
	TaxRate           *string  `json:"tax_rate" binding:"omitempty" example:"12"`
	LogoUrl           *string  `json:"logo_url" binding:"omitempty" example:"https://imgur.com/234343"`
	Rounding          *string  `json:"rounding" binding:"omitempty" example:"nearest"`
	Currency          *string  `json:"currency" binding:"omitempty" example:"NGN"`
	Timezone          *string  `json:"timezone" binding:"omitempty" example:"UTC +1"`
	Language          *string  `json:"language" binding:"omitempty" example:"en"`
	LowStockThreshold *int     `json:"low_stock_threshold" binding:"omitempty" example:"5"`
	AllowOverselling  *bool    `json:"allow_overselling" binding:"omitempty" example:"false"`
	PaymentType       []string `json:"payment_type" binding:"omitempty" example:"cash,pos,room_charge,transfer"`
	Font              *string  `json:"font" binding:"omitempty"`
	PrimaryColor      *string  `json:"primary_color" binding:"omitempty"`
	Motto             *string  `json:"motto" binding:"omitempty"`
	Country           *string  `json:"country" binding:"omitempty" example:"Nigeria"`
}

type CreateBusinessResponse struct {
	Name              string   `json:"name" example:"Palmwineexpress hotels" binding:"required"`
	Email             string   `json:"email" binding:"omitempty" example:"admin@palmwinexpress.com"`
	Website           string   `json:"website" binding:"omitempty" example:"https://palmwinexpress.com"`
	TaxID             string   `json:"tax_id" binding:"omitempty" example:"123456789"`
	TaxRate           string   `json:"tax_rate" binding:"omitempty" example:"12"`
	LogoUrl           string   `json:"logo_url" binding:"omitempty" example:"https://imgur.com/234343"`
	Rounding          string   `json:"rounding" binding:"omitempty" example:"nearest"`
	Currency          string   `json:"currency" binding:"omitempty" example:"NGN"`
	Timezone          string   `json:"timezone" binding:"omitempty" example:"UTC +1"`
	Language          string   `json:"language" binding:"omitempty" example:"en"`
	LowStockThreshold int32    `json:"low_stock_threshold" binding:"omitempty" example:"5"`
	AllowOverselling  bool     `json:"allow_overselling" binding:"omitempty" example:"false"`
	PaymentType       []string `json:"payment_type" binding:"omitempty" example:"cash,pos,room_charge,transfer"`
	Font              string   `json:"font" binding:"omitempty"`
	PrimaryColor      string   `json:"primary_color" binding:"omitempty"`
	Motto             string   `json:"motto" binding:"omitempty"`
	Country           string   `json:"country" binding:"omitempty" example:"Nigeria"`
}

// CreateBusiness godoc
// @Summary Create business
// @Description Create a new business
// @Tags core
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param business body CreateBusinessParams true "Business details"
// @Success 201 {object} CreateBusinessResponse
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /core/business [post]
func (h *Handler) createBusiness(c *gin.Context) {
	var req CreateBusinessParams
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("error binding creating business request data: %v", err)
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	var params db.CreateBusinessParams
	err := copier.Copy(&params, &req)
	if err != nil {
		h.logger.Errorf("error copying business request data: %v", err)
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	business, err := h.service.CreateBusiness(c, params)
	if err != nil {
		h.logger.Errorf("error creating a business: %v", err)
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	// Convert []db.PaymentType to []string
	paymentTypes := make([]string, len(business.PaymentType))
	for i, pt := range business.PaymentType {
		paymentTypes[i] = string(pt)
	}

	utils.SuccessResponse(c, 201, "Business created", CreateBusinessResponse{
		Name: business.Name,

		Email:             business.Email.String,
		Website:           business.Website.String,
		TaxID:             business.TaxID.String,
		TaxRate:           business.TaxRate.String,
		LogoUrl:           business.LogoUrl.String,
		Rounding:          business.Rounding.String,
		Currency:          business.Currency.String,
		Timezone:          business.Timezone.String,
		Language:          business.Language.String,
		LowStockThreshold: business.LowStockThreshold.Int32,
		AllowOverselling:  business.AllowOverselling.Bool,
		PaymentType:       paymentTypes,
		Font:              business.Font.String,
		PrimaryColor:      business.PrimaryColor.String,
		Motto:             business.Motto.String,
		Country:           business.Country,
	})
}

// GetBusiness godoc
// @Summary Get business
// @Description Get a business
// @Tags core
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CreateBusinessResponse
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /core/business/:id [get]
func (h *Handler) getBusiness(c *gin.Context) {
	id := c.Param("id")
	bid, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Errorf("get business id str conv err: %v", err)
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	business, err := h.service.GetBusiness(c, int32(bid))
	if err != nil {
		h.logger.Errorf("error getting business with is %d: %v", bid, err)
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "get business successful", business)
}

func (h *Handler) updateBusiness(c *gin.Context) {
	// Implementation goes here
}

// DeleteBusiness godoc
// @Summary Delete business
// @Description Delete a new business
// @Tags core
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {string} string "business deleted"
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /core/business/:id [delete]
func (h *Handler) deleteBusiness(c *gin.Context) {
	id := c.Param("id")
	bid, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Errorf("get business id str conv err: %v", err)
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	err = h.service.DeleteBusiness(c, int32(bid))
	if err != nil {
		h.logger.Errorf("error deleteing business with is %d: %v", bid, err)
		utils.ErrorResponse(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, 200, "business deleted", nil)
}

// ListBusinesses godoc
// @Summary Get a list of businesses
// @Description Get a list of businesses
// @Tags core
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []CreateBusinessResponse
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /core/business/all [get]
func (h *Handler) listBusinesses(c *gin.Context) {
	businesses, err := h.service.ListBusinesses(c)
	if err != nil {
		h.logger.Errorf("error listing businesses: %v", err)
		utils.ErrorResponse(c, 500, err.Error())
		return
	}
	utils.SuccessResponse(c, 200, "A list of your businesses", businesses)
}

func (h *Handler) createBranch(c *gin.Context) {
	// Implementation goes here
}

func (h *Handler) getBranch(c *gin.Context) {
	// Implementation goes here
}

func (h *Handler) updateBranch(c *gin.Context) {
	// Implementation goes here
}

func (h *Handler) deleteBranch(c *gin.Context) {
	// Implementation goes here
}

func (h *Handler) listBranches(c *gin.Context) {
	// Implementation goes here
}
