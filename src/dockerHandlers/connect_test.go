/*
author Vyacheslav Kryuchenko
*/
package dockerHandlers

import (
	"testing"
)

func TestCreateClientConnection(t *testing.T) {
	const (
		dockerhost = "tcp://192.168.56.101:2376"
		apiversion = "1.32"
	)
	testClient, err := CreateClientConnection(dockerhost, apiversion)
	defer testClient.Close()
	if err != nil {
		t.Error(err)
		return
	}
}
