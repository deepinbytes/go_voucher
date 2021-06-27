package voucherservice

import (
	"github.com/deepinbytes/go_voucher/domain/voucher"
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

func (repo *repoMock) GetByID(id uint) (*voucher.Voucher, error) {
	args := repo.Called(id)
	return args.Get(0).(*voucher.Voucher), args.Error(1)
}

func (repo *repoMock) UseCode(code string) (*voucher.Voucher, error) {
	args := repo.Called(code)
	return args.Get(0).(*voucher.Voucher), args.Error(1)
}

func (repo *repoMock) Create(voucher *voucher.Voucher) error {
	args := repo.Called(voucher)
	return args.Error(0)
}

func (repo *repoMock) Update(voucher *voucher.Voucher) error {
	args := repo.Called(voucher)
	return args.Error(0)
}
