package session

import (
	"context"
	"fmt"
)

// Errors
var (
	ErrSessionNotFound = fmt.Errorf("session not found")
	ErrUnexpected      = fmt.Errorf("nexpected session error")
)

// Session is collection of behavior of session.
type Session interface {
	Set(ctx context.Context, key string, value []byte) (err error)
	Get(ctx context.Context, key string) (value []byte, err error)
	Update(ctx context.Context, key string, value []byte) (err error)
	Delete(ctx context.Context, key string) (err error)
}
