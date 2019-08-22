package service

import (
	"context"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"
)

// CompileService serve compile
type CompileService struct {
	context.Context
	*rpcx.CompileRequest
}

// Serve serve a request
func (srv *CompileService) Serve() (*rpcx.CompileReply, error) {
	return &rpcx.CompileReply{
		ResponseCode: rpcx.CompileResponseCode_Ok,
	}, nil
}
