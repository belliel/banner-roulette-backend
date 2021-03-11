.PHONY: build, local
.SILENT:

build:
	go build -o main ./cmd/app/main.go

local:
	docker run --name some-mongo -d mongo:tag
	docker run --name some-redis -d redis
	build
	./main

