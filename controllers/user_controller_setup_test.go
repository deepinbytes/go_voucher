package controllers

import (
	"errors"

	"github.com/deepinbytes/go_voucher/domain/user"
	"github.com/jinzhu/gorm"
)

type userSvc struct{}

var alice = &user.User{
	Model:     gorm.Model{ID: uint(1)},
	Email:     "alice@cc.cc",
	FirstName: "",
	LastName:  "",
}

var david = &user.User{
	Model:     gorm.Model{ID: uint(1)},
	Email:     "david@cc.cc",
	FirstName: "",
	LastName:  "",
}

func (us *userSvc) GetByID(id uint) (*user.User, error) {
	if id >= uint(100) {
		return nil, errors.New("Ugh")
	}
	if id < uint(1) {
		return nil, errors.New("Nop")
	}
	if id >= uint(10) {
		return nil, errors.New("Record not found")
	}
	return alice, nil
}

func (us *userSvc) GetByEmail(email string) (*user.User, error) {
	if email == "bob@cc.cc" {
		return nil, errors.New("Nop")
	}
	if email == david.Email {
		return david, nil
	}
	return alice, nil
}

func (us *userSvc) Create(user *user.User) error {
	if user.Email == "bob@cc.cc" {
		return errors.New("Nop")
	}
	return nil
}

func (us *userSvc) Update(user *user.User) error {
	if user.Email == "bob@cc.cc" {
		return errors.New("Nop")
	}
	return nil
}

func (us *userSvc) ListAll() ([]*user.User, error) {
	var x []*user.User
	return x, nil

}
