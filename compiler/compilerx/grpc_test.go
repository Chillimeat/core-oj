package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"

	"google.golang.org/grpc"
)

const (
	mPort    = ":23369"
	mAddress = "127.0.0.1:23369"
)

func TestCompiler(t *testing.T) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(mAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rpcx.NewCompilerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := c.Info(
		ctx,
		&rpcx.InfoRequest{},
	)
	if err != nil {
		t.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Info: %v\n", r)
}

func TestCompile(t *testing.T) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(mAddress, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rpcx.NewCompilerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r, err := c.Compile(
		ctx, &rpcx.CompileRequest{
			CompilerType: "c++11",
			CodePath:     "/codes/test.cpp",
			AimPath:      "/codes/test",
		},
	)
	if err != nil {
		t.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Compile: %v\n", r)
}
