package persistent

import (
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/andrei-cloud/go-devops/internal/mocks"
)

func Test_storage_Ping(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockDB := mocks.NewMockPersistentDB(ctl)
	gomock.InOrder(
		mockDB.EXPECT().Ping().Return(nil),
	)

	mockDB.Ping()
}
