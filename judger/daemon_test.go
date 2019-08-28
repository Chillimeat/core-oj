package judger

import (
	"context"
	"fmt"
	"testing"

	client "github.com/Myriad-Dreamin/core-oj/docker-client"
)

func TestDaemon(t *testing.T) {
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
	go dae.Run(ctx, func(js *Judger) {

		// select {
		// case <-ctx.Done():
		// 	return
		// }

		// if _, ok := <-ctx.Done(); ok {
		// 	return
		// }
		js.conn.Send([]byte(`{"cn":1,"tp":"/codes/main","ops":0,"inp":"/problems/1001/in.txt","soup":"/problems/1001/ans.txt","tl":2001000000,"ml":262144,"spj":false}`))
		b, err := js.conn.Receive()
		fmt.Println(string(b), err)
		orzz <- true
	})

	fmt.Println("running")

	<-orzz
}
