package ssh

import "golang.org/x/crypto/ssh"

func executeCommand(session *ssh.Session, cmd string) (string, int, error) {
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), 1, err
	}
	return string(output), 0, nil
}
