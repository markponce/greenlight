# greenlight

## initalize module
go mod init github.com/markponce/greenlight

## Generating the skeleton directory structure
mkdir -p bin cmd/api internal migrations remote
touch Makefile
touch cmd/api/main.go

## Curl Command
curl localhost:4000/v1/healthcheck
curl -X POST localhost:4000/v1/movies
curl -i -X POST localhost:4000/v1/healthcheck

## Install migration tool
 brew install golang-migrate

## Create migrations
migrate create -seq -ext=.sql -dir=./migrations create_movies_table

## Run migration 
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up

## Rollback migration
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN down

## Postgres Check table definitions
\d <table name>

## Search running process
pgrep -l api

