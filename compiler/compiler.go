package compiler

import (
	"context"
	"fmt"
	"time"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"

	"github.com/Myriad-Dreamin/core-oj/types"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"google.golang.org/grpc"
)

type Compiler struct {
	cli       *client.Client
	container *dockertypes.Container
	conn      *grpc.ClientConn
	rpcx.CompilerClient
}

func StartCompiler(containerInfo *dockertypes.Container, cli *client.Client, cconfig *types.CompilerConfig, config *client.ContainerConfig) (cp *Compiler, err error) {
	if containerInfo.Status != "running" {
		err = cli.ContainerStart(context.Background(), containerInfo.ID, dockertypes.ContainerStartOptions{})
		if err != nil {
			return nil, err
		}
	}

	cp = new(Compiler)

	fmt.Printf("Container %s is started\n", containerInfo.ID)
	cp.cli = cli
	cp.container = containerInfo
	cp.conn, err = grpc.Dial(cconfig.GrpcAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	cp.CompilerClient = rpcx.NewCompilerClient(cp.conn)
	return
}

// BuildAndStartCompiler create a new compiler container
func BuildAndStartCompiler(name string, cli *client.Client, cconfig *types.CompilerConfig, config *client.ContainerConfig) (cp *Compiler, err error) {
	containerInfo, err := client.SearchContainerByName(cli, "/"+name)
	if err != nil {
		return nil, err
	}

	if containerInfo != nil {
		return StartCompiler(containerInfo, cli, cconfig, config)
	}

	_, err = cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "core-oj/compiler",
		},
		&container.HostConfig{
			PortBindings: nat.PortMap(*config.PortMap),
			Mounts:       []mount.Mount(*config.VolumeMap),
			Resources: container.Resources{
				Memory:     1024 * 1024 * 256,
				MemorySwap: 1024 * 1024 * 256,
				CPUQuota:   10000,
				CPUPeriod:  50000,
			},
		}, nil, name,
	)
	if err != nil {
		return nil, err
	}
	containerInfo, err = client.SearchContainerByName(cli, "/"+name)

	return StartCompiler(containerInfo, cli, cconfig, config)
}

func (cp *Compiler) Close() error {
	cp.conn.Close()
	var timeout = time.Second * 20
	return cp.cli.ContainerStop(context.Background(), cp.container.ID, &timeout)
}

// func StartCompiler() {
// }
