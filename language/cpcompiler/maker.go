package cpcompiler

import (
	"io"
	"os"

	"github.com/Myriad-Dreamin/core-oj/types"
)

// Compiler is for python2/python3 configuration
type Compiler struct {
	Source string
	Target string
}

// Compile (Copy) source to destination, this function only remove the write permission of files
func (cp *Compiler) Compile(srcPath, tarPath string) error {
	src, err := os.Open(srcPath + cp.Source)
	if err != nil {
		var ce = new(types.CompileError)
		ce.ProcErr = err.Error()
		return ce
	}

	dst, err := os.OpenFile(tarPath+cp.Target, os.O_CREATE|os.O_WRONLY, 0774)
	if err != nil {
		src.Close()
		return err
	}
	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()
	return err
}
