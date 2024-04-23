FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/

RUN go build -o /build main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /build /app/build
COPY --from=build /app/cmd/config.yml /app

CMD ["/app/build"]