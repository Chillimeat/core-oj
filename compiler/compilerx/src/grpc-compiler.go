package server

import (
	"context"
	"fmt"
	"net"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"

	service "github.com/Myriad-Dreamin/core-oj/compiler/compilerx/src/service"

	cpp "github.com/Myriad-Dreamin/core-oj/compiler/compilerx/src/cpp"
	types "github.com/Myriad-Dreamin/core-oj/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is the grpc server for compiling codes
type Server struct {
	Compilers map[string]types.Compiler
}

// NewServer return a pointer of grpc server
func NewServer() (*Server, error) {
	var srv = new(Server)

	srv.Compilers = make(map[string]types.Compiler)
	srv.Compilers["c++11"] = &cpp.Compiler{
		Path: "g++",
		Args: []string{"-std=c++11"},
	}

	return srv, nil
}

func (srv *Server) Compile(ctx context.Context, in *rpcx.CompileRequest) (*rpcx.CompileReply, error) {
	return (&service.CompileService{
		Compilers:      srv.Compilers,
		Context:        ctx,
		CompileRequest: in,
	}).Serve()
}

func (srv *Server) Info(ctx context.Context, in *rpcx.InfoRequest) (*rpcx.InfoReply, error) {
	return (&service.InfoService{
		Context:     ctx,
		InfoRequest: in,
	}).Serve()
}

func (srv *Server) ListenAndServe(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	rpcx.RegisterCompilerServer(s, srv)
	reflection.Register(s)

	fmt.Println("listening on", port)

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil

}
