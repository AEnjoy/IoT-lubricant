package utils

func IsClosed[T any](ch <-chan T) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}

func CloseChannel[T any](ch chan T) {
	if IsClosed(ch) {
		return
	}
	defer func() {
		_ = recover()
	}()
	close(ch)
}
func ConvertRecvToSend[T any](recvChan <-chan T) chan T {
	sendChan := make(chan T)

	go func() {
		defer close(sendChan)
		for val := range recvChan {
			sendChan <- val
		}
	}()

	return sendChan
}
