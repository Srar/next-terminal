package proxy

import (
	"fmt"
	"net"
)

type Type string

const (
	TypeNone   Type = ""
	TypeSocks5 Type = "socks5"
	TypeHTTP   Type = "http"
	TypeHTTPS  Type = "https"
	TypeSSH    Type = "ssh"
)

func (p Type) Valid() bool {
	switch p {
	case TypeNone, TypeSocks5, TypeHTTP, TypeSSH, TypeHTTPS:
		return true
	}
	return false
}

func Dial(t Type, c *Config) (net.Conn, error) {
	switch t {
	case TypeSocks5:
		return DialWithSocks5(c)
	case TypeHTTP:
		return DialWithHTTP(c)
	case TypeHTTPS:
		return DialWithHTTPS(c)
	case TypeSSH:
		return DialWithSSH(c)
	}
	return nil, fmt.Errorf("unsupported proxy. ")
}
