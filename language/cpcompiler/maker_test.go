package cpcompiler

import (
	"testing"

	"github.com/Myriad-Dreamin/core-oj/types"
)

func TestCompilePython(t *testing.T) {
	if err := (&Compiler{}).Compile("./test/test.py", "./test/main"); err != nil {
		t.Error(err, "\n", string(err.(types.CodeError).JudgeError()))
		return
	}
}
