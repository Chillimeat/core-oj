package judger

import (
	"context"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"
	types "github.com/Myriad-Dreamin/core-oj/types"
)

type Daemon struct {
	cli        *client.Client
	judgerPool chan *Judger

	// TaskJudgerhan chan *TaskJudger	// TaskInfoChan    chan *TaskInfo
}

func NewDaemon(cli *client.Client) (dae *Daemon, err error) {
	dae = new(Daemon)
	dae.judgerPool = make(chan *Judger, 10)
	dae.cli = cli

	config := client.NewContainerConfig()
	cconfig := &types.JudgerConfig{UnixAddress: "/home/kamiyoru/data/judger_tools/socks/judger0.sock"}
	config.VolumeMap.InsertBind("home/kamiyoru/data/test", "/codes")
	config.VolumeMap.InsertBind("home/kamiyoru/data/judger_tools", "/judger_tools")
	cp, err := BuildAndStartJudger("judger0", dae.cli, cconfig, config)
	if err != nil {
		return nil, err
	}
	dae.judgerPool <- cp
	return
}

// func (dae *Daemon) ExposeJudger) chan<- *TaskJudger{
// 	return dae.TaskJudgerhan
// }

// func (dae *Daemon) ExposeInfo() chan<- *TaskInfo {
// 	return dae.TaskInfoChan
// }

// func (dae *Daemon) Run(ctx context.Context) error {
// 	for {
// 		fmt.Println("QwQ")
// 		select {
// 		case task := <-dae.TaskJudgerhan:
// 			fmt.Printf("tasking %v %T\n", task, task)
// 			if _, ok := <-task.Ctx.Done(); ok {
// 				continue
// 			}
// 			go func() {
// 				select {
// 				case worker := <-dae.judgerPool:

// 					fmt.Println("worker", worker)
// 					ctx, cancel := context.WithTimeout(task.Ctx, time.Second*10)
// 					ret, err := worker.c.Judgerctx, task.Req)
// 					cancel()
// 					dae.judgerPool <- worker
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
// 				case worker := <-dae.judgerPool:

// 					fmt.Println("worker", worker)
// 					ctx, cancel := context.WithTimeout(task.Ctx, time.Second*10)
// 					ret, err := worker.c.Info(ctx, task.Req)
// 					cancel()
// 					dae.judgerPool <- worker
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
	worker, ok := <-dae.judgerPool
	for ok {
		worker.Close()
		worker, ok = <-dae.judgerPool
	}
}

func (dae *Daemon) Run(ctx context.Context, withWorker func(*Judger)) {
	for {
		select {
		case worker := <-dae.judgerPool:
			go func() {
				withWorker(worker)
				dae.judgerPool <- worker
			}()
		case <-ctx.Done():
			dae.Close()
			return
		}
	}
}
