package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	def "github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
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
