package cpp

import (
	"bytes"
	"io/ioutil"
	"os/exec"

	"github.com/Myriad-Dreamin/core-oj/types"
)

// Compiler is for c++(gcc) configuration
type Compiler struct {
	Path string
	Args []string
}

// Compile source to destination
func (cp *Compiler) Compile(source, destination string) error {
	var err error

	var gccer = exec.Command(cp.Path, append(cp.Args, []string{source, "-o", destination}...)...)

	var stderr = bytes.NewBuffer(make([]byte, 65536))
	gccer.Stderr = stderr

	if err = gccer.Run(); err != nil {
		var ce = new(types.CompileError)
		ce.ProcErr = err.Error()
		ce.Info, err = ioutil.ReadAll(stderr)
		if err != nil {
			ce.Info = []byte(err.Error())
		}
		return ce
	}
	return nil
}
