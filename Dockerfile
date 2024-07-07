FROM golang:1.22-bullseye as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-s -w" -o app

# Production stage
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=builder /app/app .

CMD ["./app"]

# Local development stage
FROM golang:1.22 as dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

CMD ["air"]
