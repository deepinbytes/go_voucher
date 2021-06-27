package offerservice

import (
	"github.com/deepinbytes/go_voucher/domain/offer"
	"github.com/stretchr/testify/mock"
)

var (
	testID10  = uint(10)
	testID100 = uint(100)
	testName  = "test"
)

type repoMock struct {
	mock.Mock
}

func (repo *repoMock) GetByID(id uint) (*offer.Offer, error) {
	args := repo.Called(id)
	return args.Get(0).(*offer.Offer), args.Error(1)
}

func (repo *repoMock) GetByName(email string) (*offer.Offer, error) {
	args := repo.Called(email)
	return args.Get(0).(*offer.Offer), args.Error(1)
}

func (repo *repoMock) Create(offer *offer.Offer) error {
	args := repo.Called(offer)
	return args.Error(0)
}

func (repo *repoMock) Update(offer *offer.Offer) error {
	args := repo.Called(offer)
	return args.Error(0)
}
