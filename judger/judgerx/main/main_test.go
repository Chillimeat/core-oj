package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	uaddr, err := net.ResolveUnixAddr("unix", "./qwq.sock")
	if err != nil {
		t.Error(err)
		return
	}

	conn, err := net.DialUnix("unix", nil, uaddr)
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	var lenBuffer int32
	var buffer = make([]byte, 3333)
	var precBuffer, suffBuffer = bytes.NewBuffer(buffer[0:4]), bytes.NewBuffer(buffer[4:])

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

	send([]byte(`{"cn":1,"tp":"/codes/main","ops":0,"inp":"in.txt","soup":"ans.txt","tl":2001000000,"ml":262144,"spj":false}`))
	fmt.Println(string(receive()))
}
