FROM golang:1.25-alpine3.22 AS build

WORKDIR /build

COPY go.mod go.sum main.go .

RUN go mod download

COPY . .

RUN go build -o user-service .

FROM alpine:3.22 AS app

WORKDIR /app

COPY --from=build /build/user-service /app

COPY configs/config-docker.yaml /etc/user-service/config.yaml

ENTRYPOINT ["/app/user-service", "server"]
