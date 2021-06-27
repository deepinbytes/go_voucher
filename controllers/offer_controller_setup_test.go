package controllers

import (
	"errors"
	"github.com/deepinbytes/go_voucher/domain/offer"
	"github.com/jinzhu/gorm"
)

type offerSvc struct{}

var of1 = &offer.Offer{
	Model:              gorm.Model{ID: uint(1)},
	Name:               "offer1",
	DiscountPercentage: 98,
}

var of2 = &offer.Offer{
	Model:              gorm.Model{ID: uint(1)},
	Name:               "offer2",
	DiscountPercentage: 27,
}

func (os *offerSvc) GetByID(id uint) (*offer.Offer, error) {
	if id >= uint(100) {
		return nil, errors.New("Ugh")
	}
	if id < uint(1) {
		return nil, errors.New("Nop")
	}
	if id >= uint(10) {
		return nil, errors.New("Record not found")
	}
	return of1, nil
}

func (os *offerSvc) GetByName(name string) (*offer.Offer, error) {
	if name == "non_existent_offer" {
		return nil, errors.New("Nop")
	}
	if name == of1.Name {
		return of1, nil
	}
	return of1, nil
}

func (os *offerSvc) Create(offer *offer.Offer) error {
	if offer.Name == "new_offer" {
		return errors.New("Nop")
	}
	return nil
}

func (os *offerSvc) Update(offer *offer.Offer) error {
	if offer.Name == "non_existing_offer" {
		return errors.New("Nop")
	}
	return nil
}

func (os *offerSvc) ListAll() ([]*offer.Offer, error) {
	var x []*offer.Offer
	return x, nil

}
