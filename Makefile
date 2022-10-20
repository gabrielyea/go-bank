sqlc:
	sqlc generate

migrate-up:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/bank?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://user:password@localhost:5432/bank?sslmode=disable" -verbose down

test-up:
	migrate -path db/migration -database "postgresql://test:password@localhost:5500/test_bank?sslmode=disable" -verbose up

test-down:
	migrate -path db/migration -database "postgresql://test:password@localhost:5500/test_bank?sslmode=disable" -verbose down

test:
	go test -v -cover ./repo