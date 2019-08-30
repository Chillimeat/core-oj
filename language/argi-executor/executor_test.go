package argiexecutor

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/Myriad-Dreamin/core-oj/types"
)

func TestExecutePython(t *testing.T) {
	var e = (&Executor{"go", "main.go", "run {w}/{d}"})
	var c = e.Context()
	c.Stdin = bytes.NewBufferString("3")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if cmd, err := c.Execute("./test/"); err != nil {
		t.Error(err, "\n", string(err.(types.CodeError).JudgeError()))
		// return
	} else {
		fmt.Println(cmd)
	}

}

// func TestExecuteC(t *testing.T) {
// 	if err := (&Executor{Path: "g++"}).Execute("./test/test.c", "./test/test2"); err != nil {
// 		t.Error(err, "\n", string(err.(types.CodeError).JudgeError()))
// 		return
// 	}
// }
