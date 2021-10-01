package guacd

import (
	"fmt"
	"io"
	"net"
	"next-terminal/pkg/proxy"
	"sync"
	"time"
)

func newDisposableProxyForwarder(timeout time.Duration, localPort int, proxyType proxy.Type, proxyConfig *proxy.Config) error {
	forwarderHandler := func(conn net.Conn) {
		proxyConnection, err := proxy.Dial(proxyType, proxyConfig)
		if err != nil {
			return
		}
		defer proxyConnection.Close()
		defer conn.Close()

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			io.Copy(proxyConnection, conn)
			wg.Done()
		}()

		go func() {
			io.Copy(conn, proxyConnection)
			wg.Done()
		}()

		wg.Wait()
	}

	err := newDisposableForwarder(timeout, localPort, forwarderHandler)
	if err != nil {
		return err
	}

	return nil
}


// newDisposableForwarder 一次性转发器
func newDisposableForwarder(timeout time.Duration, localPort int, handler func(conn net.Conn)) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		return err
	}

	waiting := true
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		waiting = false
		handler(conn)
		listener.Close()
	}()

	go func() {
		time.Sleep(timeout)
		if waiting {
			listener.Close()
		}
	}()

	return nil
}