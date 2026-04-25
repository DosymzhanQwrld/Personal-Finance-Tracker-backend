DB_URL=postgres://postgres:postgres@localhost:5432/finance_db?sslmode=disable
MIGRATE_PATH=db/migrations

create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

migrate-force:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" force $(v)

migrate-drop:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" drop