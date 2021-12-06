package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sangianpatrick/devoria-article-service/domain/account"
)

var (
	username = "root"
	password = "14qwafzx"
	host     = "localhost"
	port     = "3306"
	database = "devoria_article_service"
)

var location *time.Location

func init() {
	location, _ = time.LoadLocation("Asia/Jakarta")
}

func main() {
	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	connVal := url.Values{}
	connVal.Add("parseTime", "1")
	connVal.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", dbConnectionString, connVal.Encode())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	accountRepository := account.NewAccountRepository(db, "account")
	newAccountPassword := "123456"
	newAccount := account.Account{
		Email:          "johndoe@mail.com",
		Password:       &newAccountPassword,
		FirstName:      "John",
		LastName:       "Doe",
		CreatedAt:      time.Now().In(location),
		LastModifiedAt: nil,
	}

	ID, err := accountRepository.Save(context.TODO(), newAccount)
	if err != nil {
		log.Fatal(err)
	}

	newAccount.ID = ID

	fmt.Println("result", newAccount)
}
