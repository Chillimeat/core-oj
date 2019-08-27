package compiler

import (
	"context"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"
	types "github.com/Myriad-Dreamin/core-oj/types"
)

type Daemon struct {
	cli          *client.Client
	compilerPool chan *Compiler

	// TaskCompileChan chan *TaskCompile
	// TaskInfoChan    chan *TaskInfo
}

func NewDaemon(cli *client.Client) (dae *Daemon, err error) {
	dae = new(Daemon)
	dae.compilerPool = make(chan *Compiler, 10)
	dae.cli = cli

	config := client.NewContainerConfig()
	cconfig := &types.CompilerConfig{GrpcAddress: "127.0.0.1:23367"}
	config.PortMap.Insert("127.0.0.1", "23367", "23366")
	config.VolumeMap.InsertBind("/home/kamiyoru/data/test", "/codes")
	config.VolumeMap.InsertBind("/home/kamiyoru/data/compiler_tools", "/compiler_tools")
	cp, err := BuildAndStartCompiler("compiler2", dae.cli, cconfig, config)
	if err != nil {
		return nil, err
	}
	dae.compilerPool <- cp
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

func (dae *Daemon) Close() {
	worker, ok := <-dae.compilerPool
	for ok {
		worker.Close()
		worker, ok = <-dae.compilerPool
	}
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
