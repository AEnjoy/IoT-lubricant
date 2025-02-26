FROM golang:1.23-alpine AS builder

ARG VERSION
ARG BUILD_TIME
ARG GIT_COMMIT
ARG FEATURES
ARG BUILD_HOST_PLATFORM
ARG PLATFORM_VERSION

WORKDIR /app
COPY . .
# if you are in China, you can uncomment the following line
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
WORKDIR /app/cmd/agent
RUN GO_VERSION=`go version | awk '{print $3}'` && \
	go build -v -tags=sonic -tags=avx -ldflags "\
	-w -s \
	-X 'main.Version=${VERSION}' \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.GoVersion=$GO_VERSION' \
	-X 'main.GitCommit=${GIT_COMMIT}' \
	-X 'main.Features=${FEATURES}' \
	-X 'main.BuildHostPlatform=${BUILD_HOST_PLATFORM}' \
	-X 'main.PlatformVersion=${PLATFORM_VERSION}' \
    " -o /agent

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /agent .
EXPOSE 5436

# COPY scripts/files ./files

CMD ["./agent"]

