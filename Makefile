ifneq ("$(wildcard .env)","")
	include .env
	export $(shell sed 's/=.*//' .env)
endif

.PHONY: test
test:
	go test -cover ./...

.PHONY: test-verbose
test-verbose:
	go test -v -cover ./...

.PHONY: start
start:
	./bin/main

.PHONY: dev
dev:
	air

.PHONY: build
build:
	go build -o=./bin/main ./cmd

.PHONY: models
models:
	pg_dump --schema-only identity > schema.sql
	sqlc generate -f sqlc.yaml

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-up
migrate-up:
	migrate -path=./migrations -database="$(POSTGRES_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path=./migrations -database="$(POSTGRES_URL)" down 1

.PHONY: minio
minio:
	minio server  --console-address :9001 /Users/michaelcorrigan/minio
