package voucherrepo

import (
	"database/sql/driver"
	"errors"
	"github.com/deepinbytes/go_voucher/domain/voucher"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("can't create sqlmock: %s", err)
	}

	gormDB, gerr := gorm.Open("postgres", db)
	if gerr != nil {
		log.Fatalf("can't open gorm connection: %s", err)
	}
	gormDB.LogMode(true)
	return gormDB, mock
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestRedeem(t *testing.T) {
	gormDB, mock := setupDB()
	defer gormDB.Close()

	t.Run("Redeem a Voucher", func(t *testing.T) {
		expected := &voucher.Voucher{
			Code: "aliceSDS",
		}

		u := NewVoucherRepo(gormDB)

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM "vouchers" WHERE "vouchers"."deleted_at" IS NULL AND ((vouchers.code = $1)) ORDER BY "vouchers"."id" ASC LIMIT 1`)).
			WithArgs("aliceSDS").
			WillReturnRows(
				sqlmock.NewRows([]string{"code"}).
					AddRow("aliceSDS"))

		mock.ExpectCommit()

		result, err := u.UseCode("aliceSDS")
		assert.EqualValues(t, expected, result)
		assert.Nil(t, err)
	})
}

func TestRedeemFail(t *testing.T) {
	gormDB, mock := setupDB()
	defer gormDB.Close()

	t.Run("Redeem a user fails", func(t *testing.T) {
		exp := errors.New("oops")

		u := NewVoucherRepo(gormDB)

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM "vouchers" WHERE "vouchers"."deleted_at" IS NULL AND ((vouchers.code = $1)) ORDER BY "vouchers"."id" ASC LIMIT 1`)).
			WithArgs("aliceSDS").
			WillReturnError(exp)

		mock.ExpectCommit()

		_, err := u.UseCode("aliceSDS")
		assert.NotNil(t, err)
		assert.EqualValues(t, exp, err)

	})
}

func TestCreate(t *testing.T) {
	gormDB, mock := setupDB()
	defer gormDB.Close()

	t.Run("Create a Voucher", func(t *testing.T) {
		v := &voucher.Voucher{
			Model:   gorm.Model{},
			UsedAt:  time.Time{},
			IsUsed:  false,
			Code:    "aliceSDS",
			OfferID: 1,
			UserID:  1,
		}

		u := NewVoucherRepo(gormDB)

		mock.ExpectBegin()

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`INSERT INTO "vouchers" ("created_at","updated_at","deleted_at","used_at","code","offer_id","user_id","expire_time") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "vouchers"."id"`)).
			WithArgs(AnyTime{}, AnyTime{}, nil, AnyTime{}, "aliceSDS", 1, 1, AnyTime{}).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(1))

		mock.ExpectCommit()

		err := u.Create(v)
		assert.Nil(t, err)
	})

	t.Run("Create a user fails", func(t *testing.T) {
		exp := errors.New("oops")
		v := &voucher.Voucher{
			Model:   gorm.Model{},
			UsedAt:  time.Time{},
			IsUsed:  false,
			Code:    "aliceSDS",
			OfferID: 1,
			UserID:  1,
		}

		u := NewVoucherRepo(gormDB)

		mock.ExpectBegin()

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`INSERT INTO "vouchers" ("created_at","updated_at","deleted_at","used_at","code","offer_id","user_id","expire_time") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "vouchers"."id"`)).
			WithArgs(AnyTime{}, AnyTime{}, nil, AnyTime{}, "aliceSDS", 1, 1, AnyTime{}).
			WillReturnError(exp)

		mock.ExpectCommit()

		err := u.Create(v)
		assert.NotNil(t, err)
		assert.EqualValues(t, exp, err)
	})
}
