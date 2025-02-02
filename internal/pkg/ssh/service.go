package ssh

import (
	"errors"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"golang.org/x/crypto/ssh"
)

type RemoteClient interface {
	Close() error
	// todo...
}

func NewSSHClient(gateway *model.GatewayHost) (RemoteClient, error) {
	if gateway.UserName == "" {
		return nil, errors.New("ssh: target host user name is empty")
	}

	var sshConfig = &ssh.ClientConfig{
		User:            gateway.UserName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if gateway.PrivateKey != "" {
		key, err := ssh.ParsePrivateKey([]byte(gateway.PrivateKey))
		if err != nil {
			return nil, err
		}
		sshConfig.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(key),
		}
	} else if gateway.PassWd != "" {
		sshConfig.Auth = []ssh.AuthMethod{
			ssh.Password(gateway.PassWd),
		}
	} else {
		return nil, errors.New("ssh: unknown login method")
	}

	dial, err := ssh.Dial("tcp", gateway.Host, sshConfig)
	if err != nil {
		return nil, err
	}

	var client = &client{
		GatewayHost: gateway,
	}
	session, err := dial.NewSession()
	if err != nil {
		_ = dial.Close()
		return nil, err
	}

	client.sshClient = dial
	client.session = session
	return client, nil
}
