package offerrepo

import (
	"github.com/deepinbytes/go_voucher/domain/offer"

	"github.com/jinzhu/gorm"
)

// Repo interface
type Repo interface {
	GetByID(id uint) (*offer.Offer, error)
	GetByName(name string) (*offer.Offer, error)
	Create(offer *offer.Offer) error
	Update(offer *offer.Offer) error
}

type offerRepo struct {
	db *gorm.DB
}

// NewUserRepo will instantiate User Repository
func NewOfferRepo(db *gorm.DB) Repo {
	return &offerRepo{
		db: db,
	}
}

func (u *offerRepo) GetByID(id uint) (*offer.Offer, error) {
	var offer offer.Offer
	if err := u.db.First(&offer, id).Error; err != nil {
		return nil, err
	}
	return &offer, nil
}

func (u *offerRepo) GetByName(name string) (*offer.Offer, error) {
	var offer offer.Offer
	if err := u.db.Where("name = ?", name).First(&offer).Error; err != nil {
		return nil, err
	}
	return &offer, nil
}

func (u *offerRepo) Create(offer *offer.Offer) error {
	return u.db.Create(offer).Error
}

func (u *offerRepo) Update(offer *offer.Offer) error {
	return u.db.Save(offer).Error
}
