FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/echo-service/cmd/

RUN go build -o /build main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /build /app/build

CMD ["/app/build"]