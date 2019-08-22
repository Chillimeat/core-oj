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

type Daemon struct {
	cli          *client.Client
	compilerPool chan *Compiler

	TaskCompileChan chan *TaskCompile
}

func NewDaemon() (dae *Daemon, err error) {
	dae = new(Daemon)
	dae.compilerPool = make(chan *Compiler, 1)

	dae.cli, err = client.Connect("unix:///var/run/docker.sock", "v1.40")
	if err != nil {
		return nil, err
	}

	config := NewContainerConfig()
	config.PortMap.Insert("127.0.0.1", "23366", "23367")
	config.GrpcAddress = "127.0.0.1:23366"
	cp, err := BuildAndStartCompiler("compiler2", dae.cli, config)
	if err != nil {
		return nil, err
	}
	dae.compilerPool <- cp
	return
}

func (dae *Daemon) Expose() chan<- *TaskCompile {
	return dae.TaskCompileChan
}

func (dae *Daemon) Run() error {
	for {
		select {
		case task := <-dae.TaskCompileChan:
			fmt.Println("tasking", task)
			go func() {
				select {
				case worker := <-dae.compilerPool:

					fmt.Println("worker", worker)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
					ret, err := worker.c.Compile(ctx, task.Req)
					cancel()
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
