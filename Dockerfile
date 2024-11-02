FROM golang:1.23-alpine AS builder

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server cmd/main.go && \
    go clean -cache -modcache


FROM alpine:latest

COPY --from=builder /server /server

COPY ./internal/migrations/ /internal/migrations

CMD ["/server"]