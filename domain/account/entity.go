package account

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Account is a collection of proprty of account.
type Account struct {
	ID             int64
	Email          string
	Password       *string
	FirstName      string
	LastName       string
	CreatedAt      time.Time
	LastModifiedAt *time.Time
}

// CustomerStandardJWTClaims is a model.
type AccountStandardJWTClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}
