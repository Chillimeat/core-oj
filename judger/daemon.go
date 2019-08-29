package judger

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"sync"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"
	types "github.com/Myriad-Dreamin/core-oj/types"
)

type Daemon struct {
	cli        *client.Client
	judgerPool chan *Judger

	// TaskJudgerhan chan *TaskJudger	// TaskInfoChan    chan *TaskInfo
}

type DaemonConfig struct {
	Number           int
	JudgerName       string
	ProblemsPath     string
	CodesPath        string
	CheckerToolsPath string
	JudgerToolsPath  string
	Env              []string
}

func NewDaemon(cli *client.Client, config *DaemonConfig) (dae *Daemon, err error) {
	dae = new(Daemon)
	dae.judgerPool = make(chan *Judger, config.Number)
	dae.cli = cli

	cfg := client.NewContainerConfig()
	// cfg.VolumeMap.InsertBind("/home/kamiyoru/data/test", "/codes")
	// cfg.VolumeMap.InsertBind("/home/kamiyoru/data/judger_tools", "/judger_tools")
	// cfg.VolumeMap.InsertBind("/home/kamiyoru/data/problems", "/problems")
	// cfg.VolumeMap.InsertBind("/home/kamiyoru/data/checker_tools", "/checker_tools")
	cfg.VolumeMap.InsertBind(config.CodesPath, "/codes")
	cfg.VolumeMap.InsertBind(config.JudgerToolsPath, "/judger_tools")
	cfg.VolumeMap.InsertBind(config.ProblemsPath, "/problems")
	cfg.VolumeMap.InsertBind(config.CheckerToolsPath, "/checker_tools")
	cfg.Env = config.Env
	var wg = new(sync.WaitGroup)
	var errs Errors
	for idx := 0; idx < config.Number; idx++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			cconfig := &types.JudgerConfig{UnixAddress: "/home/kamiyoru/data/judger_tools/socks/" + config.JudgerName + strconv.Itoa(idx) + ".sock"}
			ccfg := client.NewContainerConfig()
			ccfg.VolumeMap = cfg.VolumeMap
			ccfg.Env = append(cfg.Env, "NAME="+config.JudgerName+strconv.Itoa(idx))
			js, err := BuildAndStartJudger(config.JudgerName+strconv.Itoa(idx), cli, cconfig, ccfg)
			if err != nil {
				errs = append(errs, err)
				return
			}
			dae.judgerPool <- js
		}(idx)
	}

	wg.Wait()

	fmt.Println("here")

	if len(errs) != 0 {
		dae = nil
		err = errs
		return
	}
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

type Errors []error

func (errs Errors) Error() string {
	var b = new(bytes.Buffer)
	for _, err := range errs {
		b.WriteString(err.Error())
	}
	return b.String()
}

func (dae *Daemon) Close() error {
	worker, ok := <-dae.judgerPool
	var errs Errors
	for ok {
		err := worker.Close()
		if err != nil {
			errs = append(errs, err)
		}
		worker, ok = <-dae.judgerPool
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
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
