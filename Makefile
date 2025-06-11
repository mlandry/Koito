.PHONY: all test clean client

db.up:
	GOOSE_MIGRATION_DIR=db/migrations GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://postgres:secret@localhost:5432 goose up

db.down:
	GOOSE_MIGRATION_DIR=db/migrations GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://postgres:secret@localhost:5432 goose down

db.reset:
	GOOSE_MIGRATION_DIR=db/migrations GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://postgres:secret@localhost:5432 goose down-to 0

db.schemadump:
	docker run --rm --network=host --env PGPASSWORD=secret -v "./db:/tmp/dump" \
	postgres pg_dump \
	--schema-only \
	--host=192.168.0.153 \
	--port=5432 \
	--username=postgres \
	-v --dbname="koitodb" -f "/tmp/dump/schema.sql"

postgres.run:
	docker run --name koito-db -p 5432:5432 -e POSTGRES_PASSWORD=secret -d postgres

postgres.start:
	docker start koito-db

postgres.stop:
	docker stop koito-db

postgres.rm:
	docker rm bamsort-db

api.debug:
	KOITO_ALLOWED_HOSTS=* KOITO_LOG_LEVEL=debug KOITO_CONFIG_DIR=test_config_dir KOITO_DATABASE_URL=postgres://postgres:secret@192.168.0.153:5432/koitodb?sslmode=disable go run cmd/api/main.go

api.test:
	go test ./... -timeout 60s

client.dev:
	cd client && yarn run dev

docs.dev:
	cd docs && yarn dev

client.build:
	cd client && yarn run build

test: api.test