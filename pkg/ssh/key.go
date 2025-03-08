package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
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
func GetLocalSSHPrivateKey() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", exception.ErrNewException(err, exceptionCode.ErrorIO, exception.WithMsg("failed to get user information"))
	}

	sshDir := filepath.Join(user.HomeDir, ".ssh")
	keyFiles := []string{"id_rsa", "id_ecdsa", "id_ed25519"}

	for _, file := range keyFiles {
		privateKeyFile := filepath.Join(sshDir, file)
		content, err := os.ReadFile(privateKeyFile)
		if err == nil {
			return string(content), nil
		}
		if !os.IsNotExist(err) {
			return "", exception.ErrNewException(err, exceptionCode.ErrorIO, exception.WithMsg(fmt.Sprintf("failed to read %s", privateKeyFile)))
		}
	}

	return "", exception.New(exceptionCode.ErrorIO, exception.WithMsg("no private key file found"))
}

// SaveSSHKey 保存ssh密钥到本机.ssh/ name 是"id_rsa", "id_ecdsa", "id_ed25519" 之一
func SaveSSHKey(pri, pub, name string) error {
	user, err := user.Current()
	if err != nil {
		return exception.ErrNewException(err, exceptionCode.ErrorIO, exception.WithMsg("failed to get user information"))
	}
	sshDir := filepath.Join(user.HomeDir, ".ssh")

	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		if err := os.MkdirAll(sshDir, 0700); err != nil {
			return exception.ErrNewException(err, exceptionCode.ErrorIO, exception.WithMsg("failed to create .ssh directory"))
		}
	}
	if err := os.WriteFile(filepath.Join(sshDir, name), []byte(pri), 0600); err != nil {
		return exception.ErrNewException(err, exceptionCode.ErrorIO, exception.WithMsg("failed to save private key"))
	}
	if err := os.WriteFile(filepath.Join(sshDir, name+".pub"), []byte(pub), 0600); err != nil {
		return exception.ErrNewException(err, exceptionCode.ErrorIO, exception.WithMsg("failed to save public key"))
	}
	return nil
}
func CreateSSHKeyBySshKeygen() error {
	if _, err := os.Stat("~/.ssh/id_ed25519"); err == nil {
		return nil
	}

	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", "~/.ssh/id_ed25519", "-N", "", "-q")
	if err := cmd.Run(); err != nil {
		return exception.ErrNewException(err, exceptionCode.ErrorIO, exception.WithMsg("failed to create ssh key"))
	}
	return nil
}
