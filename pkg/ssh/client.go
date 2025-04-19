package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	def "github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

var _ RemoteClient = (*client)(nil)

type client struct {
	model     *model.GatewayHost
	sshClient *ssh.Client
	session   *ssh.Session
}

func (c *client) Close() error {
	return errors.Join(c.session.Close(), c.sshClient.Close())
}
func (c *client) UploadFiles(paths []string, target string) error {
	sftpClient, err := sftp.NewClient(c.sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	for _, path := range paths {
		func(path string) {
			file, err := os.Open(path)
			if err != nil {
				logger.Error("failed to open file:", err.Error(), path, "ignore")
				return
			}
			defer file.Close()

			remoteFile, err := sftpClient.Create(fmt.Sprintf("%s/%s", target, file.Name()))
			if err != nil {
				logger.Error("failed to create remote file:", err.Error(), target, "ignore")
				return
			}
			defer remoteFile.Close()

			_, err = remoteFile.ReadFrom(file)
			if err != nil {
				logger.Error("failed to upload file:", err.Error(), path, "ignore")
				return
			}
		}(path)
	}
	return nil
}
func (c *client) IoUploadFile(reader io.ReadCloser, target string) error {
	sftpClient, err := sftp.NewClient(c.sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	createFile, err := sftpClient.Create(target)
	if err != nil {
		return err
	}
	defer createFile.Close()

	_, err = io.Copy(createFile, reader)
	return err
}
func (c *client) Download(target, local string) error {
	sftpClient, err := sftp.NewClient(c.sshClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	targetFile, err := sftpClient.Open(target)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	localFile, err := os.Create(local)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = targetFile.WriteTo(localFile)
	return err
}

func (c *client) DeployGateway(hostinfo *model.ServerInfo) error {
	resp, err := http.Get(def.GatewayDeployScripts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	hostInfoByte, err := yaml.Marshal(hostinfo)
	if err != nil {
		return err
	}

	if err = c.IoUploadFile(io.NopCloser(bytes.NewBuffer(hostInfoByte)), "/tmp/lubricant_server_config.yaml"); err != nil {
		return err
	}

	if err = c.IoUploadFile(resp.Body, "/tmp/lubricant_gateway.sh"); err != nil {
		return err
	}
	var out string
	var exitCode int
	out, exitCode, err = c.executeCommandAuto("bash /tmp/lubricant_gateway.sh init /tmp/lubricant_server_config.yaml")
	if err != nil || exitCode != 0 {
		return fmt.Errorf("failed to deploy gateway: %s, exit code: %d", out, exitCode)
	}
	return nil
}
func (c *client) UpdateConfig(hostinfo *model.ServerInfo) error {
	hostInfoByte, err := yaml.Marshal(hostinfo)
	if err != nil {
		return err
	}

	if err = c.IoUploadFile(io.NopCloser(bytes.NewBuffer(hostInfoByte)), "/tmp/lubricant_server_config.yaml"); err != nil {
		return err
	}

	// todo: 需要支持自定义conf路径
	out, exitCode, err := c.executeCommandAuto("cp /tmp/lubricant_server_config.yaml /opt/lubricant/gateway/")
	if err != nil || exitCode != 0 {
		return fmt.Errorf("failed to update config: %s, exit code: %d", out, exitCode)
	}

	out, exitCode, err = c.executeCommandAuto("systemctl restart Lubricant-Gateway")
	if err != nil || exitCode != 0 {
		return fmt.Errorf("failed to restart gateway: %s, exit code: %d", out, exitCode)
	}
	return nil
}

func (c *client) GetConfig() (ret *model.ServerInfo, err error) {
	ret = new(model.ServerInfo)
	out, exitCode, err := c.executeCommandAuto("cat /opt/lubricant/gateway/lubricant_server_config.yaml")
	if err != nil || exitCode != 0 {
		return nil, fmt.Errorf("failed to get config: %s, exit code: %d", out, exitCode)
	}
	if err = yaml.Unmarshal([]byte(out), ret); err != nil {
		return nil, err
	}
	return ret, nil
}
