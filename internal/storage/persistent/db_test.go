package persistent

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/andrei-cloud/go-devops/internal/repo"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/suite"

	"github.com/DATA-DOG/go-sqlmock"
)

type DBTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo repo.Repository
}

func (s *DBTestSuite) SetupTest() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	s.repo = &storage{s.db}
	s.NoError(err)
}

func (s *DBTestSuite) TearDownTest() {
	s.repo.Close()
}

func (s *DBTestSuite) TestPing() {
	s.NoError(s.repo.Ping())
}

func (s *DBTestSuite) TestUpdateGauge() {
	query := "^insert into metrics (.+)"

	s.mock.ExpectExec(query).WithArgs("test", 1.234).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec(query).WithArgs("fail", 1.234).WillReturnError(fmt.Errorf("DB error"))
	s.NoError(s.repo.UpdateGauge(context.Background(), "test", 1.234))
	s.Error(s.repo.UpdateGauge(context.Background(), "fail", 1.234))
}

func (s *DBTestSuite) TestUpdateCounter() {
	query := "^insert into metrics (.+)"

	s.mock.ExpectExec(query).WithArgs("test", 1234).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectExec(query).WithArgs("fail", 1234).WillReturnError(fmt.Errorf("DB error"))
	s.NoError(s.repo.UpdateCounter(context.Background(), "test", 1234))
	s.Error(s.repo.UpdateCounter(context.Background(), "fail", 1234))
}

func (s *DBTestSuite) TestGetCounter() {
	query := "^SELECT delta FROM metrics WHERE mtype = 'counter' (.+)"

	rows := sqlmock.NewRows([]string{"delta"}).AddRow(123)
	s.mock.ExpectQuery(query).WithArgs("test").WillReturnRows(rows)
	s.mock.ExpectQuery(query).WithArgs("fail").WillReturnError(fmt.Errorf("DB error"))
	value, err := s.repo.GetCounter(context.Background(), "test")
	s.Equal(value, int64(123))
	s.NoError(err)

	_, err = s.repo.GetCounter(context.Background(), "fail")
	s.Error(err)
}

func (s *DBTestSuite) TestGetGauge() {
	query := "^SELECT value FROM metrics WHERE mtype = 'gauge' (.+)"

	rows := sqlmock.NewRows([]string{"value"}).AddRow(1.234)
	s.mock.ExpectQuery(query).WithArgs("test").WillReturnRows(rows)
	s.mock.ExpectQuery(query).WithArgs("fail").WillReturnError(fmt.Errorf("DB error"))
	value, err := s.repo.GetGauge(context.Background(), "test")
	s.Equal(value, float64(1.234))
	s.NoError(err)

	_, err = s.repo.GetGauge(context.Background(), "fail")
	s.Error(err)

}

func (s *DBTestSuite) TestGetGaugeAll() {
	query := "^SELECT id, value FROM metrics WHERE mtype = 'gauge'"

	rows := sqlmock.NewRows([]string{"id", "value"}).
		AddRow("one", 1.234).
		AddRow("two", 4.321)
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.mock.ExpectQuery(query).WillReturnError(fmt.Errorf("DB error"))
	values, err := s.repo.GetGaugeAll(context.Background())
	s.Equal(values["one"], float64(1.234))
	s.NoError(err)

	_, err = s.repo.GetGaugeAll(context.Background())
	s.Error(err)
}

func (s *DBTestSuite) TestGetCounterAll() {
	query := "^SELECT id, delta FROM metrics WHERE mtype = 'counter'"

	rows := sqlmock.NewRows([]string{"id", "value"}).
		AddRow("one", 1234).
		AddRow("two", 4321)
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.mock.ExpectQuery(query).WillReturnError(fmt.Errorf("DB error"))
	values, err := s.repo.GetCounterAll(context.Background())
	s.Equal(values["one"], int64(1234))
	s.NoError(err)

	_, err = s.repo.GetCounterAll(context.Background())
	s.Error(err)
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
