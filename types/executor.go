package types

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"
)

type Executor interface {
	Context() *ExecContext
	Execute(ctx *ExecContext, workdir string) (cmd *exec.Cmd, err error)
}

type ExecContext struct {
	et          Executor
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	Context     context.Context
	SysProcAttr *syscall.SysProcAttr
	ExtraFiles  []*os.File
}

func (ec *ExecContext) Execute(workdir string) (cmd *exec.Cmd, err error) {
	return ec.et.Execute(ec, workdir)
}

func NewExecContext(et Executor) *ExecContext {
	return &ExecContext{et: et}
}
