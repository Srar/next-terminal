package proxy

import (
	"fmt"
	"net"
)

type Type string

const (
	TypeNone = ""
	TypeSocks5 = "socks5"
)

func (p Type) Valid() bool {
	switch p {
	case TypeNone, TypeSocks5:
		return true
	}
	return false
}

func Dial(t Type, c *Config) (net.Conn, error)  {
	switch t {
	case TypeSocks5:
		return DialWithSocks5(c)
	}
	return nil, fmt.Errorf("unsupported proxy. ")
}
