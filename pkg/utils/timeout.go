package utils

import (
	"context"
	"time"
)

const (
	// s
	DefaultTimeout_Oper time.Duration = 3
	DefaultTimeout_Req  time.Duration = 30
)

// CreateTimeOutContext 创建超时的context 支持当ctx为nil时，使用默认的context
func CreateTimeOutContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if ctx == nil {
		return context.WithTimeout(context.Background(), timeout*time.Second)
	}
	return context.WithTimeout(ctx, timeout*time.Second)
}
