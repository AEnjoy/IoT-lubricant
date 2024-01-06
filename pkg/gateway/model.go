package gateway

import "sync"

type uploadModel struct {
	gatewayId string
	edgeId    string
	data      [][]byte
	time      []string // ex.2023-01-01 00:00:00
	len       int
}

var dataSet = sync.Map{} // map[edgeId] []uploadModel
