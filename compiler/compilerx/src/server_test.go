package server

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"

	"google.golang.org/grpc"
)

const (
	mPort    = ":23366"
	mAddress = "localhost:23366"
)

func TestBuildCompiler(t *testing.T) {
	var srv *Server
	var err error
	if srv, err = NewServer(); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := srv.ListenAndServe(mPort); err != nil {
			log.Fatal(err)
		}
	}()

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
