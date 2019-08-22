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
	mPort    = ":23367"
	mAddress = "localhost:23367"
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
