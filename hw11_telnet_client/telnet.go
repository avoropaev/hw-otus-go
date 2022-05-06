package main

import (
	"errors"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func (c *telnetClient) Close() error {
	return c.conn.Close()
}

func (c *telnetClient) Send() error {
	if c.conn == nil {
		return errors.New("need to connect first")
	}

	_, err := io.Copy(c.conn, c.in)
	if err != nil {
		return err
	}

	if _, err = os.Stderr.Write([]byte("...EOF\n")); err != nil {
		return err
	}

	return nil
}

func (c *telnetClient) Receive() error {
	if c.conn == nil {
		return errors.New("need to connect first")
	}

	_, err := io.Copy(c.out, c.conn)
	if err != nil {
		return err
	}

	if _, err = os.Stderr.Write([]byte("...Connection was closed by peer\n")); err != nil {
		return err
	}

	return nil
}

func (c *telnetClient) Connect() error {
	if c.conn != nil {
		return nil
	}

	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}
