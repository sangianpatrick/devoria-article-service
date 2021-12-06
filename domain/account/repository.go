package account

import "context"

type AccountRepository interface {
	Save(ctx context.Context, account Account) (ID int64, err error)
	Update(ctx context.Context, ID int64, updatedAccount Account) (err error)
	FindByEmail(ctx context.Context, email string) (account Account, err error)
	FindByID(ctx context.Context, ID int64) (account Account, err error)
}
