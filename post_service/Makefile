CURRENT_DIR =$(shell pwd)
DB_URL := "postgres://postgres:nodirbek@localhost:5432/exam?sslmode=disable"

proto:
	chmod +x ./scripts/gen-proto.sh
	./scripts/gen-proto.sh

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up 

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down 

migrate-file:
	migrate create -ext sql -dir migrations/ create_posts_table
