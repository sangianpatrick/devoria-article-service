package account

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type AccountRepository interface {
	Save(ctx context.Context, account Account) (ID int64, err error)
	Update(ctx context.Context, ID int64, updatedAccount Account) (err error)
	FindByEmail(ctx context.Context, email string) (account Account, err error)
	FindByID(ctx context.Context, ID int64) (account Account, err error)
}

type accountRepositoryImpl struct {
	db        *sql.DB
	tableName string
}

func NewAccountRepository(db *sql.DB, tableName string) AccountRepository {
	return &accountRepositoryImpl{
		db:        db,
		tableName: tableName,
	}
}

func (r *accountRepositoryImpl) Save(ctx context.Context, account Account) (ID int64, err error) {
	command := fmt.Sprintf("INSERT INTO %s (email, password, firstName, lastName, createdAt) VALUES (?, ?, ?, ?, ?)", r.tableName)
	stmt, err := r.db.PrepareContext(ctx, command)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		account.Email,
		*account.Password,
		account.FirstName,
		account.LastName,
		account.CreatedAt,
	)

	if err != nil {
		log.Println(err)
		return
	}

	ID, _ = result.LastInsertId()

	return
}

func (r *accountRepositoryImpl) Update(ctx context.Context, ID int64, updatedAccount Account) (err error) {
	command := fmt.Sprintf(`UPDATE %s SET password = ?, firstName = ?, lastName = ?, lastModifiedAt = ? WHERE id = ?`, r.tableName)
	stmt, err := r.db.PrepareContext(ctx, command)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		*updatedAccount.Password,
		updatedAccount.FirstName,
		updatedAccount.LastName,
		updatedAccount.LastModifiedAt,
	)

	if err != nil {
		log.Println(err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		return fmt.Errorf("not found")
	}

	return
}

func (r *accountRepositoryImpl) FindByEmail(ctx context.Context, email string) (account Account, err error) {
	query := fmt.Sprintf(`SELECT id, email, password, firstName, lastName, createdAt, lastModifiedAt FROM %s WHERE email = ?`, r.tableName)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, email)

	var password sql.NullString
	var lastModifiedAt sql.NullTime

	err = row.Scan(
		&account.ID,
		&account.Email,
		&password,
		&account.FirstName,
		&account.LastName,
		&account.CreatedAt,
		&lastModifiedAt,
	)

	if err != nil {
		log.Println(err)
		return
	}

	if password.Valid {
		account.Password = &password.String
	}

	if lastModifiedAt.Valid {
		account.LastModifiedAt = &lastModifiedAt.Time
	}

	return
}

func (r *accountRepositoryImpl) FindByID(ctx context.Context, ID int64) (account Account, err error) {
	query := fmt.Sprintf(`SELECT id, email, password, firstName, lastName, createdAt, lastModifiedAt FROM %s WHERE id = ?`, r.tableName)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, ID)

	var password sql.NullString
	var lastModifiedAt sql.NullTime

	err = row.Scan(
		&account.ID,
		&account.Email,
		&password,
		&account.FirstName,
		&account.LastName,
		&account.CreatedAt,
		&lastModifiedAt,
	)

	if err != nil {
		log.Println(err)
		return
	}

	if password.Valid {
		account.Password = &password.String
	}

	if lastModifiedAt.Valid {
		account.LastModifiedAt = &lastModifiedAt.Time
	}

	return
}
