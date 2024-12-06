package uds

import (
	"context"
	"log"
	"net"
	"os"
	"time"
)

type UdsServer struct {
	addr string
	buf  []byte
}

func NewUdsServer(addr string) *UdsServer {
	return &UdsServer{addr: addr, buf: make([]byte, 1024)}
}

func (server *UdsServer) GetBuf() string {
	return string(server.buf)
}

func (server *UdsServer) Run(ctx context.Context) error {
	os.Remove(server.addr)

	conn, err := net.ListenPacket("unixgram", server.addr)

	// go
	t := time.Now().Add(time.Second * 20)

	// c++思路
	t = time.Now()
	t.Add(time.Second * 20)

	conn.SetReadDeadline(t)

	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		n, err := server.read(conn)
		if err != nil {
			log.Println("read error: ", err)
			continue
		}
		if n <= 0 {
			log.Println("read size is: ", n)
			continue
		}
	}
}

func (server *UdsServer) read(conn net.PacketConn) (int, error) {
	n, _, err := conn.ReadFrom(server.buf)
	if err != nil {
		return 0, err
	}
	return n, nil
}
