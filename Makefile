migration_name?=
dsn?=
postgres_user?=
postgres_password?=
db_name?=

serve:
	go run cmd/main.go

test:
	go test -cover ./... 

postgres:
	docker run --name postgresdb -p 5432:5432 -e POSTGRES_USER=${postgres_user} -e POSTGRES_PASSWORD=${postgres_password} -d postgres

createdb:
	docker exec -it postgresdb createdb --username=root --owner=root ${db_name}

dropdb:
	docker exec -it postgresdb dropdb ${db_name}

# dsn?= the database to run the migration. This is option is to be set as a flag for security reasons.
create-migration: 
	migrate create -ext sql -dir db/migration -seq $(migration_name)

migrateup:
	migrate -path db/migration -database $(dsn) -verbose up

# (Optional) version?=.... to rollback to a previous version in a case where a migration fails.
migratedown:
	if [ $(version) ]; then \
		migrate -path db/migration -database $(dsn) -verbose force $(version); \
	else \
		migrate -path db/migration -database $(dsn) -verbose down; \
	fi 

mockery:
	mockgen --build_flags=--mod=mod --package mock --destination mock/sqlstore.go github.com/adetunjii/ohlc/store Sqlstore 


PHONY: serve test create_migrate migrateup migratedown