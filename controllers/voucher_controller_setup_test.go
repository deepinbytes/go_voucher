package controllers

import (
	"errors"

	"github.com/deepinbytes/go_voucher/domain/voucher"
	"github.com/jinzhu/gorm"
)

type voucherSvc struct{}

var voucher1 = &voucher.Voucher{
	Model: gorm.Model{ID: uint(1)},
	Code:  "TEST1",
}

var voucher2 = &voucher.Voucher{
	Model:   gorm.Model{ID: uint(1)},
	Code:    "TEST2",
	UserID:  1,
	OfferID: 1,
}

func (vs *voucherSvc) GetByID(id uint) (*voucher.Voucher, error) {
	if id >= uint(100) {
		return nil, errors.New("Ugh")
	}
	if id < uint(1) {
		return nil, errors.New("Nop")
	}
	if id >= uint(10) {
		return nil, errors.New("Record not found")
	}
	return voucher1, nil
}

func (vs *voucherSvc) UseCode(code string) (*voucher.Voucher, error) {
	if code == "non_existent_code" {
		return nil, errors.New("Nop")
	}
	if code == voucher1.Code {
		return voucher1, nil
	}
	return voucher2, nil
}

func (vs *voucherSvc) Create(voucher *voucher.Voucher) error {
	if voucher.Code == "existing_code" {
		return errors.New("Nop")
	}
	return nil
}

func (vs *voucherSvc) Update(voucher *voucher.Voucher) error {
	if voucher.Code == "non_existing_code" {
		return errors.New("Nop")
	}
	return nil
}
