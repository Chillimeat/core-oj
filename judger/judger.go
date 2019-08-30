package judger

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"

	config "github.com/Myriad-Dreamin/core-oj/config"
	problemconfig "github.com/Myriad-Dreamin/core-oj/problem-config"
	types "github.com/Myriad-Dreamin/core-oj/types"
	kvorm "github.com/Myriad-Dreamin/core-oj/types/kvorm"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

type Judger struct {
	cli       *client.Client
	container *dockertypes.Container
	conn      *SocketY
}

const (
	jtp = "/home/kamiyoru/data/judger_tools/socks/"
)

type SocketY struct {
	*net.UnixConn
	addr       *net.UnixAddr
	lastDial   time.Time
	lenBuffer  int32
	buffer     []byte
	precBuffer *bytes.Buffer
	suffBuffer *bytes.Buffer
}

func NewSocketY(addr *net.UnixAddr) (sy *SocketY, err error) {
	time.Sleep(time.Second * 2)
	sy = new(SocketY)
	sy.addr = addr
	sy.buffer = make([]byte, 3333)
	sy.precBuffer, sy.suffBuffer = bytes.NewBuffer(sy.buffer[0:4]), bytes.NewBuffer(sy.buffer[4:])

	sy.UnixConn, err = net.DialUnix("unix", nil, addr)
	if err != nil {
		return nil, err
	}
	sy.lastDial = time.Now()
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
	sy.lastDial = time.Now()
	return sy.buffer[0:sy.lenBuffer], nil
}

func (sy *SocketY) Send(b []byte) (err error) {
	if time.Now().Sub(sy.lastDial) > 400*time.Millisecond {
		sy.UnixConn, err = net.DialUnix("unix", nil, sy.addr)
		if err != nil {
			return err
		}
	}
	sy.precBuffer.Reset()
	sy.suffBuffer.Reset()
	binary.Write(sy.precBuffer, binary.BigEndian, int32(len(b)))
	binary.Write(sy.suffBuffer, binary.BigEndian, b)
	_, err = sy.UnixConn.Write(sy.buffer[0 : 4+len(b)])
	if err != nil {
		return err
	}

	return nil
}

func calcuateSPJ(lt, fp string) string {
	return "todo"
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func maxTime(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (js *Judger) Judge(code *morm.Code, problem *morm.Problem) ([]*kvorm.ProcState, error) {
	var path = config.Config().ProblemPath + strconv.Itoa(problem.ID)
	var inpath = "/problems/" + strconv.Itoa(problem.ID)
	var outpath = inpath
	var cfg = new(problemconfig.ProblemConfig)

	err := problemconfig.Load(cfg, path+"/problem-config")

	if err != nil {
		return nil, err
	}

	var testCase = new(types.TestCase)
	testCase.TestType = code.CodeType
	testCase.TestWorkDir = "/codes/" + hex.EncodeToString(code.Hash)
	testCase.OptionStream = 0
	if cfg.SpecialJudgeConfig.SpecialJudge {
		testCase.SpecialJudge = true
		testCase.SpecialJudgePath = calcuateSPJ(
			cfg.SpecialJudgeConfig.LanguageType,
			cfg.SpecialJudgeConfig.FilePath,
		)
	}

	switch cfg.JudgeConfig.Type {
	case "acm":
		if len(cfg.JudgeConfig.Tasks) != 1 {
			return nil, errors.New("problem cfg error?(judge-config.tasks")
		}
		var retstat = make([]*kvorm.ProcState, 0, len(cfg.JudgeConfig.Tasks))
		var task = cfg.JudgeConfig.Tasks[0]
		var intaskpath = inpath + task.InputPath + "in"
		var outtaskPath = outpath + task.OutputPath + "out"
		testCase.TimeLimit = task.TimeLimit
		testCase.MemoryLimit = task.MemoryLimit
		var totstat = new(kvorm.ProcState)
		for testCase.CaseNumber = 1; testCase.CaseNumber <= task.CaseCount; testCase.CaseNumber++ {
			testCase.InputPath = intaskpath + strconv.Itoa(testCase.CaseNumber) + ".txt"
			testCase.StdOutputPath = outtaskPath + strconv.Itoa(testCase.CaseNumber) + ".txt"

			b, err := json.MarshalIndent(testCase, "", "    ")
			if err != nil {
				return nil, err
			}

			b, err = json.Marshal(testCase)
			if err != nil {
				return nil, err
			}
			err = js.conn.Send(b)
			if err != nil {
				return nil, err
			}
			b, err = js.conn.Receive()
			if err != nil {
				return nil, err
			}

			var stat kvorm.ProcState
			err = json.Unmarshal(b, &stat)
			if err != nil {
				return nil, err
			}
			totstat.CodeError = stat.CodeError
			totstat.TimeUsed = maxTime(totstat.TimeUsed, stat.TimeUsed)
			totstat.MemoryUsed = maxFloat64(totstat.MemoryUsed, stat.MemoryUsed)

			if totstat.CodeError.ErrorCode() != types.StatusAccepted {
				return append(retstat, totstat), nil
			}
		}
		return append(retstat, totstat), nil
	default:
		return nil, errors.New("problem config error?(judge-config.judge-type")
	}

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
	}

	cp = new(Judger)
	cp.cli = cli
	cp.container = containerInfo

	uaddr, err := net.ResolveUnixAddr("unix", cconfig.UnixAddress)
	if err != nil {
		return nil, err
	}

	cp.conn, err = NewSocketY(uaddr)
	if err != nil {
		return nil, err
	}
	return
}

// BuildAndStartJudger create a new Judger container
func BuildAndStartJudger(name string, cli *client.Client, cconfig *types.JudgerConfig, config *client.ContainerConfig) (cp *Judger, err error) {
	containerInfo, err := client.SearchContainerByName(cli, "/"+name)
	if err != nil {
		return nil, err
	}

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

func (cp *Judger) Close() error {
	cp.conn.Close()
	var timeout = time.Second * 20
	return cp.cli.ContainerStop(context.Background(), cp.container.ID, &timeout)
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
