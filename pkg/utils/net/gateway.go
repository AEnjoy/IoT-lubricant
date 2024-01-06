package net

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
)

const DefaultGateway = "0.0.0.0"

func GetGateway() (string, error) {
	file, err := os.Open("/proc/net/route")
	if err != nil {
		logger.Errorln("Error opening route file:", err)
		return DefaultGateway, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// 第一个字段是网络接口名称，第二个字段是目标网络
		// 如果目标网络是0.0.0.0 (0x00000000)，表示是默认路由
		if fields[1] == "00000000" {
			// 第三列是网关，值是以十六进制表示的反向字节序
			gatewayHex := fields[2]

			// 将十六进制反转为正确的字节顺序
			gatewayIP := fmt.Sprintf("%d.%d.%d.%d",
				toDecimal(gatewayHex[6:8]),
				toDecimal(gatewayHex[4:6]),
				toDecimal(gatewayHex[2:4]),
				toDecimal(gatewayHex[0:2]),
			)
			return gatewayIP, nil
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Error("Error reading route file:", err)
	}
	return DefaultGateway, fmt.Errorf("no default gateway found")
}
func toDecimal(hexStr string) int {
	var result int
	_, err := fmt.Sscanf(hexStr, "%x", &result)
	if err != nil {
		return 0
	}
	return result
}
