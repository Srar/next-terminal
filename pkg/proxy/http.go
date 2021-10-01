package proxy

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"time"
)

func DialWithHTTP(config *Config) (conn net.Conn, err error) {
	conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), time.Second * 5)
	if err != nil {
		return
	}
	conn.SetDeadline(time.Now().Add(time.Second * 30))

	host := fmt.Sprintf("%s:%d", config.DialHost, config.DialPort)
	conn.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n", host)))
	conn.Write([]byte(fmt.Sprintf("Host: %s\r\n", host)))
	if config.Username != "" || config.Password != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.Username, config.Password)))
		conn.Write([]byte(fmt.Sprintf("Proxy-Authorization: basic %s\r\n", encoded)))
	}
	conn.Write([]byte("\r\n"))

	br := bufio.NewReader(conn)
	response, err := http.ReadResponse(br, nil)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http proxy error: %d ", response.StatusCode)
	}

	conn.SetDeadline(time.Time{})
	return
}
