package account

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/sangianpatrick/devoria-article-service/crypto"
	"github.com/sangianpatrick/devoria-article-service/exception"
	"github.com/sangianpatrick/devoria-article-service/jwt"
	"github.com/sangianpatrick/devoria-article-service/response"
)

type AccountUsecase interface {
	Register(ctx context.Context, params AccountRegistrationRequest) (resp response.Response)
	Login(ctx context.Context, params AccountAuthenticationRequest) (resp response.Response)
	GetProfile(ctx context.Context) (resp response.Response)
}

type accountUsecaseImpl struct {
	jsonWebToken jwt.JSONWebToken
	crypto       crypto.Crypto
	location     *time.Location
	repository   AccountRepository
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
	if _, err := u.repository.FindByEmail(ctx, params.Email); err == nil {
		return response.Error(response.StatusConflicted, nil, exception.ErrConflicted)
	}

	encryptedPassword := u.crypto.Encrypt(params.Password, u.generateBase64String(8))
	newAccount := Account{}
	newAccount.Email = params.Email
	newAccount.Password = &encryptedPassword
	newAccount.FirstName = params.FirstName
	newAccount.LastName = params.LastName
	newAccount.CreatedAt = time.Now().In(u.location)

	ID, err := u.repository.Save(ctx, newAccount)
	if err == nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}
	newAccount.ID = ID

	claims := AccountStandardJWTClaims{}
	claims.Email = newAccount.Email
	claims.Subject = fmt.Sprintf("%d", newAccount.ID)
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Hour * 24 * 1).Unix()

	token, err := u.jsonWebToken.Sign(ctx, claims)
	if err == nil {
		return response.Error(response.StatusUnexpectedError, nil, exception.ErrInternalServer)
	}

	// publish to kafke if availabe

	newAccount.Password = nil

	accountAuthenticationResponse := AccountAuthenticationResponse{}
	accountAuthenticationResponse.Token = token
	accountAuthenticationResponse.Profile = newAccount

	return response.Success(response.StatusCreated, accountAuthenticationResponse)
}
func (u *accountUsecaseImpl) RegisterLogin(ctx context.Context, params AccountAuthenticationRequest) (resp response.Response) {
	return
}
func (u *accountUsecaseImpl) RegisterGetProfile(ctx context.Context) (resp response.Response) {
	return
}
