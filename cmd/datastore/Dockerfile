FROM aenjoy/debian-tdengine-driver-gobuilder:latest AS builder

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

WORKDIR /app/cmd/datastore
RUN GO_VERSION=`go version | awk '{print $3}'` && CGO_ENABLED=1 &&\
	go build -v -tags=sonic -tags=avx -ldflags "\
	-w -s \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.Version=${VERSION}' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.BuildTime=${BUILD_TIME}' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.GoVersion=$GO_VERSION' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.GitCommit=${GIT_COMMIT}' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.Features=${FEATURES}' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.BuildHostPlatform=${BUILD_HOST_PLATFORM}' \
	-X 'github.com/aenjoy/iot-lubricant/pkg/version.PlatformVersion=${PLATFORM_VERSION}' \
	" -o /datastore

FROM aenjoy/debian-tdengine-driver:latest
WORKDIR /root/
COPY --from=builder /datastore .

# COPY scripts/files ./files

CMD ["./datastore"]
