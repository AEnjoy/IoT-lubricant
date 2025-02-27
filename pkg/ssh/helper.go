package ssh

import "fmt"

func (c *client) executeCommand(cmd string) (string, int, error) {
	output, err := c.session.CombinedOutput(cmd)
	if err != nil {
		return string(output), 1, err
	}
	return string(output), 0, nil
}

func (c *client) executeCommandSudo(cmd string) (string, int, error) {
	if c.model.PassWd == "" {
		return "", 1, fmt.Errorf("no password set")
	}
	//  echo "passwd" | sudo -S cmd
	return c.executeCommand(fmt.Sprintf("echo \"%s\" | sudo -S %s", c.model.PassWd, cmd))
}
func (c *client) executeCommandAuto(cmd string) (string, int, error) {
	if c.model.UserName == "root" {
		return c.executeCommand(cmd)
	} else {
		return  c.executeCommandSudo(cmd)
	}
}
