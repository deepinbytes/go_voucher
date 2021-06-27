package controllers

import (
	"errors"
	"github.com/deepinbytes/go_voucher/services/offerservice"
	"github.com/deepinbytes/go_voucher/services/userservice"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/deepinbytes/go_voucher/domain/voucher"
	"github.com/deepinbytes/go_voucher/services/voucherservice"

	"github.com/gin-gonic/gin"
)

// UserInput represents login/register request body format
type VoucherInput struct {
	Code    string `json:"code"`
	UserID  uint   `json:"user_id"`
	OfferID uint   `json:"offer_id"`
}

// VoucherOutput represents returning user
type VoucherOutput struct {
	ID                 uint `json:"id"`
	IsUsed             bool
	ExpireTime         time.Time
	Code               string
	DiscountPercentage uint
	UsedAt             time.Time
	OfferID            uint
	UserID             uint
}

type RedeemVoucherInput struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

// UserController interface
type VoucherController interface {
	Create(*gin.Context)
	Redeem(*gin.Context)
	GetByID(*gin.Context)
	Update(*gin.Context)
}

type voucherController struct {
	voucherSvc voucherservice.VoucherService
	usrSvc     userservice.UserService
	offerSvc   offerservice.OfferService
}

// NewUserController instantiates User Controller
func NewVoucherController(
	voucherSvc voucherservice.VoucherService,
	usrSvc userservice.UserService,
	offerSvc offerservice.OfferService) VoucherController {
	return &voucherController{
		voucherSvc: voucherSvc,
		usrSvc:     usrSvc,
		offerSvc:   offerSvc,
	}
}

// @Summary Create New Voucher
// @Produce  json
// @Param code body string true "Code"
// @Param user_id body string true "UserID"
// @Param offer_id body string true "OfferID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/voucher/create [post]
func (ctl *voucherController) Create(c *gin.Context) {
	// Read user input
	var voucherGenerateInput VoucherInput
	if err := c.ShouldBindJSON(&voucherGenerateInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	u := ctl.inputToVoucher(voucherGenerateInput)

	// Create voucher
	if err := ctl.voucherSvc.Create(&u); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

}

// @Summary Get voucher info of given id
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/voucher/{id} [get]
func (ctl *voucherController) GetByID(c *gin.Context) {
	id, err := ctl.getVoucherID(c.Param(("id")))
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	voucher, err := ctl.voucherSvc.GetByID(id)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			HTTPRes(c, http.StatusNotFound, es, nil)
			return
		}
		HTTPRes(c, http.StatusInternalServerError, es, nil)
		return
	}
	voucherOutput := ctl.mapToVoucherOutput(voucher)
	HTTPRes(c, http.StatusOK, "ok", voucherOutput)
}

// @Summary Redeem Voucher
// @Produce  json
// @Param email body string true "Email"
// @Param code body string false "Code"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/voucher/redeem [post]
func (ctl *voucherController) Redeem(c *gin.Context) {

	var redeemVoucherInput RedeemVoucherInput
	if err := c.ShouldBindJSON(&redeemVoucherInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Retrieve voucher given the code
	voucher, err := ctl.voucherSvc.UseCode(redeemVoucherInput.Code)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), "Invalid Voucher")
		return
	}
	offer, err := ctl.offerSvc.GetByID(voucher.OfferID)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), "Offer Not Available Anymore")
		return
	}

	// Retrieve user for given voucher
	user, err := ctl.usrSvc.GetByID(voucher.UserID)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), "Invalid User")
		return
	}
	if redeemVoucherInput.Email != user.Email {
		HTTPRes(c, http.StatusInternalServerError, "", "Code not valid for this user")
		return
	}

	if voucher.IsUsed {
		HTTPRes(c, http.StatusOK, "Used Voucher", "")
		return
	}

	if time.Now().After(voucher.ExpireTime) {
		HTTPRes(c, http.StatusOK, "Expired Voucher", "")
		return
	}

	// Updates voucher record
	voucher.UsedAt = time.Now()
	voucher.IsUsed = true

	if err := ctl.voucherSvc.Update(voucher); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Response
	voucherOutput := ctl.mapToVoucherOutput(voucher)
	voucherOutput.DiscountPercentage = offer.DiscountPercentage
	HTTPRes(c, http.StatusOK, "ok", voucherOutput)
}

// @Summary Update voucher info
// @Produce  json
// @Param voucher_id body string true "ID"
// @Param code body string false "Code"
// @Param user_id body string false "UserID"
// @Param offer_id body string false "OfferID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/voucher/update [post]
func (ctl *voucherController) Update(c *gin.Context) {
	// Get voucher id from context
	id, exists := c.Get("voucher_id")
	if exists == false {
		HTTPRes(c, http.StatusBadRequest, "Invalid Voucher ID", nil)
		return
	}

	// Retrieve voucher given id
	voucher, err := ctl.voucherSvc.GetByID(id.(uint))
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Read voucher input
	var voucherInput VoucherInput
	if err := c.ShouldBindJSON(&voucherInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Check voucher
	if voucher.ID != id {
		HTTPRes(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Update voucher record
	voucher.UsedAt = time.Now()
	voucher.IsUsed = true

	if err := ctl.voucherSvc.Update(voucher); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Response
	voucherOutput := ctl.mapToVoucherOutput(voucher)
	HTTPRes(c, http.StatusOK, "ok", voucherOutput)
}

/*******************************/
//       PRIVATE METHODS
/*******************************/

func (ctl *voucherController) getVoucherID(voucherIDParam string) (uint, error) {
	voucherID, err := strconv.Atoi(voucherIDParam)
	if err != nil {
		return 0, errors.New("voucher id should be a number")
	}
	return uint(voucherID), nil
}

func (ctl *voucherController) inputToVoucher(input VoucherInput) voucher.Voucher {

	return voucher.Voucher{
		Code:    input.Code,
		UserID:  input.UserID,
		OfferID: input.OfferID,
	}
}

func (ctl *voucherController) mapToVoucherOutput(u *voucher.Voucher) *VoucherOutput {
	return &VoucherOutput{
		ID:         u.ID,
		Code:       u.Code,
		IsUsed:     u.IsUsed,
		UsedAt:     u.UsedAt,
		ExpireTime: u.ExpireTime,
		UserID:     u.UserID,
		OfferID:    u.OfferID,
	}
}
