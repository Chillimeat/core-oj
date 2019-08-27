package dockerx

import (
	"fmt"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
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
	PortMap   *PortMap
	VolumeMap *VolumeMap
	Env       []string
}

// NewContainerConfig return a pointer of ContainerConfig
func NewContainerConfig() *ContainerConfig {
	return &ContainerConfig{
		PortMap:   NewPortMap(),
		VolumeMap: NewVolumeMap(),
	}
}
