FROM golang:1.22.2-alpine

ENV CGO_ENABLED=1

WORKDIR /app

COPY go.mod go.sum ./

RUN apk add --no-cache gcc musl-dev
RUN go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd
COPY ./config ./config
COPY ./migrations ./migrations

RUN go build -o security_service ./cmd/security_service

RUN go test ./...

EXPOSE 8080

CMD ["./security_service", "--config=./config/dev.yaml"]