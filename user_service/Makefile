CURRENT_DIR=$(shell pwd)
DB_URL := "postgres://postgres:nodirbek@localhost:5432/postgres?sslmode=disable"

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto:
	chmod +x ./scripts/gen-proto.sh
	./scripts/gen-proto.sh

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-file:
	migrate create -ext sql -dir migrations/ -seq create_users_table