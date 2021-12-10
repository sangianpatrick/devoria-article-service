package account_test

import (
	"context"
	"testing"
	"time"

	cryptoMocks "github.com/sangianpatrick/devoria-article-service/crypto/mocks"
	"github.com/sangianpatrick/devoria-article-service/domain/account"
	"github.com/sangianpatrick/devoria-article-service/domain/account/mocks"
	"github.com/sangianpatrick/devoria-article-service/exception"
	jsonWebTokenMocks "github.com/sangianpatrick/devoria-article-service/jwt/mocks"
	sessionMocks "github.com/sangianpatrick/devoria-article-service/session/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	location *time.Location
)

func TestMain(m *testing.M) {
	location, _ = time.LoadLocation("Asia/Jakarta")

	m.Run()
}

func TestUsecaseRegister_Success(t *testing.T) {
	sess := new(sessionMocks.Session)
	sess.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)

	jsonWebToken := new(jsonWebTokenMocks.JSONWebToken)
	jsonWebToken.On("Sign", mock.Anything, mock.AnythingOfType("account.AccountStandardJWTClaims")).Return("mock token", nil)

	crypto := new(cryptoMocks.Crypto)
	crypto.On("Encrypt", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("encrypted password")

	accountRepository := new(mocks.AccountRepository)
	accountRepository.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(account.Account{}, exception.ErrNotFound)
	accountRepository.On("Save", mock.Anything, mock.AnythingOfType("account.Account")).Return(int64(1), nil)

	accountUsecase := account.NewAccountUsecase(
		"globalIVTest",
		sess,
		jsonWebToken,
		crypto,
		location,
		accountRepository,
	)

	ctx := context.TODO()
	params := account.AccountRegistrationRequest{
		Email:     "john.doe@email.com",
		Password:  "P@ssw0rdTest",
		FirstName: "John",
		LastName:  "Doe",
	}
	resp := accountUsecase.Register(ctx, params)

	assert.NoError(t, resp.Err())

	sess.AssertExpectations(t)
	jsonWebToken.AssertExpectations(t)
	crypto.AssertExpectations(t)
	accountRepository.AssertExpectations(t)
}
