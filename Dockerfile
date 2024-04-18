FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/view

RUN go build -o resumeview

EXPOSE 3007

CMD ["./resumeview"]