package offerservice

import (
	"errors"
	"github.com/deepinbytes/go_voucher/domain/offer"
	"github.com/deepinbytes/go_voucher/repositories/offerrepo"
)

// OfferService interface
type OfferService interface {
	GetByID(id uint) (*offer.Offer, error)
	GetByName(name string) (*offer.Offer, error)
	Create(*offer.Offer) error
	Update(*offer.Offer) error
}

type offerService struct {
	Repo offerrepo.Repo
}

// NewOfferService will instantiate User Service
func NewOfferService(
	repo offerrepo.Repo,
) OfferService {

	return &offerService{
		Repo: repo,
	}
}

func (os *offerService) GetByID(id uint) (*offer.Offer, error) {
	if id == 0 {
		return nil, errors.New("id param is required")
	}
	offer, err := os.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return offer, nil
}

func (os *offerService) GetByName(name string) (*offer.Offer, error) {
	if name == "" {
		return nil, errors.New("Name(string) is required")
	}
	user, err := os.Repo.GetByName(name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (os *offerService) Create(offer *offer.Offer) error {
	return os.Repo.Create(offer)
}

func (os *offerService) Update(offer *offer.Offer) error {
	return os.Repo.Update(offer)
}
