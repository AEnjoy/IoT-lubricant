package ssh

import (
	"errors"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"golang.org/x/crypto/ssh"
)

var _ RemoteClient = (*client)(nil)

type client struct {
	*model.GatewayHost
	sshClient *ssh.Client
	session   *ssh.Session
}

func (c *client) Close() error {
	return errors.Join(c.session.Close(), c.sshClient.Close())
}
