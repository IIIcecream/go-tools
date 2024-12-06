package uds

import (
	"log"
	"testing"
)

func TestUdsClient(t *testing.T) {
	uds_client, err := NewUdsClient("/tmp/uds.sock")
	if err != nil {
		log.Fatalln(err.Error())
		t.FailNow()
		return
	}
	defer uds_client.Close()

	uds_client.SendMsg("hello, world")
}
