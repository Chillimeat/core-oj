package compiler

import (
	"context"
	"errors"
	"fmt"
	"time"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"
)

type TaskCompile struct {
	CallBack func(*rpcx.CompileReply)
	Cerr     func(error)
	Req      *rpcx.CompileRequest
}

type TaskInfo struct {
	CallBack func(*rpcx.InfoReply)
	Cerr     func(error)
	Req      *rpcx.InfoRequest
}

type Daemon struct {
	cli          *client.Client
	compilerPool chan *Compiler

	TaskCompileChan chan *TaskCompile
	TaskInfoChan    chan *TaskInfo
}

func NewDaemon() (dae *Daemon, err error) {
	dae = new(Daemon)
	dae.compilerPool = make(chan *Compiler, 1)
	dae.TaskCompileChan = make(chan *TaskCompile)
	dae.TaskInfoChan = make(chan *TaskInfo)

	dae.cli, err = client.Connect("unix:///var/run/docker.sock", "v1.40")
	if err != nil {
		return nil, err
	}

	config := NewContainerConfig()
	config.PortMap.Insert("127.0.0.1", "23367", "23366")
	config.VolumeMap.InsertBind("/home/kamiyoru/data/test", "/codes")
	config.VolumeMap.InsertBind("/home/kamiyoru/data/compiler_tools", "/compiler_tools")
	config.GrpcAddress = "127.0.0.1:23367"
	cp, err := BuildAndStartCompiler("compiler2", dae.cli, config)
	if err != nil {
		return nil, err
	}
	dae.compilerPool <- cp
	return
}

func (dae *Daemon) ExposeCompile() chan<- *TaskCompile {
	return dae.TaskCompileChan
}

func (dae *Daemon) ExposeInfo() chan<- *TaskInfo {
	return dae.TaskInfoChan
}

func (dae *Daemon) Run() error {
	for {
		fmt.Println("QwQ")
		select {
		case task := <-dae.TaskCompileChan:
			fmt.Printf("tasking %v %T\n", task, task)
			go func() {
				select {
				case worker := <-dae.compilerPool:

					fmt.Println("worker", worker)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					ret, err := worker.c.Compile(ctx, task.Req)
					cancel()
					dae.compilerPool <- worker
					if err != nil {
						task.Cerr(err)
						return
					}
					task.CallBack(ret)
				case <-time.After(time.Second * 10):
					task.Cerr(errors.New("submit timeout"))
				}
			}()
		case task := <-dae.TaskInfoChan:
			fmt.Printf("tasking %v %T\n", task, task)
			go func() {
				select {
				case worker := <-dae.compilerPool:

					fmt.Println("worker", worker)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					ret, err := worker.c.Info(ctx, task.Req)
					cancel()
					dae.compilerPool <- worker
					if err != nil {
						task.Cerr(err)
						return
					}
					task.CallBack(ret)
				case <-time.After(time.Second * 10):
					task.Cerr(errors.New("submit timeout"))
				}
			}()
		}
	}
}
