package core

import (
	db "herp/db/sqlc"
	"herp/internal/auth"
	"herp/internal/config"
	"herp/internal/utils"
	"herp/pkg/monitoring/logging"

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
		business.GET("", auth.PermissionMiddleware(authSvc, "core:view_business"), h.listBusinesses)
	}
}

type CreateBusinessParams struct {
	Name              string   `json:"name" example:"Palmwineexpress hotels" binding:"required"`
	AddressOne        string   `json:"address_one" example:"32 Ander avenue" binding:"required"`
	AddresTwo         *string  `json:"addres_two" binding:"omitempty" example:"2nd floor"`
	Country           string   `json:"country" example:"Nigeria" binding:"required"`
	Phone             *string  `json:"phone" binding:"omitempty" example:"+2348123456789"`
	Email             *string  `json:"email" binding:"omitempty" example:"admin@palmwinexpress.com"`
	Website           *string  `json:"website" binding:"omitempty" example:"https://palmwinexpress.com"`
	City              *string  `json:"city" example:"Aba" binding:"omitempty"`
	State             *string  `json:"state" example:"Abia" binding:"omitempty"`
	ZipCode           *string  `json:"zip_code" binding:"omitempty" example:"23432"`
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
}

type CreateBusinessResponse struct {
	Name              string   `json:"name" example:"Palmwineexpress hotels" binding:"required"`
	AddressOne        string   `json:"address_one" example:"32 Ander avenue" binding:"required"`
	AddresTwo         string  `json:"address_two" binding:"omitempty" example:"2nd floor"`
	Country           string   `json:"country" example:"Nigeria" binding:"required"`
	Phone             string  `json:"phone" binding:"omitempty" example:"+2348123456789"`
	Email             string  `json:"email" binding:"omitempty" example:"admin@palmwinexpress.com"`
	Website           string  `json:"website" binding:"omitempty" example:"https://palmwinexpress.com"`
	City              string  `json:"city" example:"Aba" binding:"omitempty"`
	State             string  `json:"state" example:"Abia" binding:"omitempty"`
	ZipCode           string  `json:"zip_code" binding:"omitempty" example:"23432"`
	TaxID             string  `json:"tax_id" binding:"omitempty" example:"123456789"`
	TaxRate           string  `json:"tax_rate" binding:"omitempty" example:"12"`
	LogoUrl           string  `json:"logo_url" binding:"omitempty" example:"https://imgur.com/234343"`
	Rounding          string  `json:"rounding" binding:"omitempty" example:"nearest"`
	Currency          string  `json:"currency" binding:"omitempty" example:"NGN"`
	Timezone          string  `json:"timezone" binding:"omitempty" example:"UTC +1"`
	Language          string  `json:"language" binding:"omitempty" example:"en"`
	LowStockThreshold int32     `json:"low_stock_threshold" binding:"omitempty" example:"5"`
	AllowOverselling  bool    `json:"allow_overselling" binding:"omitempty" example:"false"`
	PaymentType       []string `json:"payment_type" binding:"omitempty" example:"cash,pos,room_charge,transfer"`
	Font              string  `json:"font" binding:"omitempty"`
	PrimaryColor      string  `json:"primary_color" binding:"omitempty"`
	Motto             string  `json:"motto" binding:"omitempty"`
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
		Name:              business.Name,
		AddressOne:        business.AddressOne,
		AddresTwo:         business.AddresTwo.String,
		Country:           business.Country,
		Phone:             business.Phone.String,
		Email:             business.Email.String,
		Website:           business.Website.String,
		City:              business.City.String,
		State:             business.State.String,
		ZipCode:           business.ZipCode.String,
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
	})
}

func (h *Handler) getBusiness(c *gin.Context) {
	// Implementation goes here
}

func (h *Handler) updateBusiness(c *gin.Context) {
	// Implementation goes here
}

func (h *Handler) deleteBusiness(c *gin.Context) {
	// Implementation goes here
}

func (h *Handler) listBusinesses(c *gin.Context) {
	// Implementation goes here
}
