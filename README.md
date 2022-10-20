# Migration

## To create migration

- download `go migrate`

- run `migrate create -ext sql -dir db/migration -seq init_schema` This will create empty migration files with up and down migration.

## To run up migration

- run `migrate -path db/migration -database "postgresql://user:password@localhost:5432/bank?sslmode=disable" -verbose up`
 ### OR

- run `make migrate-up`

## Test database must be running to pass integration tests
Testing script should automatically create migrations, just run:
- 1 `docker compse up`
- 2 Run your tests!
