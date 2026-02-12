# Include variables from the .envrc file
include .envrc

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	 @sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api \
		-db-dsn=${NAILIT_DB_DSN} \
		-limiter-enabled=${LIMITER_ENABLED} \
		-cors-trusted-origins="${CORS_TRUSTED_ORIGINS}"



## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${NAILIT_DB_DSN}


## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}


## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${NAILIT_DB_DSN} up