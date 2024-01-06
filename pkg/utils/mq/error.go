package mq

import "errors"

var (
	ErrNotFoundTopic      = errors.New("topic not found")
	ErrNotFoundSubscriber = errors.New("subscriber not found")
	ErrHasBeenSubscribed  = errors.New("has been subscribed")
)
