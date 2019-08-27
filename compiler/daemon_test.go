package compiler

import (
	"context"
	"fmt"
	"testing"
	"time"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"
	client "github.com/Myriad-Dreamin/core-oj/docker-client"
)

type TaskCompile struct {
	Ctx      context.Context
	CallBack func(*rpcx.CompileReply)
	Cerr     func(error)
	Req      *rpcx.CompileRequest
}

type TaskInfo struct {
	Ctx      context.Context
	CallBack func(*rpcx.InfoReply)
	Cerr     func(error)
	Req      *rpcx.InfoRequest
}

func TestDaemon(t *testing.T) {
	cli, err := client.Connect("unix:///var/run/docker.sock", "v1.40")
	if err != nil {
		t.Error(err)
		return
	}

	dae, err := NewDaemon(cli)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	var TaskCompileChan = make(chan *TaskCompile)
	var TaskInfoChan = make(chan *TaskInfo)

	go dae.Run(context.Background(), func(worker *Compiler) {
		fmt.Println("worker", worker)
		select {
		case task := <-TaskCompileChan:
			fmt.Printf("tasking %v %T\n", task, task)
			ctx, cancel := context.WithTimeout(task.Ctx, time.Second*10)
			ret, err := worker.Compile(ctx, task.Req)
			cancel()
			if err != nil {
				task.Cerr(err)
				return
			}
			task.CallBack(ret)
		case task := <-TaskInfoChan:
			fmt.Printf("tasking %v %T\n", task, task)
			ctx, cancel := context.WithTimeout(task.Ctx, time.Second*10)
			ret, err := worker.Info(ctx, task.Req)
			cancel()
			if err != nil {
				task.Cerr(err)
				return
			}
			task.CallBack(ret)
		case <-ctx.Done():
			return
		}
	})

	fmt.Println("running")

	orzz := make(chan bool)

	TaskInfoChan <- &TaskInfo{
		Ctx: context.Background(),
		CallBack: func(cr *rpcx.InfoReply) {
			fmt.Println("infoing...", cr.CompilerVersion)
			orzz <- true
		},
		Cerr: func(err error) {
			t.Error(err)
			orzz <- true
			return
		},
		Req: &rpcx.InfoRequest{},
	}
	<-orzz
	fmt.Println("orz")

	TaskCompileChan <- &TaskCompile{
		Ctx: context.Background(),
		CallBack: func(cr *rpcx.CompileReply) {
			fmt.Println("compiled...", cr.ResponseCode, string(cr.Info))
			orzz <- true
		},
		Cerr: func(err error) {
			t.Error(err)
			orzz <- true
			return
		},
		Req: &rpcx.CompileRequest{
			CompilerType: "c++11",
			CodePath:     "/codes/test.cpp",
			AimPath:      "/codes/test",
		},
	}
	<-orzz
}
