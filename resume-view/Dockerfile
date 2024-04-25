FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .. .

WORKDIR /app/cmd/

RUN go build -o /build main.go

FROM alpine:latest

WORKDIR /app

RUN apk add curl make

RUN curl -fsSL \
        https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
        sh

COPY --from=build /build /app/build
COPY --from=build /app/migrations /app/migrations

COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait

CMD /wait && /app/build