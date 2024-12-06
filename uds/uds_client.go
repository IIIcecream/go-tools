package uds

import (
	"errors"
	"log"
	"net"
)

type UdsClient struct {
	conn *net.UnixConn
}

func (c *UdsClient) SendMsg(msg string) error {
	n, err := c.conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	if n != len(msg) {
		log.Fatalf("uds send msg failed, n:%d, msg:%s\n", n, msg)
		return errors.New("uds send msg failed")
	}
	return nil
}

func (c *UdsClient) Close() error {
	return c.conn.Close()
}

func NewUdsClient(addr string) (*UdsClient, error) {
	unixAddr, err := net.ResolveUnixAddr("unixgram", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUnix("unixgram", nil, unixAddr)
	if err != nil {
		return nil, err
	}

	return &UdsClient{conn: conn}, nil
}
