package compiler

import (
	"testing"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"
)

func TestBuildAndStartCompiler(t *testing.T) {
	cli, err := client.Connect("unix:///var/run/docker.sock", "v1.40")
	if err != nil {
		t.Error(err)
		return
	}
	config := NewContainerConfig()
	config.PortMap.Insert("127.0.0.1", "23366", "23367")
	BuildAndStartCompiler("compiler", cli, config)
}
