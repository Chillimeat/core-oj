package service

import (
	"context"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"
	types "github.com/Myriad-Dreamin/core-oj/types"
)

// CompileService serve compile
type CompileService struct {
	ReverseCompilers map[int64]types.Compiler
	context.Context
	*rpcx.CompileRequest
}

// Serve serve a request
func (srv *CompileService) Serve() (*rpcx.CompileReply, error) {

	if cp, ok := srv.ReverseCompilers[srv.CompilerType]; ok {
		err := cp.Compile(srv.CodePath, srv.AimPath)
		if err != nil {
			if v, ok := err.(types.CodeError); ok {
				return &rpcx.CompileReply{
					ResponseCode: rpcx.CompileResponseCode(v.ErrorCode()),
					Info:         v.JudgeError(),
				}, nil
			}
			return nil, err
		}
		return &rpcx.CompileReply{
			ResponseCode: rpcx.CompileResponseCode_Ok,
		}, nil
	}
	return &rpcx.CompileReply{
		ResponseCode: rpcx.CompileResponseCode_CompileError,
		Info:         []byte("compiler not found"),
	}, nil
}
