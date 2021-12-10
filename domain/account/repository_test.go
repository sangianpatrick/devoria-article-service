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

func TestFindByEmail_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()

	ctx := context.TODO()
	expectedAccountPassowrd := "P@sswordTest"
	expectedAccountEmail := "john.doe@email.com"
	expectedRows := sqlmock.NewRowsWithColumnDefinition(
		sqlmock.NewColumn("id"),
		sqlmock.NewColumn("email"),
		sqlmock.NewColumn("password"),
		sqlmock.NewColumn("firstName"),
		sqlmock.NewColumn("lastName"),
		sqlmock.NewColumn("createdAt"),
		sqlmock.NewColumn("lastModified"),
	).AddRow(
		int64(1),
		expectedAccountEmail,
		&expectedAccountPassowrd,
		"John",
		"Doe",
		time.Now(),
		nil,
	)

	expectedQuery := fmt.Sprintf(`SELECT id, email, password, firstName, lastName, createdAt, lastModified FROM %s WHERE email = ?`, tableName)
	expectedArgs := make([]driver.Value, 0)
	expectedArgs = append(expectedArgs, expectedAccountEmail)

	mock.ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs(expectedArgs...).
		WillReturnRows(expectedRows)

	accountRepository := account.NewAccountRepository(db, tableName)
	existingAccount, err := accountRepository.FindByEmail(ctx, expectedAccountEmail)

	assert.NoError(t, err, "should not be error")
	assert.Equal(t, expectedAccountEmail, existingAccount.Email, fmt.Sprintf("email should be '%s'", expectedAccountEmail))
	assert.Nil(t, existingAccount.LastModifiedAt, "last modified at should be nil")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
