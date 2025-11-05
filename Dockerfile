FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go mod download

COPY . .

RUN go build -o spy-cats ./cmd/server

EXPOSE 8080
CMD ["./spy-cats"]