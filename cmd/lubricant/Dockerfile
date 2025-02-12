FROM golang:1.23-bookworm  AS builder

ARG VERSION
ARG BUILD_TIME
ARG GIT_COMMIT
ARG FEATURES
ARG BUILD_HOST_PLATFORM
ARG PLATFORM_VERSION

RUN apt-get update && \
    apt-get install -y wget build-essential

RUN wget https://static.iotroom.top/TDengine-client-3.3.5.2-Linux-x64.tar.gz -O /tmp/tdengine.tar.gz && \
    tar -xzf /tmp/tdengine.tar.gz -C /tmp/ && \
    cd /tmp/TDengine-client-3.3.5.2 && \
    bash install_client.sh

ENV CGO_ENABLED=1

WORKDIR /app

# if you are in China, you can uncomment the following line
# RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download
COPY . .

WORKDIR /app/cmd/lubricant
RUN GO_VERSION=`go version | awk '{print $3}'` && CGO_ENABLED=0 &&\
	go build -v -tags=sonic -tags=avx -ldflags "\
	-w -s \
	-X 'main.Version=${VERSION}' \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.GoVersion=$GO_VERSION' \
	-X 'main.GitCommit=${GIT_COMMIT}' \
	-X 'main.Features=${FEATURES}' \
	-X 'main.BuildHostPlatform=${BUILD_HOST_PLATFORM}' \
	-X 'main.PlatformVersion=${PLATFORM_VERSION}' \
	" -o /lubricant

FROM debian:bookworm-slim

COPY --from=builder /tmp/tdengine.tar.gz /tmp/
RUN apt-get update && \
    apt-get install -y tar && \
    tar -xzf /tmp/tdengine.tar.gz -C /tmp/ && \
    cd /tmp/TDengine-client-3.3.5.2 && \
    bash install_client.sh && \
    rm -rf /tmp/TDengine-client-3.3.5.2* /tmp/tdengine.tar.gz && \
    apt-get purge -y tar && \
    apt-get autoremove -y && \
    apt-get clean

WORKDIR /root/
COPY --from=builder /lubricant .

# COPY scripts/files ./files

EXPOSE 8080
CMD ["./lubricant"]
