FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
WORKDIR /app/cmd/core
RUN go build -o /core

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /core .

# COPY scripts/files ./files

EXPOSE 8080
CMD ["./core"]
