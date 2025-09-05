.PHONY: build clean migrate create-admin serve dev

build:
	go build -o main .

clean:
	rm -f main

migrate:
	go run migrate/migrate.go

create-admin:
	go run scripts/create_admin.go

deps:
	go mod tidy

serve:
	go run .

dev: deps migrate create-admin serve
