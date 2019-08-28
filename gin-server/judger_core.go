package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
	"unsafe"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"
	"github.com/Myriad-Dreamin/core-oj/log"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"

	compiler "github.com/Myriad-Dreamin/core-oj/compiler"
	judger "github.com/Myriad-Dreamin/core-oj/judger"

	types "github.com/Myriad-Dreamin/core-oj/types"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"
)

// JudgeService defines handler functions of code router
type JudgeService struct {
	cdae   *compiler.Daemon
	jdae   *judger.Daemon
	cr     *morm.Coder
	pr     *morm.Problemer
	logger log.TendermintLogger
}

// NewJudgeService return a pointer of JudgeService
func NewJudgeService(coder *morm.Coder, problemer *morm.Problemer, logger log.TendermintLogger) (*JudgeService, error) {
	cli, err := client.Connect("unix:///var/run/docker.sock", "v1.40")
	if err != nil {
		return nil, err
	}

	cdae, err := compiler.NewDaemon(cli)
	if err != nil {
		return nil, err
	}

	jdae, err := judger.NewDaemon(cli)
	if err != nil {
		return nil, err
	}

	return &JudgeService{
		cdae:   cdae,
		jdae:   jdae,
		cr:     coder,
		pr:     problemer,
		logger: logger,
	}, nil
}

// ProcessAllCodes just run all daemon
func (js *JudgeService) ProcessAllCodes(ctx context.Context) {

	var (
		waitingCodes = js.cr.ExposeWaitingCodes()
		runningCodes = js.cr.ExposeRunningCodes()
	)
	go js.cdae.Run(ctx, func(worker *compiler.Compiler) {
		var req *rpcx.CompileRequest
		var code *morm.Code
		for req == nil {
			code = <-waitingCodes
			switch code.CodeType {
			case morm.CodeTypeCpp11:
				req = &rpcx.CompileRequest{
					CompilerType: "c++11",
					CodePath:     "/codes/" + hex.EncodeToString(code.Hash) + "/main.cpp",
					AimPath:      "/codes/" + hex.EncodeToString(code.Hash) + "main",
				}
			default:
				continue
			}
		}
		code.Status = types.StatusCompiling
		js.cr.StartToExecuteTask(code)
		jsctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		ret, err := worker.Compile(jsctx, req)

		// todo: then
		cancel()
		if err != nil {
			if err, ok := err.(types.CodeError); ok {
				// Warning: unsafeConvert
				atomic.StoreInt64((*int64)(unsafe.Pointer(&code.Status)), int64(err.ErrorCode()))
				js.logger.Debug("catch code error", "error", err)
			} else {
				// Warning: unsafeConvert
				atomic.StoreInt64((*int64)(unsafe.Pointer(&code.Status)), int64(types.StatusUnknownError))
				js.logger.Debug("catch unknown code error", "error", err)
			}
			return
		}
		fmt.Println("compiled...", ret.ResponseCode, string(ret.Info))
		runningCodes <- code
	})
	for {
		select {
		case <-ctx.Done():
			return
		case code := <-runningCodes:
			fmt.Println("running", code)
		}
	}
}
