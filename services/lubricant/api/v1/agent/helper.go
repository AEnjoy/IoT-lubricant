package agent

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type operator uint8

const (
	unknownOperator operator = iota
	startAgent
	stopAgent
	restartAgent
	startGather
	stopGather
	getOpenapiDoc
	getGatherStatus
)

func (a Api) _getOperator(c *gin.Context) operator {
	o := c.Query("operator")
	o = strings.ToLower(o)
	switch o {
	case "start-agent":
		return startAgent
	case "stop-agent":
		return stopAgent
	case "start-gather":
		return startGather
	case "stop-gather":
		return stopGather
	case "get-openapidoc":
		return getOpenapiDoc
	case "get-gather-status":
		return getGatherStatus
	default:
		return unknownOperator
	}
}
