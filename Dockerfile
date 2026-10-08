FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o social-platform \
    ./cmd/api


FROM alpine:3.22

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/social-platform .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./social-platform"]