FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/resume-view/cmd/

RUN go build -o /build main.go

FROM alpine:latest

WORKDIR /app


COPY --from=build /build /app/build

COPY --from=ghcr.io/ufoscout/docker-compose-wait:latest /wait /wait

CMD /wait && /app/build