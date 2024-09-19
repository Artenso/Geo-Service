.PHONY: generate
generate:
	mkdir -p pkg/user_service
	protoc --proto_path vendor.protogen \
	--proto_path api \
	--go_out=pkg/user_service \
	--go-grpc_out=pkg/user_service \
	api/user_service.proto
		make move

.PHONY: move
move:
	mv pkg/user_service/github.com/Artenso/Geo-Service/user/pkg/user_service/* pkg/user_service &&\
	rm -rf pkg/user_service/github.com/ 

.PHONY: install-goose
install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

include .env
LOCAL_MIGRATION_DIR=migrations
LOCAL_MIGRATION_DSN='host=localhost port=${DB_PORT} dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} sslmode=disable'
M_NAME=new_column

.PHONY: local-migration-create
local-migration-create:
	goose -dir ${LOCAL_MIGRATION_DIR} -v -s postgres ${LOCAL_MIGRATION_DSN} create ${M_NAME} sql

.PHONY: local-migration-up
local-migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up

.PHONY: local-migration-down
local-migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down

# usage make local-migration-up-to V=5
.PHONY: local-migration-up-to
local-migration-up-to:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up-to ${V}

.PHONY: local-migration-status
local-migration-status:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status

