CURRENT_DIR=$(shell pwd)
DB_URL := "postgres://postgres:nodirbek@localhost:5432/commentdb?sslmode=disable"

proto:
	chmod +x ./scripts/gen-proto.sh
	./scripts/gen-proto.sh
migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-force:
	migrate -path migrations -database "$(DB_URL)" -verbose force 1

migrate-file:
	migrate create -ext sql -dir migrations/ -seq create_comments_table

