package dockerHandlers

import (
	"log"
	"testing"
)

func TestList(t *testing.T) {
	const (
		dockerhost = "tcp://localhost:2376"
		apiversion = "1.32"
	)
	testClient, err := CreateClientConnection(dockerhost, apiversion)
	if err != nil {
		t.Error(err)
		return
	}
	defer testClient.Close()
	containers, err := List(testClient, true)
	if err != nil {
		t.Error(err)
		return
	}
	if len(containers) == 0 {
		t.Error("empty container list")
		return
	}
	log.Print(containers)
}
