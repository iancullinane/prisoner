.PHONY: build docker

build:
	mkdir -p bin
	go build -o bin/prisoner .

docker:
	docker build -t prisoner .
