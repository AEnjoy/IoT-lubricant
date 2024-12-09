package exception

import "fmt"

type ErrLogInfo struct {
	User    string
	Agent   string
	Message error
}

func (r *ErrLogInfo) String() string {
	return fmt.Sprintf("User:%s 's Agent:%s report an error:%s ", r.User, r.Agent, r.Message.Error())
}
