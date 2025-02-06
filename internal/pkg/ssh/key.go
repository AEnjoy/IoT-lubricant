package ssh

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception"
	exceptionCode "github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
)

func GetLocalSSHPublicKey() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", exception.ErrNewException(err, exceptionCode.ErrorIO, exception.WithMsg("failed to get user information"))
	}

	sshDir := filepath.Join(user.HomeDir, ".ssh")
	publicKeyFiles := []string{"id_rsa.pub", "id_ecdsa.pub", "id_ed25519.pub"}

	for _, file := range publicKeyFiles {
		publicKeyFile := filepath.Join(sshDir, file)
		content, err := os.ReadFile(publicKeyFile)
		if err == nil {
			return string(content), nil
		}
		if !os.IsNotExist(err) {
			return "", exception.ErrNewException(err, exceptionCode.ErrorIO, exception.WithMsg(fmt.Sprintf("failed to read %s", publicKeyFile)))
		}
	}

	return "", exception.New(exceptionCode.ErrorIO, exception.WithMsg("no public key file found"))
}
