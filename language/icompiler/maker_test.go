package icompiler

import (
	"testing"

	"github.com/Myriad-Dreamin/core-oj/types"
)

func TestCompileCXX(t *testing.T) {
	if err := (&Compiler{Path: "g++", Args: []string{
		"-std=c++11",
	}}).Compile("./test/test.cpp", "./test/test"); err != nil {
		t.Error(err, "\n", string(err.(types.CodeError).JudgeError()))
		return
	}
}

func TestCompileC(t *testing.T) {
	if err := (&Compiler{Path: "g++"}).Compile("./test/test.c", "./test/test2"); err != nil {
		t.Error(err, "\n", string(err.(types.CodeError).JudgeError()))
		return
	}
}
