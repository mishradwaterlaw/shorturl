# Stage 1: Builder
FROM golang:1.26 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o shorturl .

# Stage 2: Runner
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/shorturl ./
EXPOSE 8080
CMD ["./shorturl"]