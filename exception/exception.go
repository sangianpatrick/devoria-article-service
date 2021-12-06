package exception

import "fmt"

var (
	ErrConflicted     = fmt.Errorf("conflicted")
	ErrInternalServer = fmt.Errorf("internal server error")
)
