.PHONY: run, build-server, create_tables, drop_tables

run:
	docker network create mynetwork || true
	docker-compose up

build-server:
	docker-compose build server

create_tables:
	psql postgresql://postgres:postgres@localhost:13080/postgres -f build/package/postgres/init.sql

drop_tables:
	psql postgresql://postgres:postgres@localhost:13080/postgres -f build/package/postgres/drop_all.sql

# generate_data:
# 	cd gen && python3 gen.py

# fill_tables: create_tables
# 	psql postgresql://postgres:postgres@localhost:13080/postgres -f gen/load_data.sql

# curl -X POST http://localhost:8080/users/signup -H 'Content-Type: application/json' -d '{"login": "larchik", "password": "zhopa332"}'

# curl -X POST http://localhost:8080/users/signin -H 'Content-Type: application/json' -d '{"login": "larchik", "password": "newpasssas"}'

# curl -X GET http://localhost:8080/users/me -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoyLCJsb2dpbiI6ImxhcmNoaWsifSwiZXhwIjoxNzA3NzYyNTI3LCJpYXQiOjE3MDc1MDMzMjd9.FDN-nbv_OWxViX7y1tgzsogDWv4gBxhT7vGrLtuk9aI' -d '{"login": "larchik", "password": "zhopa332"}'

# curl -X POST http://localhost:8080/users/changepass -H 'Content-Type: application/json' -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoyLCJsb2dpbiI6ImxhcmNoaWsifSwiZXhwIjoxNzA3NzYyNTI3LCJpYXQiOjE3MDc1MDMzMjd9.FDN-nbv_OWxViX7y1tgzsogDWv4gBxhT7vGrLtuk9aI' -d '{"oldPassword": "zhopa332", "newPassword": "newpasssas"}'