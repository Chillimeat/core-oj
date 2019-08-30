package argiexecutor

import (
	"os/exec"
	"strings"

	types "github.com/Myriad-Dreamin/core-oj/types"
)

// Executor is for c++(gcc) configuration
type Executor struct {
	Path     string
	DestName string
	Regp     string
}

func (et *Executor) Context() *types.ExecContext {
	return types.NewExecContext(et)
}

// Execute source in workdir
func (et *Executor) Execute(ctx *types.ExecContext, workdir string) (cmd *exec.Cmd, err error) {
	regp := strings.ReplaceAll(et.Regp, "{w}", workdir)
	regp = strings.ReplaceAll(regp, "{d}", et.DestName)
	if ctx.Context != nil {
		cmd = exec.CommandContext(ctx.Context, et.Path, strings.Fields(regp)...)
	} else {
		cmd = exec.Command(et.Path, strings.Fields(regp)...)
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
