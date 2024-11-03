package edge

import "time"

type DataPacket struct {
	Data      []byte
	Timestamp time.Time
}
