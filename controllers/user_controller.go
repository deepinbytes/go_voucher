package controllers

import (
	"errors"
	"github.com/deepinbytes/go_voucher/domain/voucher"
	"net/http"
	"strconv"
	"strings"

	"github.com/deepinbytes/go_voucher/domain/user"
	"github.com/deepinbytes/go_voucher/services/userservice"

	"github.com/gin-gonic/gin"
)

// UserInput represents register request body format
type UserInput struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// UserOutput represents returning user
type UserOutput struct {
	ID        uint              `json:"id"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Email     string            `json:"email"`
	Vouchers  []voucher.Voucher `json:"vouchers"`
}

// UserUpdateInput represents updating profile request body format
type UserUpdateInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// UserController interface
type UserController interface {
	Register(*gin.Context)
	GetByID(*gin.Context)
	GetByEmail(*gin.Context)
	ListUsers(*gin.Context)
	Update(*gin.Context)
}

type userController struct {
	us userservice.UserService
}

// NewUserController instantiates User Controller
func NewUserController(
	us userservice.UserService) UserController {
	return &userController{
		us: us,
	}
}

// @Summary Register new user
// @Produce  json
// @Param email body string true "Email"
// @Param firstName body string true "FirstName"
// @Param lastName body string true "LastName"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/register [post]
func (ctl *userController) Register(c *gin.Context) {
	// Read user input
	var userInput UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	u := ctl.inputToUser(userInput)

	// Create user
	if err := ctl.us.Create(&u); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	userOutput := ctl.mapToUserOutput(&u)
	HTTPRes(c, http.StatusOK, "ok", userOutput)
}

// @Summary Get user info of given id
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/user/{id} [get]
func (ctl *userController) GetByID(c *gin.Context) {
	id, err := ctl.getUserID(c.Param(("id")))
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	user, err := ctl.us.GetByID(id)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			HTTPRes(c, http.StatusNotFound, es, nil)
			return
		}
		HTTPRes(c, http.StatusInternalServerError, es, nil)
		return
	}
	userOutput := ctl.mapToUserOutput(user)
	HTTPRes(c, http.StatusOK, "ok", userOutput)
}

// @Summary Get user info using given email
// @Produce  json
// @Param email path int true "ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/user/{email} [get]
func (ctl *userController) GetByEmail(c *gin.Context) {
	email := c.Param(("email"))
	user, err := ctl.us.GetByEmail(email)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	userOutput := ctl.mapToUserOutput(user)
	HTTPRes(c, http.StatusOK, "ok", userOutput)
}

// @Summary Get users
// @Produce  json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/list_users [get]
func (ctl *userController) ListUsers(c *gin.Context) {

	users, err := ctl.us.ListAll()
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	HTTPRes(c, http.StatusOK, "ok", users)
}

// @Summary Update account info
// @Produce  json
// @Param user_id body string true "ID"
// @Param email body string true "Email"
// @Param firstName body string false "First Name"
// @Param lastName body string false "Last Name"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/user/update [put]
func (ctl *userController) Update(c *gin.Context) {
	// Get user id from context
	id, exists := c.Get("user_id")
	if exists == false {
		HTTPRes(c, http.StatusBadRequest, "Invalid User ID", nil)
		return
	}

	// Retrieve user given id
	user, err := ctl.us.GetByID(id.(uint))
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Read user input
	var userInput UserUpdateInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Check user
	if user.ID != id {
		HTTPRes(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Update user record
	user.FirstName = userInput.FirstName
	user.LastName = userInput.LastName
	user.Email = userInput.Email
	if err := ctl.us.Update(user); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Response
	userOutput := ctl.mapToUserOutput(user)
	HTTPRes(c, http.StatusOK, "ok", userOutput)
}

/*******************************/
//       PRIVATE METHODS
/*******************************/

func (ctl *userController) getUserID(userIDParam string) (uint, error) {
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return 0, errors.New("user id should be a number")
	}
	return uint(userID), nil
}

func (ctl *userController) inputToUser(input UserInput) user.User {
	return user.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}
}

func (ctl *userController) mapToUserOutput(u *user.User) *UserOutput {
	return &UserOutput{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Vouchers:  u.Voucher,
	}
}

// Issue token and return user
func (ctl *userController) login(c *gin.Context, u *user.User) error {
	userOutput := ctl.mapToUserOutput(u)
	out := gin.H{"user": userOutput}
	HTTPRes(c, http.StatusOK, "ok", out)
	return nil
}
