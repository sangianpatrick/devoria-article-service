package account

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const AccountSessionKeyFormat = "account:session:%s"

type AccountContextKey struct{}

// Account is a collection of proprty of account.
type Account struct {
	ID             int64      `json:"id"`
	Email          string     `json:"email"`
	Password       *string    `json:"password,omitempty"`
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	CreatedAt      time.Time  `json:"createdAt"`
	LastModifiedAt *time.Time `json:"lastModifiedAt"`
}

// CustomerStandardJWTClaims is a model.
type AccountStandardJWTClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}
