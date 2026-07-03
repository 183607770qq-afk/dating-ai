FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/api cmd/api/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o bin/scraper cmd/scraper/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin /app
COPY --from=builder /app/config.yaml /app
COPY --from=builder /app/web/public /app/web/public

EXPOSE 8080

ENV GIN_MODE=release

CMD ["sh"]
