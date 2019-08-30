package icompiler

import (
	"bytes"
	"io/ioutil"
	"os/exec"

	"github.com/Myriad-Dreamin/core-oj/types"
)

// Compiler is for c++(gcc) configuration
type Compiler struct {
	Path   string
	Args   []string
	Source string
}

// Compile source to destination
func (cp *Compiler) Compile(srcPath, _ string) error {
	var err error

	var icp = exec.Command(cp.Path, append(cp.Args, []string{srcPath + cp.Source}...)...)

	var stderr = bytes.NewBuffer(make([]byte, 0, 200))
	icp.Stderr = stderr

	if err = icp.Run(); err != nil {
		var ce = new(types.CompileError)
		ce.ProcErr = err.Error()
		ce.Info, err = ioutil.ReadAll(stderr)
		return ce
	}
	return nil
}
