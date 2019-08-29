package types

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	fmt.Println(json.Marshal(WrongAnswer{Info: []byte("wrong answer"), ProcErr: "wrong answer"}))
}

func TestDecode(t *testing.T) {
	var g BaseCodeError

	b, _ := json.Marshal(WrongAnswer{Info: []byte("wrong answer"), ProcErr: "wrong answer"})

	fmt.Println(json.Unmarshal(b, &g))
	fmt.Println(g.ErrorCode(), g.JErr, g.Err)
}
