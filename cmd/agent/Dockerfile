FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY . .
# if you are in China, you can uncomment the following line
# RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
WORKDIR /app/cmd/agent
RUN go build -o /agent

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /agent .

# COPY scripts/files ./files

CMD ["./agent"]

