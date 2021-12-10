# devoria-article-service
Workshop Devoria: Clean Architecture Implementation using Golang

# How to run test with cover profile
```bash
make cover
```

# Hot to run this application
### setup environment
```bash
APP_NAME=devoria_article_service
APP_PORT=9001
MARIADB_HOST=localhost
MARIADB_PORT=3306
MARIADB_USERNAME=root
MARIADB_PASSWORD=password
MARIADB_DATABASE=devoria_article_service
MARIADB_MAX_OPEN_CONNECTIONS=50
MARIADB_MAX_IDLE_CONNECTIONS=50
REDIS_HOST=localhost:6379
REDIS_PASSWORD=
REDIS_DATABASE=0
BASIC_AUTH_USERNAME=devoria
BASIC_AUTH_PASSWORD=challenge
AES_SECRET_KEY=279988E50A8194FCED59646B2DB90710
GLOBAL_IV=1234567890123456
```
### for development
```bash
make run.dev
```
### for production
```bash
make build
./app
```
