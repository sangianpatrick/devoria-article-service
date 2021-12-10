package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sangianpatrick/devoria-article-service/config"
	"github.com/sangianpatrick/devoria-article-service/crypto"
	"github.com/sangianpatrick/devoria-article-service/domain/account"
	"github.com/sangianpatrick/devoria-article-service/jwt"
	"github.com/sangianpatrick/devoria-article-service/middleware"
	"github.com/sangianpatrick/devoria-article-service/session"
	"go.elastic.co/apm/module/apmgoredisv8"
	"go.elastic.co/apm/module/apmgorilla"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/mysql"
)

func main() {
	location, _ := time.LoadLocation("Asia/Jakarta")
	cfg := config.New()

	db, err := apmsql.Open("mysql", cfg.Mariadb.DSN)
	db.SetMaxOpenConns(cfg.Mariadb.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.Mariadb.MaxIdleConnections)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	rc := redis.NewClient(cfg.Redis.Options)
	if _, err := rc.Ping(context.Background()).Result(); err != nil {
		log.Fatal(err)
	}
	rc.AddHook(apmgoredisv8.NewHook())

	vld := validator.New()
	encryption := crypto.NewAES256CBC(cfg.AES.SecretKey)
	jsonWebToken := jwt.NewJSONWebToken(jwt.GetRSAPrivateKey("./secret/id_rsa"), jwt.GetRSAPublicKey("./secret/id_rsa.pub"))
	sess := session.NewRedisSessionStoreAdapter(rc, time.Hour*24*1)
	basicAuthMiddleware := middleware.NewBasicAuth(cfg.BasicAuth.Username, cfg.BasicAuth.Password)

	router := mux.NewRouter()
	apmgorilla.Instrument(router)

	accountRepository := account.NewAccountRepository(db, "account")
	accountUsecase := account.NewAccountUsecase(cfg.GlobalIV, sess, jsonWebToken, encryption, location, accountRepository)
	account.NewAccountHTTPHandler(router, basicAuthMiddleware, vld, accountUsecase)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: router,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	<-sigterm

	fmt.Println("shutting down application ...")

	server.Shutdown(context.Background())
	db.Close()
	rc.Close()
}
