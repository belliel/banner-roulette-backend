.PHONY: build, local, swag
.SILENT:

build:
	go build -o main ./cmd/app/main.go

local:
	docker run --name some-mongo -d mongo:tag
	build
	./main

container:
	docker build --file ./Dockerfile --tag banner_roulette_backend:latest ./;
	docker stop $(docker ps -q) || true
	docker rm $(docker ps -aq) || true
	docker run -d \
 		--volume="$(pwd)/assets:/banner_roulette_backend/assets" banner_roulette_backend




swagger:
	swag init -g internal/app/app.go

