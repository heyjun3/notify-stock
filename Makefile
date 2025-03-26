.PHONY:db-setup\ 
	build\
	notify

db-setup:
	docker compose up -d database
	PGPASSWORD=postgres psql -d notify-stock -U postgres -p 5555 -h localhost < schema.sql
	PGPASSWORD=postgres psql -d notify-stock-test -U postgres -p 5555 -h localhost < schema.sql

build:
	/usr/local/go/bin/go build -o main  cmd/main.go

notify:
	go run cmd/main.go notify -s "N225,S&P500"
