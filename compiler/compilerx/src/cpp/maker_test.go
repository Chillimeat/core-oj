package cpp

import (
	"testing"

	"github.com/Myriad-Dreamin/core-oj/types"
)

func TestCompile(t *testing.T) {
	if err := (&Compiler{Path: "gcc", Args: []string{
		"-std=c++11",
	}}).Compile("./test/test.cpp", "./test/test"); err != nil {
		t.Error(err, "\n", string(err.(types.CodeError).JudgeError()))
		return
	}
}
