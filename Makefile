.PHONY: run, build-server, create_tables, drop_tables

run:
	docker network create mynetwork || true
	docker-compose up -d

build-server:
	docker-compose build server

create_tables:
	psql postgresql://postgres:postgres@localhost:13080/postgres -f build/package/postgres/init.sql

drop_tables:
	psql postgresql://postgres:postgres@localhost:13080/postgres -f build/package/postgres/drop_all.sql

gen-swagg:
	swag init --parseDependency --parseInternal --parseDepth 3 -g ./cmd/main.go -o ./docs

run-front:
	cd frontend && docker-compose up -d
