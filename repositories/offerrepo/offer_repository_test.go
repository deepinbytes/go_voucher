package offerrepo

import (
	"database/sql/driver"
	"errors"
	"github.com/deepinbytes/go_voucher/domain/offer"
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

func TestGetByID(t *testing.T) {
	gormDB, mock := setupDB()
	defer gormDB.Close()

	t.Run("Get a offer", func(t *testing.T) {
		expected := &offer.Offer{
			Name: "TEST",
		}

		u := NewOfferRepo(gormDB)

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM "offers"  WHERE "offers"."deleted_at" IS NULL AND (("offers"."id" = 100)) ORDER BY "offers"."id" ASC LIMIT 1 `)).
			WillReturnRows(
				sqlmock.NewRows([]string{"name"}).
					AddRow("TEST"))

		result, err := u.GetByID(100)

		assert.EqualValues(t, expected, result)
		assert.Nil(t, err)
	})

	t.Run("Error occurs", func(t *testing.T) {
		expected := errors.New("Nop")

		u := NewOfferRepo(gormDB)

		mock.
			ExpectQuery(
				regexp.QuoteMeta(`SELECT * FROM "offers" WHERE "offers"."deleted_at" IS NULL AND (("offers"."id" = 100)) ORDER BY "offers"."id" ASC LIMIT 1`)).
			WillReturnError(expected)

		result, err := u.GetByID(100)

		assert.EqualValues(t, expected, err)
		assert.Nil(t, result)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		expected := errors.New("record not found")

		u := NewOfferRepo(gormDB)

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM "offers" WHERE "offers"."deleted_at" IS NULL AND (("offers"."id" = 100)) ORDER BY "offers"."id" ASC LIMIT 1`)).
			WillReturnRows(
				sqlmock.NewRows([]string{}))

		result, err := u.GetByID(100)

		assert.EqualValues(t, expected, err)
		assert.Nil(t, result)
	})
}

func TestGetByName(t *testing.T) {
	gormDB, mock := setupDB()
	defer gormDB.Close()

	t.Run("Get a offer", func(t *testing.T) {
		expected := &offer.Offer{
			Name: "TEST",
		}

		u := NewOfferRepo(gormDB)
		sqlStr := `SELECT * FROM "offers" WHERE "offers"."deleted_at" IS NULL AND ((name = $1)) ORDER BY "offers"."id" ASC LIMIT 1`

		mock.
			ExpectQuery(regexp.QuoteMeta(sqlStr)).
			WithArgs("TEST").
			WillReturnRows(
				sqlmock.NewRows([]string{"name"}).
					AddRow("TEST"))

		result, err := u.GetByName("TEST")

		assert.EqualValues(t, expected, result)
		assert.Nil(t, err)
	})

	t.Run("Error occurs", func(t *testing.T) {
		expected := errors.New("Nop")

		u := NewOfferRepo(gormDB)
		sqlStr := `SELECT * FROM "offers" WHERE "offers"."deleted_at" IS NULL AND ((name = $1)) ORDER BY "offers"."id" ASC LIMIT 1`

		mock.
			ExpectQuery(regexp.QuoteMeta(sqlStr)).
			WithArgs("TEST").
			WillReturnError(expected)

		result, err := u.GetByName("TEST")

		assert.EqualValues(t, expected, err)
		assert.Nil(t, result)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		expected := errors.New("record not found")

		u := NewOfferRepo(gormDB)
		sqlStr := `SELECT * FROM "offers" WHERE "offers"."deleted_at" IS NULL AND ((name = $1)) ORDER BY "offers"."id" ASC LIMIT 1`

		mock.
			ExpectQuery(regexp.QuoteMeta(sqlStr)).
			WithArgs("TEST").
			WillReturnRows(
				sqlmock.NewRows([]string{}))

		result, err := u.GetByName("TEST")

		assert.EqualValues(t, expected, err)
		assert.Nil(t, result)
	})
}

func TestCreate(t *testing.T) {
	gormDB, mock := setupDB()
	defer gormDB.Close()

	t.Run("Create an offer", func(t *testing.T) {
		offer := &offer.Offer{
			Name:               "TEST",
			DiscountPercentage: 24,
		}

		u := NewOfferRepo(gormDB)

		mock.ExpectBegin()

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`INSERT INTO "offers" ("created_at","updated_at","deleted_at","name","discount_percentage") VALUES ($1,$2,$3,$4,$5) RETURNING "offers"."id`)).
			WithArgs(AnyTime{}, AnyTime{}, nil, "TEST", 24).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(1))

		mock.ExpectCommit()

		err := u.Create(offer)
		assert.Nil(t, err)
	})

	t.Run("Create a offer fails", func(t *testing.T) {
		exp := errors.New("oops")
		offer := &offer.Offer{
			Name:               "TEST",
			DiscountPercentage: 24,
		}

		u := NewOfferRepo(gormDB)

		mock.ExpectBegin()

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`INSERT INTO "offers" ("created_at","updated_at","deleted_at","name","discount_percentage") VALUES ($1,$2,$3,$4,$5) RETURNING "offers"."id`)).
			WithArgs(AnyTime{}, AnyTime{}, nil, "TEST", 24).
			WillReturnError(exp)

		mock.ExpectCommit()

		err := u.Create(offer)
		assert.NotNil(t, err)
		assert.EqualValues(t, exp, err)
	})
}

func TestUpdate(t *testing.T) {
	gormDB, mock := setupDB()
	defer gormDB.Close()

	t.Run("Update a user", func(t *testing.T) {
		offer := &offer.Offer{
			Name:               "TEST",
			DiscountPercentage: 24,
		}

		u := NewOfferRepo(gormDB)

		mock.ExpectBegin()

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`INSERT INTO "offers" ("created_at","updated_at","deleted_at","name","discount_percentage") VALUES ($1,$2,$3,$4,$5) RETURNING "offers"."id`)).
			WithArgs(AnyTime{}, AnyTime{}, nil, "TEST", 24).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(1))

		mock.ExpectCommit()

		err := u.Update(offer)
		assert.Nil(t, err)
	})

	t.Run("Update a user fails", func(t *testing.T) {
		exp := errors.New("oops")
		offer := &offer.Offer{
			Name:               "TEST",
			DiscountPercentage: 24,
		}

		u := NewOfferRepo(gormDB)

		mock.ExpectBegin()

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`INSERT INTO "offers" ("created_at","updated_at","deleted_at","name","discount_percentage") VALUES ($1,$2,$3,$4,$5) RETURNING "offers"."id`)).
			WithArgs(AnyTime{}, AnyTime{}, nil, "TEST", 24).
			WillReturnError(exp)

		mock.ExpectCommit()

		err := u.Update(offer)
		assert.NotNil(t, err)
		assert.EqualValues(t, exp, err)
	})
}
