package account_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sangianpatrick/devoria-article-service/domain/account"
	"github.com/sangianpatrick/devoria-article-service/domain/account/mocks"
	"github.com/sangianpatrick/devoria-article-service/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type responseBody struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}

func TestHandlerRegister_Success(t *testing.T) {
	newAccountRegistrationRequest := account.AccountRegistrationRequest{
		Email:     "john.doe@email.com",
		Password:  "P@ssw0rdTest",
		FirstName: "John",
		LastName:  "Doe",
	}
	accountAuthenticationResponse := account.AccountAuthenticationResponse{
		Token: "token",
		Profile: account.Account{
			ID:        1,
			Email:     newAccountRegistrationRequest.Email,
			Password:  &newAccountRegistrationRequest.Password,
			FirstName: newAccountRegistrationRequest.FirstName,
			LastName:  newAccountRegistrationRequest.LastName,
			CreatedAt: time.Now(),
		},
	}

	resp := response.Success(response.StatusCreated, accountAuthenticationResponse)

	validate := validator.New()

	accountUsecase := new(mocks.AccountUsecase)
	accountUsecase.On("Register", mock.Anything, mock.AnythingOfType("account.AccountRegistrationRequest")).Return(resp)

	newAccountRegistrationRequestBuff, _ := json.Marshal(newAccountRegistrationRequest)

	accountHTTPHandler := account.AccountHTTPHandler{
		Validate: validate,
		Usecase:  accountUsecase,
	}

	r := httptest.NewRequest(http.MethodPost, "/just/for/testing", bytes.NewReader(newAccountRegistrationRequestBuff))
	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(accountHTTPHandler.Register)
	handler.ServeHTTP(recorder, r)

	rb := responseBody{}
	if err := json.NewDecoder(recorder.Body).Decode(&rb); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, response.StatusCreated, rb.Status, fmt.Sprintf("should be status '%s'", response.StatusCreated))
	assert.NotNil(t, rb.Data, "should not be nil")

	data, ok := rb.Data.(map[string]interface{})
	if !ok {
		t.Error("response data isn't a type of 'map[string]interface{}'")
		return
	}

	assert.Equal(t, accountAuthenticationResponse.Token, data["token"], fmt.Sprintf("token should be '%s'", accountAuthenticationResponse.Token))

	accountUsecase.AssertExpectations(t)
}
