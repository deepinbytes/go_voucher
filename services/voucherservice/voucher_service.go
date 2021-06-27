package voucherservice

import (
	"errors"
	"github.com/deepinbytes/go_voucher/domain/voucher"
	"github.com/deepinbytes/go_voucher/repositories/voucherrepo"
)

// voucherService interface
type VoucherService interface {
	GetByID(id uint) (*voucher.Voucher, error)
	UseCode(code string) (*voucher.Voucher, error)
	Create(*voucher.Voucher) error
	Update(*voucher.Voucher) error
}

type voucherService struct {
	Repo voucherrepo.Repo
}

// NewVoucherService will instantiate Voucher Service
func NewVoucherService(
	repo voucherrepo.Repo,
) VoucherService {

	return &voucherService{
		Repo: repo,
	}
}

func (vs *voucherService) GetByID(id uint) (*voucher.Voucher, error) {
	if id == 0 {
		return nil, errors.New("id param is required")
	}
	voucher, err := vs.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return voucher, nil
}

func (vs *voucherService) UseCode(code string) (*voucher.Voucher, error) {
	if code == "" {
		return nil, errors.New("Code(string) is required")
	}
	voucher, err := vs.Repo.UseCode(code)
	if err != nil {
		return nil, err
	}
	return voucher, nil
}

func (vs *voucherService) Create(voucher *voucher.Voucher) error {
	return vs.Repo.Create(voucher)
}

func (vs *voucherService) Update(voucher *voucher.Voucher) error {
	return vs.Repo.Update(voucher)
}
