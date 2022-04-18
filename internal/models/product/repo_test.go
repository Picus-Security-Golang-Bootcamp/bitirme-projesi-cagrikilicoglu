package product

import (
	"database/sql"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cagrikilicoglu/shopping-basket/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository *ProductRepository
	product    *models.Product
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(s.T(), err)
	s.repository = NewProductRepository(s.DB)
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

var (
	id   = uuid.New()
	name = "test"
)

var stock = models.Stock{SKU: "TESTSKU", Number: 10}
var product = models.Product{
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	DeletedAt: gorm.DeletedAt{},
	ID:        id,
	Name:      &name,
	Price:     float32(12.0),
	Stock:     stock,
}

func (s *Suite) TestProductRepository_GetBySKU() {
	var (
		query_1 = `SELECT * FROM "products" WHERE "products"."sku" = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT 1`

		row_1 = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "name", "price", "sku", "number"}).
			AddRow(product.CreatedAt, product.UpdatedAt, product.DeletedAt, product.ID.String(), product.Name, product.Price, product.Stock.SKU, product.Stock.Number)
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_1)).
		WithArgs(product.Stock.SKU).
		WillReturnRows(row_1)

	res, err := s.repository.GetBySKU(stock.SKU)

	require.NoError(s.T(), err)
	require.True(s.T(), reflect.DeepEqual(&product, res))
}

// func (s *Suite) TestProductRepository_GetByID() {
// 	var (
// 		query_1 = `SELECT * FROM "products" WHERE 	"products.id" = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT 1`

// 		row_1 = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "name", "price", "sku", "number"}).
// 			AddRow(product.CreatedAt, product.UpdatedAt, product.DeletedAt, product.ID.String(), product.Name, product.Price, product.Stock.SKU, product.Stock.Number)
// 	)

// 	s.mock.ExpectQuery(regexp.QuoteMeta(
// 		query_1)).
// 		WithArgs(product.ID).
// 		WillReturnRows(row_1)

// 	res, err := s.repository.getByID(product.ID.String())

// 	require.NoError(s.T(), err)
// 	require.True(s.T(), reflect.DeepEqual(&product, res))
// }

// func (s *Suite) TestProductRepository_GetByName() {
// 	var (
// 		query_1 = `SELECT * FROM "products" WHERE Name ILIKE $1 AND "products"."deleted_at" IS NULL ORDER BY name`

// 		row_1 = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "name", "price", "sku", "number"}).
// 			AddRow(product.CreatedAt, product.UpdatedAt, product.DeletedAt, product.ID.String(), product.Name, product.Price, product.Stock.SKU, product.Stock.Number)
// 		productSrc = "%" + name + "%"
// 	)

// 	s.mock.ExpectQuery(regexp.QuoteMeta(
// 		query_1)).
// 		WithArgs(productSrc).
// 		WillReturnRows(row_1)

// 	res, err := s.repository.getByName(name)

// 	require.NoError(s.T(), err)
// 	require.True(s.T(), reflect.DeepEqual(&product, res))
// }
