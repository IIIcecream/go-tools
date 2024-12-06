package uds

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestUdsServer(t *testing.T) {
	server := NewUdsServer("/tmp/uds.sock")

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()

		server.Run(ctx)
		s := server.GetBuf()
		if len(s) == 0 {
			fmt.Println("buf is empty")
		}
		fmt.Println(s)
	}()

	time.Sleep(time.Second * 1)

	cancel()
	wg.Wait()
}
