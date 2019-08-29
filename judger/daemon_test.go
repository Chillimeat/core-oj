package judger

import (
	"context"
	"fmt"
	"testing"
	"time"

	config "github.com/Myriad-Dreamin/core-oj/config"
	client "github.com/Myriad-Dreamin/core-oj/docker-client"
	morm "github.com/Myriad-Dreamin/core-oj/types/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func TestDaemon(t *testing.T) {
	var err error
	engine, err := xorm.NewEngine(config.Config().DriverName, config.Config().MasterDataSourceName)
	if err != nil {
		t.Error("prepare failed", "error", err)
		return
	}

	morm.RegisterEngine(engine)

	cli, err := client.Connect("unix:///var/run/docker.sock", "v1.40")
	if err != nil {
		t.Error(err)
		return
	}

	dae, err := NewDaemon(cli)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	orzz := make(chan bool)

	problemer, err := morm.NewProblemer()
	if err != nil {
		t.Error(err)
		return
	}

	// affected, err := (&morm.Problem{
	// 	Name:     "A+B problem",
	// 	OwnerUID: 1,
	// }).Insert()
	// if affected == 0 || err != nil {
	// 	t.Error(affected, err)
	// 	return
	// }
	// fmt.Println(affected)

	go dae.Run(ctx, func(js *Judger) {
		defer func() {
			time.Sleep(time.Second * 1)
			orzz <- true
		}()
		problem, err := problemer.Query(1)
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(problem)
		code := new(morm.Code)
		code.Hash = []byte("")
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(code)
		results, err := js.Judge(code, problem)
		if err != nil {
			t.Error(err)
			return
		}
		for _, result := range results {
			fmt.Println(*result)
		}
	})

	fmt.Println("running")

	<-orzz
	<-orzz
}
