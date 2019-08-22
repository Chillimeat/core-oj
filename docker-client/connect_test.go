package dockerx

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	cli, err := Connect("unix:///var/run/docker.sock", "v1.40")
	if err != nil {
		t.Error(err)
		return
	}

	containers, err := ListAllContainers(cli)
	if err != nil {
		t.Error(err)
		return
	}
	for _, container := range containers {
		fmt.Println(container.ID)
	}
}
