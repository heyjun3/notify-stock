.PHONY:db-setup\ 
	build\
	notify\
	gqlgen\
	update\
	update-all\
	gen-local-key\

db-setup:
	docker compose up -d database
	PGPASSWORD=postgres psql -d notify-stock -U postgres -p 5555 -h localhost < schema.sql
	PGPASSWORD=postgres psql -d notify-stock-test -U postgres -p 5555 -h localhost < schema.sql

build:
	docker compose build notify

notify:
	docker compose run --rm notify notify -s "^N225,^GSPC"
update:
	docker compose run --rm notify stock update
update-all:
	docker compose run --rm notify stock update -a
gqlgen:
	go run github.com/99designs/gqlgen generate
