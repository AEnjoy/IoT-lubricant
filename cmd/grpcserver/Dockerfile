FROM golang:1.24-alpine AS builder

ARG VERSION
ARG BUILD_TIME
ARG GIT_COMMIT
ARG FEATURES
ARG BUILD_HOST_PLATFORM
ARG PLATFORM_VERSION
ARG CGO_ENABLED

WORKDIR /app

# if you are in China, you can uncomment the following line
# RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /app/cmd/grpcserver
RUN GO_VERSION=`go version | awk '{print $3}'` && CGO_ENABLED=$CGO_ENABLED &&\
	go build -v -tags=sonic -tags=avx -ldflags "\
	-w -s \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.Version=${VERSION}' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.BuildTime=${BUILD_TIME}' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.GoVersion=$GO_VERSION' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.GitCommit=${GIT_COMMIT}' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.Features=${FEATURES}' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.BuildHostPlatform=${BUILD_HOST_PLATFORM}' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.PlatformVersion=${PLATFORM_VERSION}' \
	" -o /grpcserver

FROM alpine:3.21.3
WORKDIR /root/
COPY --from=builder /grpcserver .

# COPY scripts/files ./files

EXPOSE 5423
CMD ["./grpcserver"]
