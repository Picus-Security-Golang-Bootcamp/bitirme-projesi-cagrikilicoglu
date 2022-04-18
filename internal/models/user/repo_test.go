package user

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

	repository *UserRepository
	user       *models.User
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
	s.repository = NewUserRepository(s.DB)
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

var (
	id       = uuid.New()
	email    = "test@testmail.com"
	password = "testPassword"
)

var user = models.User{
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	DeletedAt: gorm.DeletedAt{},
	ID:        id,
	Email:     &email,
	Password:  &password,
	FirstName: "John",
	LastName:  "Doe",
	ZipCode:   "10001",
	Role:      "user",
}

func (s *Suite) TestUserRepository_get() {
	var (
		query_1 = `SELECT * FROM "users" WHERE "users"."email" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`

		row_1 = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "email", "password", "first_name", "last_name", "zip_code", "role"}).
			AddRow(user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.ID.String(), user.Email, user.Password, user.FirstName, user.LastName, user.ZipCode, user.Role)
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_1)).
		WithArgs(email).
		WillReturnRows(row_1)

	res, err := s.repository.get(email)

	require.NoError(s.T(), err)
	require.True(s.T(), reflect.DeepEqual(&user, res))
}

func (s *Suite) TestUserRepository_get_Error() {
	var (
		query_1 = `SELECT * FROM "users" WHERE "users"."email" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`

		row_1 = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "email", "password", "first_name", "last_name", "zip_code", "role"})
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		query_1)).
		WithArgs(email).
		WillReturnRows(row_1)

	res, err := s.repository.get(email)

	require.Error(s.T(), err)
	require.Empty(s.T(), res)
}

// func (s *Suite) TestUserRepository_Create() {
// 	var (
// 		query_1 = `INSERT INTO "user" ("created_at","updated_at","deleted_at","id","email","password","first_name","last_name","zip_code","role")
//        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "user"."id"`

// 		row_1 = sqlmock.NewRows([]string{"created_at", "updated_at", "deleted_at", "id", "email", "password", "first_name", "last_name", "zip_code", "role"}).
// 			AddRow(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.ID.String(), user.Email, user.Password, user.FirstName, user.LastName, user.ZipCode, user.Role)
// 	)
// 	// s.mock.ExpectBegin()
// 	// prep := s.mock.ExpectPrepare(query_1)
// 	// prep.ExpectExec().WithArgs(user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.ID.String(), user.Email, user.Password, user.FirstName, user.LastName, user.ZipCode, user.Role).WillReturnResult(sqlmock.NewResult(0, 1))
// 	// s.mock.ExpectCommit()
// 	s.mock.ExpectBegin()
// 	s.mock.ExpectQuery(regexp.QuoteMeta(
// 		query_1)).
// 		WithArgs(user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.ID.String(), user.Email, user.Password, user.FirstName, user.LastName, user.ZipCode, user.Role).
// 		WillReturnRows(
// 			row_1)
// 	s.mock.ExpectCommit()
// 	res, err := s.repository.Create(&user)

// 	require.NoError(s.T(), err)
// 	require.True(s.T(), reflect.DeepEqual(&user, res))
// }
