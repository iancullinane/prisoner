.PHONY: build docker compose-up compose-down compose-reset test test-integration run-server-postgres run-server-memory run-server run-web dev

STORE ?= memory
POSTGRES_DATABASE_URL := postgres://prisoner:prisoner@127.0.0.1:5432/prisoner?sslmode=disable

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

TEST_PKGS    := $(filter-out test,$(MAKECMDGOALS))
TEST_TARGETS := $(if $(TEST_PKGS),$(addsuffix /...,$(addprefix ./,$(TEST_PKGS))),./...)

test:
	go test $(if $(TEST_PKGS),-v -count=1) -cover $(TEST_TARGETS)

ifneq (,$(filter test -v,$(MAKECMDGOALS)))
%:
	@:
endif

test-integration:
	go test -v -tags=integration -count=1 -timeout=10m ./internal/store/...

test-integration-memory:
	go test -v ./internal/store/memory/...

test-integration-file:
	go test -v ./internal/store/file/...

test-integration-postgres:
	go test -v -tags=integration -count=1 -timeout=10m ./internal/store/postgres/...

start-db: compose-up
	@until docker compose exec -T db pg_isready -U prisoner -d prisoner >/dev/null 2>&1; do sleep 0.3; done

run-server-postgres: build start-db
	DATABASE_URL='$(POSTGRES_DATABASE_URL)' ./bin/prisoner --store postgres server

run-server-memory: build
	./bin/prisoner server

# Runs the server against STORE (memory, file, or postgres; default memory),
# e.g. `make run-server STORE=file` or `make dev STORE=postgres`.
run-server: build
ifeq ($(STORE),postgres)
run-server: start-db
endif
	$(if $(filter postgres,$(STORE)),DATABASE_URL='$(POSTGRES_DATABASE_URL)') ./bin/prisoner --store $(STORE) server

run-web:
	cd web && npm install && npm run dev

dev:
	$(MAKE) -j2 run-server run-web STORE=$(STORE)
