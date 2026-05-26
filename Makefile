.PHONY: build docker compose-up compose-down compose-reset test test-integration run-server-postgres

build:
	mkdir -p bin
	go build -o bin/prisoner .

docker:
	docker build -t prisoner .

compose-up:
	docker compose up -d db

compose-down:
	docker compose down

compose-reset:
	docker compose down -v

test:
	go test ./...

test-integration:
	go test -tags=integration -count=1 -timeout=10m ./internal/store/...

run-server-postgres: build compose-up
	@until docker compose exec -T db pg_isready -U prisoner -d prisoner >/dev/null 2>&1; do sleep 0.3; done
	DATABASE_URL='postgres://prisoner:prisoner@127.0.0.1:5432/prisoner?sslmode=disable' ./bin/prisoner --store postgres server
