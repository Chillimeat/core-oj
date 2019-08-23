package compiler

import (
	"context"
	"fmt"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"google.golang.org/grpc"
)

// PortMap helps insert port mapping
type PortMap nat.PortMap

// NewPortMap return a pointer of PortMap
func NewPortMap() (pb *PortMap) {
	pb = new(PortMap)
	*pb = PortMap(make(nat.PortMap))
	return pb
}

type VolumeMap []mount.Mount

func NewVolumeMap() (vp *VolumeMap) {
	vp = new(VolumeMap)
	return vp
}

// Insert a port mapping into the Port Map
func (pb *PortMap) Insert(ip, u, v string) error {

	containerPort, err := nat.NewPort("tcp", v)

	if err != nil {
		return fmt.Errorf("unable to get the port:%v", v)
	}

	(*pb)[containerPort] = []nat.PortBinding{nat.PortBinding{
		HostIP:   ip,
		HostPort: u,
	}}
	return nil
}

func (vp *VolumeMap) InsertBind(source, target string) {
	*vp = append(*vp, mount.Mount{
		Type:   mount.TypeBind,
		Source: source,
		Target: target,
	})
}

// ContainerConfig decides the container's configuration
type ContainerConfig struct {
	PortMap     *PortMap
	VolumeMap   *VolumeMap
	GrpcAddress string
}

// NewContainerConfig return a pointer of ContainerConfig
func NewContainerConfig() *ContainerConfig {
	return &ContainerConfig{
		PortMap:   NewPortMap(),
		VolumeMap: NewVolumeMap(),
	}
}

type Compiler struct {
	container *types.Container
	conn      *grpc.ClientConn
	c         rpcx.CompilerClient
}

func StartCompiler(containerInfo *types.Container, cli *client.Client, config *ContainerConfig) (cp *Compiler, err error) {
	if containerInfo.Status != "running" {
		err = cli.ContainerStart(context.Background(), containerInfo.ID, types.ContainerStartOptions{})
		if err != nil {
			return nil, err
		}
		fmt.Println(err)
	}

	cp = new(Compiler)

	fmt.Printf("Container %s is started\n", containerInfo.ID)
	cp.container = containerInfo
	cp.conn, err = grpc.Dial(config.GrpcAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	cp.c = rpcx.NewCompilerClient(cp.conn)
	return
}

// BuildAndStartCompiler create a new compiler container
func BuildAndStartCompiler(name string, cli *client.Client, config *ContainerConfig) (cp *Compiler, err error) {
	containerInfo, err := client.SearchContainerByName(cli, "/"+name)
	if err != nil {
		return nil, err
	}
	fmt.Println(containerInfo)

	if containerInfo != nil {
		return StartCompiler(containerInfo, cli, config)
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

	return StartCompiler(containerInfo, cli, config)
}

func (cp *Compiler) Close() {
	cp.conn.Close()
}

// func StartCompiler() {
// }
