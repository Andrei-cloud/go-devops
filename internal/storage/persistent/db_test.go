package persistent

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func Test_storage_Ping(t *testing.T) {
	db, _ := newMock()
	repo := &storage{db}
	defer func() {
		repo.Close()
	}()

	assert.NoError(t, repo.Ping())
}

func Test_storage_UpdateGauge(t *testing.T) {
	db, mock := newMock()
	repo := &storage{db}
	defer func() {
		repo.Close()
	}()

	query := "^insert into metrics (.+)"

	mock.ExpectExec(query).WithArgs("test", 1.234).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(query).WithArgs("fail", 1.234).WillReturnError(fmt.Errorf("DB error"))
	assert.NoError(t, repo.UpdateGauge(context.Background(), "test", 1.234))
	assert.Error(t, repo.UpdateGauge(context.Background(), "fail", 1.234))

}

func Test_storage_UpdateCounter(t *testing.T) {
	db, mock := newMock()
	repo := &storage{db}
	defer func() {
		repo.Close()
	}()

	query := "^insert into metrics (.+)"

	mock.ExpectExec(query).WithArgs("test", 1234).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(query).WithArgs("fail", 1234).WillReturnError(fmt.Errorf("DB error"))
	assert.NoError(t, repo.UpdateCounter(context.Background(), "test", 1234))
	assert.Error(t, repo.UpdateCounter(context.Background(), "fail", 1234))
}

func Test_storage_GetCounter(t *testing.T) {
	db, mock := newMock()
	repo := &storage{db}
	defer func() {
		repo.Close()
	}()

	query := "^SELECT delta FROM metrics WHERE mtype = 'counter' (.+)"

	rows := sqlmock.NewRows([]string{"delta"}).AddRow(123)
	mock.ExpectQuery(query).WithArgs("test").WillReturnRows(rows)
	mock.ExpectQuery(query).WithArgs("fail").WillReturnError(fmt.Errorf("DB error"))
	value, err := repo.GetCounter(context.Background(), "test")
	assert.Equal(t, value, int64(123))
	assert.NoError(t, err)

	_, err = repo.GetCounter(context.Background(), "fail")
	assert.Error(t, err)

}

func Test_storage_GetGauge(t *testing.T) {
	db, mock := newMock()
	repo := &storage{db}
	defer func() {
		repo.Close()
	}()

	query := "^SELECT value FROM metrics WHERE mtype = 'gauge' (.+)"

	rows := sqlmock.NewRows([]string{"value"}).AddRow(1.234)
	mock.ExpectQuery(query).WithArgs("test").WillReturnRows(rows)
	mock.ExpectQuery(query).WithArgs("fail").WillReturnError(fmt.Errorf("DB error"))
	value, err := repo.GetGauge(context.Background(), "test")
	assert.Equal(t, value, float64(1.234))
	assert.NoError(t, err)

	_, err = repo.GetGauge(context.Background(), "fail")
	assert.Error(t, err)

}

func Test_storage_GetGaugeAll(t *testing.T) {
	db, mock := newMock()
	repo := &storage{db}
	defer func() {
		repo.Close()
	}()

	query := "^SELECT id, value FROM metrics WHERE mtype = 'gauge'"

	rows := sqlmock.NewRows([]string{"id", "value"}).
		AddRow("one", 1.234).
		AddRow("two", 4.321)
	mock.ExpectQuery(query).WillReturnRows(rows)
	mock.ExpectQuery(query).WillReturnError(fmt.Errorf("DB error"))
	values, err := repo.GetGaugeAll(context.Background())
	assert.Equal(t, values["one"], float64(1.234))
	assert.NoError(t, err)

	_, err = repo.GetGaugeAll(context.Background())
	assert.Error(t, err)
}

func Test_storage_GetCounterAll(t *testing.T) {
	db, mock := newMock()
	repo := &storage{db}
	defer func() {
		repo.Close()
	}()

	query := "^SELECT id, delta FROM metrics WHERE mtype = 'counter'"

	rows := sqlmock.NewRows([]string{"id", "value"}).
		AddRow("one", 1234).
		AddRow("two", 4321)
	mock.ExpectQuery(query).WillReturnRows(rows)
	mock.ExpectQuery(query).WillReturnError(fmt.Errorf("DB error"))
	values, err := repo.GetCounterAll(context.Background())
	assert.Equal(t, values["one"], int64(1234))
	assert.NoError(t, err)

	_, err = repo.GetCounterAll(context.Background())
	assert.Error(t, err)
}
