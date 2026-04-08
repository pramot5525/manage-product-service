package repository

import (
	"errors"
	"regexp"
	"testing"

	"prm-product/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	cleanup := func() {
		_ = sqlDB.Close()
	}

	return db, mock, cleanup
}

func TestProductRepositoryCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		repo := NewProductRepository(db)
		desc := "Phone"
		salePrice := 99.9
		product := &model.Product{
			Name:        "iPhone",
			Description: &desc,
			SalePrice:   &salePrice,
			Price:       120.0,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "products" ("name","description","sale_price","price") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
			WithArgs(product.Name, product.Description, product.SalePrice, product.Price).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(101)))
		mock.ExpectCommit()

		id, err := repo.Create(product)

		require.NoError(t, err)
		require.Equal(t, int64(101), id)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		repo := NewProductRepository(db)
		product := &model.Product{Name: "iPhone", Price: 120.0}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "products" ("name","description","sale_price","price") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
			WithArgs(product.Name, product.Description, product.SalePrice, product.Price).
			WillReturnError(errors.New("insert failed"))
		mock.ExpectRollback()

		id, err := repo.Create(product)

		require.Error(t, err)
		require.Zero(t, id)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepositoryUpdate(t *testing.T) {
	t.Run("empty update map", func(t *testing.T) {
		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		repo := NewProductRepository(db)

		err := repo.Update(1, map[string]interface{}{})

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		repo := NewProductRepository(db)
		updateData := map[string]interface{}{"name": "iPhone Pro"}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "name"=$1 WHERE id = $2`)).
			WithArgs("iPhone Pro", int64(10)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Update(10, updateData)

		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, cleanup := newMockDB(t)
		defer cleanup()

		repo := NewProductRepository(db)
		updateData := map[string]interface{}{"price": 1500.0}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "price"=$1 WHERE id = $2`)).
			WithArgs(1500.0, int64(20)).
			WillReturnError(errors.New("update failed"))
		mock.ExpectRollback()

		err := repo.Update(20, updateData)

		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
