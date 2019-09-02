package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	profiler "github.com/Myriad-Dreamin/core-oj/judger/judgerx/src"
	types "github.com/Myriad-Dreamin/core-oj/types"
	kvorm "github.com/Myriad-Dreamin/core-oj/types/kvorm"
)

var (
	unixAddress   = flag.String("addr", "/var/run/judger-test.sock", "ipc server address")
	eachCaseDelay = 500 * time.Millisecond
)

var lenBuffer int32
var buffer = make([]byte, 1024*128)
var precBody, suffBody = buffer[0:4], buffer[4:]
var precBuffer, suffBuffer = bytes.NewBuffer(precBody), bytes.NewBuffer(suffBody)

func serve(ctx context.Context, conn *net.UnixConn) {
	defer conn.Close()

	var receive = func() []byte {
		err := binary.Read(conn, binary.BigEndian, &lenBuffer)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		_, err = io.ReadFull(conn, buffer[0:lenBuffer])
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return buffer[0:lenBuffer]
	}

	var send = func(b []byte) {
		precBuffer.Reset()
		suffBuffer.Reset()
		binary.Write(precBuffer, binary.BigEndian, int32(len(b)))
		binary.Write(suffBuffer, binary.BigEndian, b)
		_, err := conn.Write(buffer[0 : 4+len(b)])
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	for {
		conn.SetDeadline(time.Now().Add(eachCaseDelay))
		b := receive()
		if b == nil {
			return
		}

		var testCase = new(types.TestCase)

		err := json.Unmarshal(buffer[0:lenBuffer], testCase)
		if err != nil {
			b, err := json.Marshal(&kvorm.ProcState{0, types.SystemError{ProcErr: err.Error()}, 0, 0})
			if err != nil {
				panic(err)
			}
			send(b)
			return
		}

		// b, err := json.Marshal(testCase)
		// fmt.Println(string(s)tring(b), err)

		input, err := profiler.MakeInput(testCase)
		if err != nil {
			b, err = json.Marshal(&kvorm.ProcState{0, types.SystemError{ProcErr: err.Error()}, 0, 0})
			if err != nil {
				panic(err)
			}
			send(b)
			return
		}

		conn.SetDeadline(time.Now().Add(testCase.TimeLimit + eachCaseDelay))
		sctx, cancel := context.WithDeadline(ctx, time.Now().Add(testCase.TimeLimit+eachCaseDelay))
		var output = new(bytes.Buffer)

		procInfo := profiler.Profile(sctx, testCase, input, output)
		cancel()
		input.Close()

		if procInfo != nil {

			if procInfo.CodeError == nil {
				// todo spj time limit
				conn.SetDeadline(time.Now().Add(testCase.TimeLimit + eachCaseDelay))
				sctx, cancel = context.WithDeadline(ctx, time.Now().Add(testCase.TimeLimit+eachCaseDelay))
				procInfo.CodeError = profiler.Check(sctx, testCase, output)
			}
			cancel()
			b, err = json.Marshal(procInfo)
			if err != nil {
				fmt.Println(*procInfo)
				fmt.Printf("%v %T\n", procInfo.CodeError, procInfo.CodeError)
				panic(err)
			}
			send(b)
		} else {
			panic("empty procInfo")
		}
	}
}

func init() {
	flag.Parse()
}

func main() {

	os.Remove(*unixAddress)
	uaddr, err := net.ResolveUnixAddr("unix", *unixAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	unixListener, err := net.ListenUnix("unix", uaddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.Chmod(*unixAddress, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer unixListener.Close()

	fmt.Println("qwq...")

	for {
		conn, err := unixListener.AcceptUnix()
		if err != nil {
			fmt.Println(err)
			continue
		}
		serve(context.Background(), conn)
	}
}

// {"cn":1,"ops":0,"inp":"in.txt","soup":"ans.txt","tl":2001000000,"ml":262144,"spj":false}
// testCase := &types.TestCase{
// 	CaseNumber:    1,
// 	OptionStream:  0,
// 	InputPath:     "in.txt",
// 	StdOutputPath: "ans.txt",
// 	TimeLimit:     2001 * time.Millisecond,
// 	MemoryLimit:   256 * 1024,
// 	SpecialJudge:  false,
// }
