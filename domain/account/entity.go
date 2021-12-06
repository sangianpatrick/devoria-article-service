package account

import "time"

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
