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
)

func (a Api) _getOperator(c *gin.Context) operator {
	o := c.Query("operator")
	o = strings.ToLower(o)
	switch o {
	// start-agent,stop-agent,start-gather,stop-agent,get-openapidoc
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
	default:
		return unknownOperator
	}
}
