.PHONY: build, mongo, container swag
.SILENT:

build:
	go build -o main ./cmd/app/main.go

mongo:
	docker run --name some-mongo -d mongo:tag


container:
	docker build --file ./Dockerfile --tag banner_roulette_backend:latest ./;
	docker stop $(docker ps -q) || true
	docker rm $(docker ps -aq) || true
	docker run -d \
 		--volume="$PWD/assets:/banner_roulette_backend/assets" banner_roulette_backend

swagger:
	swag init -g internal/app/app.go

