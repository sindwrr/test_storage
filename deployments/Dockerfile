FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/app

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=builder /server .
COPY web/templates ./web/templates
COPY web/static ./web/static
COPY migrations ./migrations
COPY .env .env
EXPOSE 8000
CMD ["./server"]
