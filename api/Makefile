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

notify:
	go run cmd/main.go notify -s "^N225,^GSPC"
update:
	go run cmd/main.go stock update
update-all:
	go run cmd/main.go stock update -a
gqlgen:
	go run github.com/99designs/gqlgen generate
db-connect:
	psql -d notify-stock -h localhost -U postgres -p 5555
