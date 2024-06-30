package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type simpleTelnetClient struct {
	conn    net.Conn
	in      io.Reader
	out     io.Writer
	address string
	timeout time.Duration
	cancel  context.CancelFunc
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &simpleTelnetClient{in: in, out: out, address: address, timeout: timeout}
}

func (c *simpleTelnetClient) Connect() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	c.cancel = cancel
	dialer := &net.Dialer{Timeout: c.timeout}
	if c.conn, err = dialer.DialContext(ctx, "tcp", c.address); err != nil {
		return fmt.Errorf("соединение не удалось: %w", err)
	}
	return nil
}

func (c *simpleTelnetClient) Close() error {
	c.cancel()
	return c.conn.Close()
}

func (c *simpleTelnetClient) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		if _, err := c.conn.Write([]byte(scanner.Text() + "\n")); err != nil {
			return fmt.Errorf("передача данных не удалась: %w", err)
		}
	}
	if scanner.Err() != nil {
		return fmt.Errorf("ошибка сканирования при передаче: %w", scanner.Err())
	}
	return nil
}

func (c *simpleTelnetClient) Receive() error {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		if _, err := fmt.Fprintln(c.out, scanner.Text()); err != nil {
			return fmt.Errorf("приём данных не удался: %w", err)
		}
	}
	if scanner.Err() != nil {
		return fmt.Errorf("ошибка сканирования при приеме: %w", scanner.Err())
	}
	return nil
}
