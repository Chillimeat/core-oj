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

	dae.Expose() <- &TaskCompile{
		CallBack: func(cr *rpcx.CompileReply) {
			fmt.Println(cr)
		},
		Cerr: func(err error) {
			t.Error(err)
			return
		},
		Req: &rpcx.CompileRequest{
			CompilerType: "cpp",
			CodePath:     []byte{},
		},
	}

}
