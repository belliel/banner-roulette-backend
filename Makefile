.PHONY: build, mongo, container, swagger
.SILENT:

build:
	go build -o main ./cmd/app/main.go

mongo:
	docker run --name some-mongo \
		-e MONGO_INITDB_ROOT_USERNAME=mongoadmin \
		-e MONGO_INITDB_ROOT_PASSWORD=secret \
		-d mongo:4.2-bionic


container:
	docker build --file ./Dockerfile --tag banner_roulette_backend:latest ./;
	docker stop $(docker ps -q) || true
	docker rm $(docker ps -aq) || true
	docker run -d \
 		--volume="$(PWD)/assets:/banner_roulette_backend/assets" banner_roulette_backend

swagger:
	swag init -g internal/app/app.go

