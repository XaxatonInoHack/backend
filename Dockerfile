FROM golang:1.23-alpine AS builder

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server ./cmd/service/main.go
RUN go build -o /prepare ./cmd/prepare/main.go

FROM alpine:latest

COPY --from=builder /server /server
COPY --from=builder /prepare /prepare

COPY ./internal/migrations/ /internal/migrations

CMD ["/server"]