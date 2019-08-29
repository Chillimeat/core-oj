package compiler

import (
	"bytes"
	"context"
	"strconv"
	"sync"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"
	types "github.com/Myriad-Dreamin/core-oj/types"
)

type Daemon struct {
	cli          *client.Client
	compilerPool chan *Compiler

	// TaskCompileChan chan *TaskCompile
	// TaskInfoChan    chan *TaskInfo
}

type DaemonConfig struct {
	Number            int
	CompilerName      string
	Host              string
	PortDst           string
	StartPort         string
	CodesPath         string
	ComplierToolsPath string
}

func NewDaemon(cli *client.Client, config *DaemonConfig) (dae *Daemon, err error) {
	dae = new(Daemon)
	dae.compilerPool = make(chan *Compiler, config.Number)
	dae.cli = cli

	cfg := client.NewContainerConfig()

	// cfg.VolumeMap.InsertBind("/home/kamiyoru/data/test", "/codes")
	// cfg.VolumeMap.InsertBind("/home/kamiyoru/data/compiler_tools", "/compiler_tools")
	cfg.VolumeMap.InsertBind(config.CodesPath, "/codes")
	cfg.VolumeMap.InsertBind(config.ComplierToolsPath, "/compiler_tools")

	var wg sync.WaitGroup
	var errs Errors
	for idx := 0; idx < config.Number; idx++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			cconfig := &types.CompilerConfig{GrpcAddress: config.Host + ":" + strconv.Itoa(23367+idx)}
			ccfg := client.NewContainerConfig()
			ccfg.PortMap.Insert(config.Host, strconv.Itoa(23367+idx), config.StartPort)
			ccfg.VolumeMap = cfg.VolumeMap
			cp, err := BuildAndStartCompiler(config.CompilerName+strconv.Itoa(idx), dae.cli, cconfig, ccfg)
			if err != nil {
				errs = append(errs, err)
				return
			}

			dae.compilerPool <- cp
		}(idx)
	}
	wg.Wait()

	if len(errs) != 0 {
		dae = nil
		err = errs
		return
	}

	return
}

// func (dae *Daemon) ExposeCompile() chan<- *TaskCompile {
// 	return dae.TaskCompileChan
// }

// func (dae *Daemon) ExposeInfo() chan<- *TaskInfo {
// 	return dae.TaskInfoChan
// }

// func (dae *Daemon) Run(ctx context.Context) error {
// 	for {
// 		fmt.Println("QwQ")
// 		select {
// 		case task := <-dae.TaskCompileChan:
// 			fmt.Printf("tasking %v %T\n", task, task)
// 			if _, ok := <-task.Ctx.Done(); ok {
// 				continue
// 			}
// 			go func() {
// 				select {
// 				case worker := <-dae.compilerPool:

// 					fmt.Println("worker", worker)
// 					ctx, cancel := context.WithTimeout(task.Ctx, time.Second*10)
// 					ret, err := worker.c.Compile(ctx, task.Req)
// 					cancel()
// 					dae.compilerPool <- worker
// 					if err != nil {
// 						task.Cerr(err)
// 						return
// 					}
// 					task.CallBack(ret)
// 				case <-task.Ctx.Done():
// 					return
// 				}
// 			}()
// 		case task := <-dae.TaskInfoChan:
// 			fmt.Printf("tasking %v %T\n", task, task)
// 			if _, ok := <-task.Ctx.Done(); ok {
// 				continue
// 			}
// 			go func() {
// 				select {
// 				case worker := <-dae.compilerPool:

// 					fmt.Println("worker", worker)
// 					ctx, cancel := context.WithTimeout(task.Ctx, time.Second*10)
// 					ret, err := worker.c.Info(ctx, task.Req)
// 					cancel()
// 					dae.compilerPool <- worker
// 					if err != nil {
// 						task.Cerr(err)
// 						return
// 					}
// 					task.CallBack(ret)
// 				case <-task.Ctx.Done():
// 					return
// 				}
// 			}()
// 		case <-ctx.Done():
// 			return nil
// 		}
// 	}
// }

type Errors []error

func (errs Errors) Error() string {
	var b = new(bytes.Buffer)
	for _, err := range errs {
		b.WriteString(err.Error())
	}
	return b.String()
}

func (dae *Daemon) Close() error {
	worker, ok := <-dae.compilerPool
	var errs Errors
	for ok {
		err := worker.Close()
		if err != nil {
			errs = append(errs, err)
		}
		worker, ok = <-dae.compilerPool
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (dae *Daemon) Run(ctx context.Context, withWorker func(*Compiler)) {
	for {
		select {
		case worker := <-dae.compilerPool:
			go func() {
				withWorker(worker)
				dae.compilerPool <- worker
			}()
		case <-ctx.Done():
			dae.Close()
			return
		}
	}
}
