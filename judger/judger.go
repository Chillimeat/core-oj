package judger

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"

	types "github.com/Myriad-Dreamin/core-oj/types"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

type Judger struct {
	container *dockertypes.Container
	conn      *SocketY
}

const (
	jtp = "/home/kamiyoru/data/judger_tools/socks/"
)

type SocketY struct {
	*net.UnixConn
	lenBuffer  int32
	buffer     []byte
	precBuffer *bytes.Buffer
	suffBuffer *bytes.Buffer
}

func NewSocketY(conn *net.UnixConn) (sy *SocketY) {
	sy = new(SocketY)
	sy.UnixConn = conn
	sy.buffer = make([]byte, 3333)
	sy.precBuffer, sy.suffBuffer = bytes.NewBuffer(sy.buffer[0:4]), bytes.NewBuffer(sy.buffer[4:])
	return
}

func (sy *SocketY) Receive() ([]byte, error) {
	err := binary.Read(sy.UnixConn, binary.BigEndian, &sy.lenBuffer)
	if err != nil {
		return nil, err
	}
	_, err = io.ReadFull(sy.UnixConn, sy.buffer[0:sy.lenBuffer])
	if err != nil {
		return nil, err
	}
	return sy.buffer[0:sy.lenBuffer], nil
}

func (sy *SocketY) Send(b []byte) error {
	sy.precBuffer.Reset()
	sy.suffBuffer.Reset()
	binary.Write(sy.precBuffer, binary.BigEndian, int32(len(b)))
	binary.Write(sy.suffBuffer, binary.BigEndian, b)
	_, err := sy.UnixConn.Write(sy.buffer[0 : 4+len(b)])
	if err != nil {
		return err
	}
	return nil
}

func (js *Judger) Judge(problem *morm.Problem) {

	// b, err := json.Marshal(testCase)
	// js.conn.Send(b)
	// b, err = js.conn.Receive()
	// fmt.Println(string(b), err)
}

func StartJudger(containerInfo *dockertypes.Container, cli *client.Client, cconfig *types.JudgerConfig, config *client.ContainerConfig) (cp *Judger, err error) {
	if containerInfo.Status != "running" {
		err = cli.ContainerStart(context.Background(), containerInfo.ID, dockertypes.ContainerStartOptions{})
		if err != nil {
			return nil, err
		}
		fmt.Println(err)
	}

	cp = new(Judger)
	fmt.Printf("Container %s is started\n", containerInfo.ID)
	cp.container = containerInfo

	uaddr, err := net.ResolveUnixAddr("unix", cconfig.UnixAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUnix("unix", nil, uaddr)
	if err != nil {
		return nil, err
	}
	cp.conn = NewSocketY(conn)
	return
}

// BuildAndStartJudger create a new Judger container
func BuildAndStartJudger(name string, cli *client.Client, cconfig *types.JudgerConfig, config *client.ContainerConfig) (cp *Judger, err error) {
	containerInfo, err := client.SearchContainerByName(cli, "/"+name)
	if err != nil {
		return nil, err
	}
	fmt.Println(containerInfo)

	if containerInfo != nil {
		return StartJudger(containerInfo, cli, cconfig, config)
	}

	_, err = cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "core-oj/judger",
			Env:   config.Env,
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

	return StartJudger(containerInfo, cli, cconfig, config)
}

func (cp *Judger) Close() {
	cp.conn.Close()
}

// func StartJudger() {
// }

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"os/exec"
// 	"strconv"

// 	client "github.com/Myriad-Dreamin/core-oj/docker-client"

// 	"github.com/docker/docker/api/dockertypes"
// 	"github.com/docker/docker/api/dockertypes/mount"
// )

// type Judger struct {
// 	container *dockertypes.Container
// 	cli       *client.Client
// }

// func runJudger(containerInfo *dockertypes.Container, cli *client.Client, config *client.ContainerConfig) (cp *Judger, err error) {
// 	if containerInfo.Status != "running" {
// 		err = cli.ContainerStart(context.Background(), containerInfo.ID, dockertypes.ContainerStartOptions{})
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(err)
// 	}

// 	cp = new(Judger)
// 	cp.container = containerInfo
// 	fmt.Printf("Container %s is started\n", cp.container.ID)
// 	cp.cli = cli
// 	return
// }

// // RunJudger create a new Judger container
// func RunJudger(name string, cli *client.Client, config *client.ContainerConfig) (cp *Judger, err error) {
// 	containerInfo, err := client.SearchContainerByName(cli, "/"+name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(containerInfo)

// 	if containerInfo != nil {
// 		return runJudger(containerInfo, cli, config)
// 	}

// 	var args = []string{
// 		"run", "-id",
// 		"--name", name,
// 	}

// 	// for k, v := range config.PortMap {
// 	// 	args = append(args, []string{"-p", (v.HostIP + v.HostPort) + ":" + k, }...)
// 	// }

// 	for _, v := range []mount.Mount(*config.VolumeMap) {
// 		args = append(args, []string{"-v", v.Source + ":" + v.Target}...)

// 	}

// 	args = append(args, []string{
// 		"--memory", strconv.Itoa(1024 * 1024 * 630),
// 		"--memory-swap", strconv.Itoa(1024 * 1024 * 630),
// 		"--cpu-quota", strconv.Itoa(10000),
// 		"--cpu-period", strconv.Itoa(50000),
// 	}...)

// 	args = append(args, "core-oj/judger")

// 	cmd := exec.CommandContext(
// 		context.TODO(),
// 		"docker",
// 		args...,
// 	)

// 	var output = new(bytes.Buffer)
// 	cmd.Stdout = output
// 	var errput = new(bytes.Buffer)
// 	cmd.Stderr = errput

// 	if err := cmd.Run(); err != nil {
// 		fmt.Println(err, string(output.String()), string(errput.String()))
// 		return nil, err
// 	}
// 	fmt.Println("here")
// 	fmt.Println(err, string(output.String()), string(errput.String()))

// 	containerInfo, err = client.SearchContainerByName(cli, "/"+name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Printf("Container %s is running\n", containerInfo.ID)

// 	return runJudger(containerInfo, cli, config)
// }

// func (cp *Judger) Close() {
// }

// // func runJudger() {
// // }
