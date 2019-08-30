package iocompiler

import (
	"bytes"
	"io/ioutil"
	"os/exec"

	"github.com/Myriad-Dreamin/core-oj/types"
)

// Compiler is for c++(gcc) configuration
type Compiler struct {
	Path   string
	Source string
	Target string
	Args   []string
}

// Compile source to destination
func (cp *Compiler) Compile(srcPath, tarPath string) error {
	var err error

	var iocp = exec.Command(cp.Path, append(cp.Args, []string{srcPath + cp.Source, "-o", tarPath + cp.Target}...)...)

	var stderr = bytes.NewBuffer(make([]byte, 0, 200))
	iocp.Stderr = stderr

	if err = iocp.Run(); err != nil {
		var ce = new(types.CompileError)
		ce.ProcErr = err.Error()
		ce.Info, err = ioutil.ReadAll(stderr)
		return ce
	}
	return nil
}
