.PHONY:db-setup\ 
	build

db-setup:
	docker compose up -d database
	PGPASSWORD=postgres psql -d notify-stock -U postgres -p 5555 -h localhost < schema.sql
	PGPASSWORD=postgres psql -d notify-stock-test -U postgres -p 5555 -h localhost < schema.sql

build:
	/usr/local/go/bin/go build -o main  cmd/main.go
