package account_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sangianpatrick/devoria-article-service/domain/account"
	"github.com/stretchr/testify/assert"
)

var (
	tableName string = "account"
)

func TestSave_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()

	ctx := context.TODO()

	newAccountPassowrd := "P@sswordTest"
	newAccount := account.Account{
		Email:     "john.doe@email.com",
		Password:  &newAccountPassowrd,
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
	}

	expectedCommand := fmt.Sprintf("INSERT INTO %s", tableName)
	expectedArgs := make([]driver.Value, 0)
	expectedArgs = append(
		expectedArgs,
		newAccount.Email,
		*newAccount.Password,
		newAccount.FirstName,
		newAccount.LastName,
		newAccount.CreatedAt,
	)

	mock.ExpectPrepare(expectedCommand).
		ExpectExec().
		WithArgs(expectedArgs...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	accountRepository := account.NewAccountRepository(db, tableName)
	ID, err := accountRepository.Save(ctx, newAccount)

	assert.NoError(t, err, "should not be error")
	assert.Equal(t, int64(1), ID, "id should be `1`")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
