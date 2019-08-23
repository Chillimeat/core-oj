package compiler

import (
	"fmt"
	"testing"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"
)

func TestDaemon(t *testing.T) {
	dae, err := NewDaemon()
	if err != nil {
		t.Error(err)
		return
	}

	go dae.Run()

	fmt.Println("running")

	orzz := make(chan bool)

	dae.ExposeInfo() <- &TaskInfo{
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

	dae.ExposeCompile() <- &TaskCompile{
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
