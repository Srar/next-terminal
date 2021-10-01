package proxy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	authMethodNone             = byte(0x00)
	authMethodUsernamePassword = byte(0x02)
)

func DialWithSocks5(config *Config) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), time.Second*5)
	if err != nil {
		return nil, err
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	defer func() {
		if err != nil && conn != nil {
			conn.Close()
		}
	}()

	authenticationMethod := authMethodNone
	if config.Username != "" || config.Password != "" {
		authenticationMethod = authMethodUsernamePassword
	}
	_, err = conn.Write([]byte{0x05, 0x01, authenticationMethod})
	if err != nil {
		return nil, err
	}

	authenticationResponse := make([]byte, 2)
	n, err := conn.Read(authenticationResponse)
	if err != nil {
		return nil, err
	}
	if n != 2 {
		return nil, fmt.Errorf("failed to handshake with server. ")
	}
	if authenticationResponse[1] == 0xff {
		return nil, fmt.Errorf("no acceptable authentication method. ")
	}

	if authenticationMethod == authMethodUsernamePassword {
		b := bytes.Buffer{}
		b.Write([]byte{0x01})
		b.Write([]byte{byte(len(config.Username))})
		b.Write([]byte(config.Username))
		b.Write([]byte{byte(len(config.Password))})
		b.Write([]byte(config.Password))
		conn.Write(b.Bytes())

		 n, err = conn.Read(authenticationResponse)
		 if err != nil {
			 return nil, err
		 }
		 if n != 2 {
			 return nil, fmt.Errorf("failed to handshake with server. ")
		 }
		 if authenticationResponse[1] != 0x00 {
			 return nil, fmt.Errorf("failed to handshake due to wrong username or password. ")
		 }
	}

	parsedIPv4 := net.ParseIP(config.DialHost)
	if parsedIPv4 != nil {
		parsedIPv4 = parsedIPv4.To4()
	}

	var command []byte
	if parsedIPv4 == nil {
		command = make([]byte, 5+len(config.DialHost)+2)
		command[0] = 0x05
		command[1] = 0x01
		command[2] = 0x00
		command[3] = 0x03
		command[4] = byte(len(config.DialHost))
		copy(command[5:], config.DialHost)
		binary.BigEndian.PutUint16(command[5+len(config.DialHost):], uint16(config.DialPort))
	} else {
		command = make([]byte, 4+4+2)
		command[0] = 0x05
		command[1] = 0x01
		command[2] = 0x00
		command[3] = 0x01
		copy(command[4:], parsedIPv4)
		binary.BigEndian.PutUint16(command[4+4:], uint16(config.DialPort))
	}
	_, err = conn.Write(command)
	if err != nil {
		return nil, err
	}

	commandResponse := make([]byte, 10)
	n, err = conn.Read(commandResponse)
	if err != nil {
		return nil, err
	}
	if n != 10 {
		return nil, fmt.Errorf("failed to command with server. ")
	}
	if commandResponse[0] != 0x05 || commandResponse[1] != 0x00 || commandResponse[2] != 0x00 || commandResponse[3] != 0x01 {
		return nil, fmt.Errorf("the server response that can't connect to remote. ")
	}

	conn.SetReadDeadline(time.Time{})
	return conn, nil
}
