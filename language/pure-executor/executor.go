package pureexecutor

import (
	"os/exec"

	types "github.com/Myriad-Dreamin/core-oj/types"
)

// Executor is for c++(gcc) configuration
type Executor struct {
	DestName string
}

func (et *Executor) Context() *types.ExecContext {
	return types.NewExecContext(et)
}

// Execute source to workdir
func (et *Executor) Execute(ctx *types.ExecContext, workdir string) (cmd *exec.Cmd, err error) {
	if ctx.Context != nil {
		cmd = exec.CommandContext(ctx.Context, workdir+et.DestName)
	} else {
		cmd = exec.Command(workdir + et.DestName)
	}

	cmd.Stdin = ctx.Stdin
	cmd.Stdout = ctx.Stdout
	cmd.Stderr = ctx.Stderr
	cmd.SysProcAttr = ctx.SysProcAttr
	cmd.ExtraFiles = ctx.ExtraFiles

	if err = cmd.Run(); err != nil {
		return
	}
	return
}
