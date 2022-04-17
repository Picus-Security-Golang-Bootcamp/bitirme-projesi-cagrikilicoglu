package cart

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
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository CartRepo
	cart1      *models.Cart
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
	s.repository = NewCartRepository(s.DB)
}

// func (s *Suite) AfterTest(_, _ string) {
// 	require.True(s.T(), s.mock.ExpectationsWereMet())
// }

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

var id = uuid.New()
var updatedPrice = float32(60.0)

var cart = models.Cart{
	CreatedAt:  time.Now(),
	UpdatedAt:  time.Now(),
	DeletedAt:  gorm.DeletedAt{},
	ID:         id,
	UserID:     uuid.New(),
	Items:      []models.Item{item_1},
	TotalPrice: float32(32.0),
}
var item_1 = models.Item{
	CartID:    id,
	IsOrdered: false,
}

func (s *Suite) TestCartRepository_GetByCartID() {
	var (
		query_1 = `SELECT * FROM "carts" WHERE id = $1 AND "carts"."deleted_at" IS NULL ORDER BY "carts"."id" LIMIT 1`
		query_2 = `SELECT * FROM "items" WHERE "items"."cart_id" = $1 AND is_ordered = $2 AND "items"."deleted_at" IS NULL`
		row_1   = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "user_id", "total_price"}).
			AddRow(cart.CreatedAt, cart.UpdatedAt, cart.DeletedAt, cart.ID.String(), cart.UserID.String(), cart.TotalPrice)
		row_2 = sqlmock.NewRows([]string{"cart_id", "is_ordered"}).
			AddRow(item_1.CartID, item_1.IsOrdered)
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_1)).
		WithArgs(cart.ID.String()).
		WillReturnRows(row_1)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_2)).
		WithArgs(item_1.CartID.String(), item_1.IsOrdered).
		WillReturnRows(row_2)

	res, err := s.repository.GetByCartID(cart.ID.String())

	require.NoError(s.T(), err)
	require.True(s.T(), reflect.DeepEqual(&cart, res))
}

func (s *Suite) TestCartRepository_GetByCartID_Error() {
	var (
		query_1 = `SELECT * FROM "carts" WHERE id = $1 AND "carts"."deleted_at" IS NULL ORDER BY "carts"."id" LIMIT 1`
		query_2 = `SELECT * FROM "items" WHERE "items"."cart_id" = $1 AND is_ordered = $2 AND "items"."deleted_at" IS NULL`
		row_1   = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "user_id", "total_price"})
		row_2   = sqlmock.NewRows([]string{"cart_id", "is_ordered"}).
			AddRow(item_1.CartID, item_1.IsOrdered)
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_1)).
		WithArgs(cart.ID.String()).
		WillReturnRows(row_1)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_2)).
		WithArgs(item_1.CartID.String(), item_1.IsOrdered).
		WillReturnRows(row_2)

	res, err := s.repository.GetByCartID(cart.ID.String())

	require.Error(s.T(), err)
	require.Empty(s.T(), res)
}

func (s *Suite) TestCartRepository_GetByUserID() {
	var (
		query_1 = `SELECT * FROM "carts" WHERE user_id = $1 AND "carts"."deleted_at" IS NULL ORDER BY "carts"."id" LIMIT 1`
		query_2 = `SELECT * FROM "items" WHERE "items"."cart_id" = $1 AND is_ordered = $2 AND "items"."deleted_at" IS NULL`
		row_1   = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "user_id", "total_price"}).
			AddRow(cart.CreatedAt, cart.UpdatedAt, cart.DeletedAt, cart.ID.String(), cart.UserID.String(), cart.TotalPrice)
		row_2 = sqlmock.NewRows([]string{"cart_id", "is_ordered"}).
			AddRow(item_1.CartID, item_1.IsOrdered)
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_1)).
		WithArgs(cart.UserID.String()).
		WillReturnRows(row_1)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_2)).
		WithArgs(item_1.CartID.String(), item_1.IsOrdered).
		WillReturnRows(row_2)

	res, err := s.repository.GetByUserID(cart.UserID.String())
	zap.L().Debug("test", zap.Reflect("res", res))

	require.NoError(s.T(), err)
	require.True(s.T(), reflect.DeepEqual(&cart, res))
}
func (s *Suite) TestCartRepository_GetByUserID_Error() {
	var (
		query_1 = `SELECT * FROM "carts" WHERE user_id = $1 AND "carts"."deleted_at" IS NULL ORDER BY "carts"."id" LIMIT 1`
		query_2 = `SELECT * FROM "items" WHERE "items"."cart_id" = $1 AND is_ordered = $2 AND "items"."deleted_at" IS NULL`
		row_1   = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "user_id", "total_price"})
		row_2   = sqlmock.NewRows([]string{"cart_id", "is_ordered"}).
			AddRow(item_1.CartID, item_1.IsOrdered)
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_1)).
		WithArgs(cart.UserID.String()).
		WillReturnRows(row_1)
	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_2)).
		WithArgs(item_1.CartID.String(), item_1.IsOrdered).
		WillReturnRows(row_2)
	res, err := s.repository.GetByUserID(cart.UserID.String())

	require.Error(s.T(), err)
	require.Empty(s.T(), res)
}

// func (s *Suite) TestCartRepository_UpdateTotalPrice() {
// 	var (
// 		query_1 = `UPDATE "carts" SET total_price = $1 WHERE id = $2 AND "carts"."deleted_at" IS NULL`

// 		// row_1 = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "user_id", "total_price"}).
// 		// 	AddRow(cart.CreatedAt, cart.UpdatedAt, cart.DeletedAt, cart.ID.String(), cart.UserID.String(), cart.TotalPrice)
// 	)
// 	prep := s.mock.ExpectPrepare(query_1)
// 	prep.ExpectExec().
// 		WithArgs(updatedPrice, cart.ID).
// 		WillReturnResult(sqlmock.NewResult(0, 1))

// 	err := s.repository.UpdateTotalPrice(&cart, updatedPrice)

// 	require.NoError(s.T(), err)

// }
