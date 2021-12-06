package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Config struct {
	App struct {
		Name string
		Port string
	}
	Logger struct {
		Formatter logrus.Formatter
	}
	Mariadb struct {
		DSN                string
		MaxOpenConnections int
		MaxIdleConnections int
	}
	Redis struct {
		Options *redis.Options
	}
	AES struct {
		SecretKey string
	}
	BasicAuth struct {
		Username string
		Password string
	}
}

func New() *Config {
	c := new(Config)
	c.loadApp()
	c.loadMariadb()
	c.loadRedis()
	c.loadAes()
	c.loadBasicAuth()

	return c
}

func (c *Config) loadApp() *Config {
	name := os.Getenv("APP_NAME")
	port := os.Getenv("APP_PORT")

	c.App.Name = name
	c.App.Port = port

	return c
}

func (c *Config) loadMariadb() *Config {
	host := os.Getenv("MARIADB_HOST")
	port := os.Getenv("MARIADB_PORT")
	username := os.Getenv("MARIADB_USERNAME")
	password := os.Getenv("MARIADB_PASSWORD")
	database := os.Getenv("MARIADB_DATABASE")
	maxOpenConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_MAX_OPEN_CONNECTIONS"), 10, 64)
	maxIdleConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_MAX_IDLE_CONNECTIONS"), 10, 64)

	connVal := url.Values{}
	connVal.Add("parseTime", "1")
	connVal.Add("loc", "Asia/Jakarta")

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	dsn := fmt.Sprintf("%s?%s", dbConnectionString, connVal.Encode())

	c.Mariadb.DSN = dsn
	c.Mariadb.MaxOpenConnections = int(maxOpenConnections)
	c.Mariadb.MaxIdleConnections = int(maxIdleConnections)

	return c
}

func (c *Config) loadRedis() *Config {
	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.ParseInt(os.Getenv("REDIS_DATABASE"), 10, 64)

	options := &redis.Options{
		Addr:     host,
		Password: password,
		DB:       int(db),
	}

	c.Redis.Options = options

	return c
}

func (c *Config) loadAes() *Config {
	secretKey := os.Getenv("AES_SECRET_KEY")
	c.AES.SecretKey = secretKey

	return c
}

func (c *Config) loadBasicAuth() *Config {
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	c.BasicAuth.Username = username
	c.BasicAuth.Password = password

	return c
}
