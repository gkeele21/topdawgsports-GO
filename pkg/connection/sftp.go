package connection

import (
	"fmt"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// NewSFTPClient gets a client to use for SFTP operations with aes128-cbc enabled
func NewSFTPClient(host, username, password string) (*sftp.Client, error) {
	var sshCfg ssh.Config
	sshCfg.SetDefaults()
	sshCfg.Ciphers = append(sshCfg.Ciphers, "aes128-cbc")

	config := &ssh.ClientConfig{
		Config: sshCfg,
		User:   username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, fmt.Errorf("can't dial %s", err)
	}

	ftpClient, err := sftp.NewClient(client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("couldn't get ftpClient %s", err)
	}

	return ftpClient, nil
}
