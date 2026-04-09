FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/api .
COPY migrations ./migrations
RUN mkdir -p /app/storage/pdfs
RUN chown -R 1000:1000 /app/storage
EXPOSE 8080
CMD ["./api"]
