FROM golang:1.15-alpine3.12

RUN mkdir /banner_roulette_backend

WORKDIR /banner_roulette_backend

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

CMD ./main


