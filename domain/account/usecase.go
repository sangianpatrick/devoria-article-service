package account

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sangianpatrick/devoria-article-service/crypto"
	"github.com/sangianpatrick/devoria-article-service/exception"
	"github.com/sangianpatrick/devoria-article-service/jwt"
	"github.com/sangianpatrick/devoria-article-service/response"
	"github.com/sangianpatrick/devoria-article-service/session"
)

type AccountUsecase interface {
	Register(ctx context.Context, params AccountRegistrationRequest) (resp response.Response)
	Login(ctx context.Context, params AccountAuthenticationRequest) (resp response.Response)
	GetProfile(ctx context.Context) (resp response.Response)
}

type accountUsecaseImpl struct {
	globalIV     string
	session      session.Session
	jsonWebToken jwt.JSONWebToken
	crypto       crypto.Crypto
	location     *time.Location
	repository   AccountRepository
}

func NewAccountUsecase(
	globalIV string,
	session session.Session,
	jsonWebToken jwt.JSONWebToken,
	crypto crypto.Crypto,
	location *time.Location,
	repository AccountRepository,
) AccountUsecase {
	return &accountUsecaseImpl{
		globalIV:     globalIV,
		session:      session,
		jsonWebToken: jsonWebToken,
		crypto:       crypto,
		location:     location,
		repository:   repository,
	}
}

func (u *accountUsecaseImpl) generateBase64String(byteSize int) string {
	b := make([]byte, byteSize)

	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	// encoded := base64.StdEncoding.EncodeToString(b)
	encoded := hex.EncodeToString(b)

	return encoded
}

func (u *accountUsecaseImpl) Register(ctx context.Context, params AccountRegistrationRequest) (resp response.Response) {
	fmt.Println("len", len(u.generateBase64String(8)))
	fmt.Println(u.globalIV)
	_, err := u.repository.FindByEmail(ctx, params.Email)
	if err == nil {
		return response.Error(response.StatusConflicted, nil, exception.ErrConflicted)
	}

	if err != exception.ErrNotFound {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}
	encryptedPassword := u.crypto.Encrypt(params.Password, u.globalIV)
	newAccount := Account{}
	newAccount.Email = params.Email
	newAccount.Password = &encryptedPassword
	newAccount.FirstName = params.FirstName
	newAccount.LastName = params.LastName
	newAccount.CreatedAt = time.Now().In(u.location)

	ID, err := u.repository.Save(ctx, newAccount)
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}
	newAccount.ID = ID

	claims := AccountStandardJWTClaims{}
	claims.Email = newAccount.Email
	claims.Subject = fmt.Sprintf("%d", newAccount.ID)
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Hour * 24 * 1).Unix()

	token, err := u.jsonWebToken.Sign(ctx, claims)
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	newAccountBuff, _ := json.Marshal(newAccount)

	err = u.session.Set(ctx, fmt.Sprintf(AccountSessionKeyFormat, newAccount.Email), newAccountBuff)
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	// publish to kafke if availabe

	newAccount.Password = nil

	accountAuthenticationResponse := AccountAuthenticationResponse{}
	accountAuthenticationResponse.Token = token
	accountAuthenticationResponse.Profile = newAccount

	return response.Success(response.StatusCreated, accountAuthenticationResponse)
}

func (u *accountUsecaseImpl) Login(ctx context.Context, params AccountAuthenticationRequest) (resp response.Response) {
	account, err := u.repository.FindByEmail(ctx, params.Email)
	if err != nil {
		if err == exception.ErrNotFound {
			return response.Error(response.StatusInvalidPayload, nil, exception.ErrBadRequest)
		}
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	encryptedPassword := u.crypto.Encrypt(params.Password, u.globalIV)
	if encryptedPassword != *account.Password {
		return response.Error(response.StatusInvalidPayload, nil, exception.ErrBadRequest)
	}

	claims := AccountStandardJWTClaims{}
	claims.Email = account.Email
	claims.Subject = fmt.Sprintf("%d", account.ID)
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Hour * 24 * 1).Unix()

	token, err := u.jsonWebToken.Sign(ctx, claims)
	if err == nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	accountBuff, _ := json.Marshal(account)

	err = u.session.Set(ctx, fmt.Sprintf(AccountSessionKeyFormat, account.Email), accountBuff)
	if err != nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	// publish to kafke if availabe

	account.Password = nil

	accountAuthenticationResponse := AccountAuthenticationResponse{}
	accountAuthenticationResponse.Token = token
	accountAuthenticationResponse.Profile = account

	return response.Success(response.StatusOK, accountAuthenticationResponse)
}
func (u *accountUsecaseImpl) GetProfile(ctx context.Context) (resp response.Response) {
	account, ok := ctx.Value(AccountContextKey{}).(Account)
	if !ok {
		return response.Error(response.StatusUnauthorized, nil, exception.ErrUnauthorized)
	}

	return response.Success(response.StatusOK, account)
}
