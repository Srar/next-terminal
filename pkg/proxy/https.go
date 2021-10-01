package proxy

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"time"
)

func DialWithHTTPS(config *Config) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp",  fmt.Sprintf("%s:%d", config.Host, config.Port), time.Second * 5)
	if err != nil {
		return nil, err
	}
	conn.SetDeadline(time.Now().Add(time.Minute))

	tlsConn := tls.Client(conn, &tls.Config{})
	err = tlsConn.Handshake()
	if err != nil {
		tlsConn.Close()
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", config.DialHost, config.DialPort)
	tlsConn.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n", host)))
	tlsConn.Write([]byte(fmt.Sprintf("Host: %s\r\n", host)))
	if config.Username != "" || config.Password != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.Username, config.Password)))
		tlsConn.Write([]byte(fmt.Sprintf("Proxy-Authorization: basic %s\r\n", encoded)))
	}
	tlsConn.Write([]byte("\r\n"))

	br := bufio.NewReader(tlsConn)
	response, err := http.ReadResponse(br, nil)
	if err != nil {
		tlsConn.Close()
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		tlsConn.Close()
		return nil, fmt.Errorf("http proxy error: %d ", response.StatusCode)
	}

	conn.SetDeadline(time.Time{})

	return tlsConn, nil
}
