.PHONY:db-setup\ 
	build\
	notify\
	gqlgen

db-setup:
	docker compose up -d database
	PGPASSWORD=postgres psql -d notify-stock -U postgres -p 5555 -h localhost < schema.sql
	PGPASSWORD=postgres psql -d notify-stock-test -U postgres -p 5555 -h localhost < schema.sql

build:
	docker compose build notify

notify:
	docker compose run --rm notify notify -s "N225,S&P500"
register:
	docker compose run --rm notify register -s "N225,S&P500"
gqlgen:
	go run github.com/99designs/gqlgen generate
