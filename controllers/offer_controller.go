package controllers

import (
	"errors"
	"github.com/deepinbytes/go_voucher/domain/voucher"
	"github.com/deepinbytes/go_voucher/services/userservice"
	"github.com/deepinbytes/go_voucher/services/voucherservice"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/deepinbytes/go_voucher/domain/offer"
	"github.com/deepinbytes/go_voucher/services/offerservice"

	"github.com/gin-gonic/gin"
)

// UserInput represents login/register request body format
type OfferInput struct {
	Name               string `json:"name"`
	DiscountPercentage uint   `json:"discount_percentage"`
}

type OfferGenerateVoucherInput struct {
	Name       string `json:"name"`
	ExpiryTime uint   `json:"expiry_time"`
}

// UserOutput represents returning user
type OfferOutput struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name"`
	DiscountPercentage uint   `json:"discount_percentage"`
}

// UserUpdateInput represents updating profile request body format
type OfferUpdateInput struct {
	DiscountPercentage uint   `json:"discount_percentage"`
	Name               string `json:"name"`
}

// UserController interface
type OfferController interface {
	Create(*gin.Context)
	GetByID(*gin.Context)
	Update(*gin.Context)
	GenerateVouchers(*gin.Context)
}

type offerController struct {
	offerSvc offerservice.OfferService
	usrSvc   userservice.UserService
	vouchSvc voucherservice.VoucherService
}

// @Summary Generates vouchers for all the users given offer name
// @Produce  json
// @Param name body string true "Name"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/offer/generate_vouchers [post]
func (ctl *offerController) GenerateVouchers(c *gin.Context) {

	var generateVoucherInput OfferGenerateVoucherInput
	if err := c.ShouldBindJSON(&generateVoucherInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	offer, err := ctl.offerSvc.GetByName(generateVoucherInput.Name)
	if err != nil {
		HTTPRes(c, http.StatusNotFound, err.Error(), "Offer not found")
		return
	}

	users, _ := ctl.usrSvc.ListAll()
	//var v []*voucher.Voucher
	for _, t := range users {
		t := &voucher.Voucher{
			Model:      gorm.Model{},
			UsedAt:     time.Time{},
			IsUsed:     false,
			Code:       randSeq(8),
			OfferID:    offer.ID,
			UserID:     t.ID,
			ExpireTime: time.Now().Add(time.Hour * (time.Duration(generateVoucherInput.ExpiryTime) * 24)),
		}
		ctl.vouchSvc.Create(t)
	}
	//p := ctl.vouchSvc.BulkCreate(v)
	HTTPRes(c, http.StatusOK, "ok", "Generated")

}

// NewUserController instantiates User Controller
func NewOfferController(
	us offerservice.OfferService,
	usrSvc userservice.UserService,
	voucherSvc voucherservice.VoucherService) OfferController {
	return &offerController{
		offerSvc: us,
		usrSvc:   usrSvc,
		vouchSvc: voucherSvc,
	}
}

// @Summary Creates new offer
// @Produce  json
// @Param name body string true "Name"
// @Param discount_percentage body string true "DiscountPercentage"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/offer/create [post]
func (ctl *offerController) Create(c *gin.Context) {
	// Read user input
	var offerInput OfferInput
	if err := c.ShouldBindJSON(&offerInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	u := ctl.inputToUser(offerInput)

	// Create user
	if err := ctl.offerSvc.Create(&u); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

}

// @Summary Get offer info of given id
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/offer/{id} [get]
func (ctl *offerController) GetByID(c *gin.Context) {
	id, err := ctl.getOfferID(c.Param(("id")))
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	offer, err := ctl.offerSvc.GetByID(id)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			HTTPRes(c, http.StatusNotFound, es, nil)
			return
		}
		HTTPRes(c, http.StatusInternalServerError, es, nil)
		return
	}
	offerOutput := ctl.mapToOfferOutput(offer)
	HTTPRes(c, http.StatusOK, "ok", offerOutput)
}

// @Summary Update offer info
// @Produce  json
// @Param offer_id body string true "ID"
// @Param name body string false "Name"
// @Param discount_percentage body string false "DiscountPercentage"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/offer/update [post]
func (ctl *offerController) Update(c *gin.Context) {
	// Get user id from context
	id, exists := c.Get("offer_id")
	if exists == false {
		HTTPRes(c, http.StatusBadRequest, "Invalid Offer ID", nil)
		return
	}

	// Retrieve offer given id
	offer, err := ctl.offerSvc.GetByID(id.(uint))
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Read offer input
	var offerInput OfferUpdateInput
	if err := c.ShouldBindJSON(&offerInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Check offer
	if offer.ID != id {
		HTTPRes(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Update offer record
	offer.DiscountPercentage = offerInput.DiscountPercentage
	offer.Name = offerInput.Name

	if err := ctl.offerSvc.Update(offer); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Response
	offerOutput := ctl.mapToOfferOutput(offer)
	HTTPRes(c, http.StatusOK, "ok", offerOutput)
}

/*******************************/
//       PRIVATE METHODS
/*******************************/

func (ctl *offerController) getOfferID(offerIDParam string) (uint, error) {
	offerID, err := strconv.Atoi(offerIDParam)
	if err != nil {
		return 0, errors.New("offer id should be a number")
	}
	return uint(offerID), nil
}

func (ctl *offerController) inputToUser(input OfferInput) offer.Offer {
	return offer.Offer{
		Name:               input.Name,
		DiscountPercentage: input.DiscountPercentage,
	}
}

func (ctl *offerController) mapToOfferOutput(u *offer.Offer) *OfferOutput {
	return &OfferOutput{
		ID:                 u.ID,
		Name:               u.Name,
		DiscountPercentage: u.DiscountPercentage,
	}
}
