package proxy

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func DialWithSSH(config *Config) (conn net.Conn, err error) {
	sshConfig := ssh.ClientConfig{
		User:            config.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(config.Password)},
		Timeout:         time.Second * 5,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshConnection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), &sshConfig)
	if err != nil {
		return nil, err
	}

	conn, err = sshConnection.Dial("tcp", fmt.Sprintf("%s:%d", config.DialHost, config.DialPort))
	return
}
