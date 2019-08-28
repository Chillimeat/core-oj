package judger

import (
	"fmt"
	"testing"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"
	types "github.com/Myriad-Dreamin/core-oj/types"
)

func TestBuildAndStartJudger(t *testing.T) {
	cli, err := client.Connect("unix:///var/run/docker.sock", "v1.40")
	if err != nil {
		t.Error(err)
		return
	}
	config := client.NewContainerConfig()
	cconfig := &types.JudgerConfig{UnixAddress: "/home/kamiyoru/data/judger_tools/socks/judger.sock"}
	config.VolumeMap.InsertBind("/home/kamiyoru/data/test", "/codes")
	config.VolumeMap.InsertBind("/home/kamiyoru/data/judger_tools", "/judger_tools")
	config.VolumeMap.InsertBind("/home/kamiyoru/data/problems", "/problems")
	config.VolumeMap.InsertBind("/home/kamiyoru/data/checker_tools", "/checker_tools")
	config.Env = append(config.Env, "NAME=judger")
	js, err := BuildAndStartJudger("judger", cli, cconfig, config)
	if err != nil {
		t.Error(err)
		return
	}

	js.conn.Send([]byte(`{"cn":1,"tp":"/codes/main","ops":0,"inp":"/problems/1001/in.txt","soup":"/problems/1001/ans.txt","tl":2001000000,"ml":262144,"spj":false}`))
	b, err := js.conn.Receive()
	fmt.Println(string(b), err)
}
