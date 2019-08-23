package service

import (
	"context"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"
)

const version = "v0.0.4"

// InfoService serve info
type InfoService struct {
	context.Context
	*rpcx.InfoRequest
}

// Serve serve a request
func (srv *InfoService) Serve() (*rpcx.InfoReply, error) {
	return &rpcx.InfoReply{
		CompilerVersion: version,
		CompilerTools:   nil,
	}, nil
}
