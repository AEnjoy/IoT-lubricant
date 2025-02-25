#!/usr/bin/env sh
# set -e

GATEWAY_BASE_PATH="/opt/lubricant/gateway"
GATEWAY_SERVICE="/etc/systemd/system/iot-lubricant.service"
GATEWAY_SERVICERC="/etc/init.d/iot-lubricant"
GITHUB_URL="github.com"

_version="v0.0.1"
os_arch=""
command=$1
configFile=$2

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
plain='\033[0m'
export PATH="$PATH:/usr/local/bin"
err() {
    printf "${red}%s${plain}\n" "$*" >&2
}

success() {
    printf "${green}%s${plain}\n" "$*"
}

info() {
    printf "${yellow}%s${plain}\n" "$*"
}

pre_check() {
    umask 077

    ## os_arch
    if uname -m | grep -q 'x86_64'; then
        os_arch="amd64"
    elif uname -m | grep -q 'i386\|i686'; then
        os_arch="386"
    elif uname -m | grep -q 'aarch64\|armv8b\|armv8l'; then
        os_arch="arm64"
    elif uname -m | grep -q 'arm'; then
        os_arch="arm"
    elif uname -m | grep -q 's390x'; then
        os_arch="s390x"
    elif uname -m | grep -q 'riscv64'; then
        os_arch="riscv64"
    fi
}

pre_check
GATEWAY_DOWNLOAD_URL="https://github.com/aenjoy/iot-lubricant/releases/download/${_version}/lubricant-gateway-linux-${os_arch}.zip"

if [ "$command" = "install" ]; then
  sudo mkdir -p $GATEWAY_BASE_PATH
  _cmd="wget -t 2 -T 60 -O lubricant-gateway-linux-${os_arch}.zip $GATEWAY_DOWNLOAD_URL >/dev/null 2>&1"
  if ! eval "$_cmd"; then
      err "Release 下载失败，请检查本机能否连接 ${GITHUB_URL}"
      return 1
  fi
  unzip -qo lubricant-gateway-linux-${os_arch}.zip
  mv lubricant-gateway-linux-${os_arch} $GATEWAY_BASE_PATH/
  mv $configFile $GATEWAY_BASE_PATH/
  rf lubricant-gateway-linux-${os_arch}.zip
fi
