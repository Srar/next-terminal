package term

import (
	"fmt"
	"next-terminal/pkg/proxy"
	"time"

	"golang.org/x/crypto/ssh"
)

func NewSshClient(ip string, port int, proxyType proxy.Type, proxyConfig *proxy.Config, username, password, privateKey, passphrase string) (*ssh.Client, error) {
	var authMethod ssh.AuthMethod
	if username == "-" || username == "" {
		username = "root"
	}
	if password == "-" {
		password = ""
	}
	if privateKey == "-" {
		privateKey = ""
	}
	if passphrase == "-" {
		passphrase = ""
	}

	var err error
	if privateKey != "" {
		var key ssh.Signer
		if len(passphrase) > 0 {
			key, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(passphrase))
			if err != nil {
				return nil, err
			}
		} else {
			key, err = ssh.ParsePrivateKey([]byte(privateKey))
			if err != nil {
				return nil, err
			}
		}
		authMethod = ssh.PublicKeys(key)
	} else {
		authMethod = ssh.Password(password)
	}

	config := &ssh.ClientConfig{
		Timeout:         5 * time.Second,
		User:            username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", ip, port)

	// 不使用代理
	if proxyType == "" {
		return ssh.Dial("tcp", addr, config)
	}

	// 使用代理
	proxyConnection, err := proxy.Dial(proxyType, proxyConfig)
	if err != nil {
		return nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(proxyConnection, addr, config)
	if err != nil {
		return nil, err
	}
	return ssh.NewClient(c, chans, reqs), nil
}
