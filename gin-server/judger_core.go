package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	rpcx "github.com/Myriad-Dreamin/core-oj/compiler/grpc"
	helpfunc "github.com/Myriad-Dreamin/core-oj/help-func"
	"github.com/Myriad-Dreamin/core-oj/log"
	kvorm "github.com/Myriad-Dreamin/core-oj/types/kvorm"
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
	psr    *kvorm.ProcStater
	logger log.TendermintLogger
}

// NewJudgeService return a pointer of JudgeService
func NewJudgeService(coder *morm.Coder, problemer *morm.Problemer, procStater *kvorm.ProcStater, logger log.TendermintLogger) (*JudgeService, error) {
	cli, err := client.Connect("unix:///var/run/docker.sock", "v1.40")
	if err != nil {
		return nil, err
	}
	var cdae *compiler.Daemon
	var jdae *judger.Daemon
	if err := helpfunc.Walk(2, func(idx int) error {
		if idx == 0 {
			var err error
			cdae, err = compiler.NewDaemon(cli, &compiler.DaemonConfig{
				Number:            8,
				CompilerName:      "compiler",
				Host:              "127.0.0.1",
				PortDst:           "23366",
				StartPort:         "23366",
				CodesPath:         "/home/kamiyoru/data/test",
				ComplierToolsPath: "/home/kamiyoru/data/compiler_tools",
			})
			if err != nil {
				return err
			}
		} else {
			var err error
			jdae, err = judger.NewDaemon(cli, &judger.DaemonConfig{
				Number:           20,
				JudgerName:       "judger",
				ProblemsPath:     "/home/kamiyoru/data/problems",
				CodesPath:        "/home/kamiyoru/data/test",
				CheckerToolsPath: "/home/kamiyoru/data/checker_tools",
				JudgerToolsPath:  "/home/kamiyoru/data/judger_tools",
			})
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		logger.Error("start error", "error", err)
		return nil, err
	}

	return &JudgeService{
		cdae:   cdae,
		jdae:   jdae,
		cr:     coder,
		pr:     problemer,
		psr:    procStater,
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
			select {
			case <-ctx.Done():
				return
			case code = <-waitingCodes:
				req = &rpcx.CompileRequest{
					CompilerType: code.CodeType,
					CodePath:     "/codes/" + hex.EncodeToString(code.Hash),
					AimPath:      "/codes/" + hex.EncodeToString(code.Hash),
				}
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
				js.logger.Debug("catch codClosee error", "error", err)
			} else {
				// Warning: unsafeConvert
				atomic.StoreInt64((*int64)(unsafe.Pointer(&code.Status)), int64(types.StatusUnknownError))
				js.logger.Debug("catch unknown code error", "error", err)
			}
			settled, err := js.cr.SettleTask(code.ID)
			if !settled || err != nil {
				js.logger.Debug("catch settle task error", "settled", settled, "error", err)
				return
			}
			return
		}
		js.logger.Info("compiled...", "response", ret.ResponseCode.String(), "info", string(ret.Info))
		if ret.ResponseCode != 0 {
			settled, err := js.cr.SettleTask(code.ID)
			if !settled || err != nil {
				js.logger.Debug("catch settle task error", "settled", settled, "error", err)
				return
			}
		} else {
			runningCodes <- code
		}
	})

	go js.jdae.Run(ctx, func(jg *judger.Judger) {
		select {
		case <-ctx.Done():
			return
		case code := <-runningCodes:
			code.Status = types.StatusRunning
			problem, err := js.pr.Query(code.ProblemID)
			if err != nil {
				js.logger.Debug("catch query problem error", "error", err)
				return
			}
			fmt.Println(problem)
			results, err := jg.Judge(code, problem)
			if err != nil {
				js.logger.Debug("catch judge problem error", "error", err)
				return
			}

			err = js.psr.InsertP(code.ID, results)
			if err != nil {
				js.logger.Debug("catch judge problem error", "error", err)
				return
			}

			// if len(results) == 0 {
			// 	js.logger.Debug("catch judge problem error", "error", "len(results)==0")
			// 	return
			// }
			status := types.StatusAccepted
			for _, result := range results {
				if result.CodeError.ErrorCode() != types.StatusAccepted {
					status = result.CodeError.ErrorCode()
					break
				}
			}
			code.Status = status
			settled, err := js.cr.SettleTask(code.ID)
			if !settled || err != nil {
				js.logger.Debug("catch settle task error", "settled", settled, "error", err)
				return
			}
		}
	})
}

func (js *JudgeService) Close() error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := js.jdae.Close(); err != nil {
			js.logger.Debug("judger daemon close error", err)
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		if err := js.cdae.Close(); err != nil {
			js.logger.Debug("compiler daemon close error", err)
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}
