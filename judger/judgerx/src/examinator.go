package profiler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/Myriad-Dreamin/core-oj/types"
)

func MakeInput(t *types.TestCase) (io.ReadCloser, error) {
	if _, err := os.Stat(t.InputPath); err != nil {
		return nil, err
	}
	return os.Open(t.InputPath)
}

func Check(ctx context.Context, t *types.TestCase, output io.Reader) types.CodeError {
	if _, err := os.Stat(t.StdOutputPath); err != nil {
		return types.SystemError{ProcErr: err.Error()}
	}
	if t.SpecialJudge {
		sctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		cmd := exec.CommandContext(sctx, t.StdOutputPath)
		var jerrput = new(bytes.Buffer)
		cmd.Stdout = jerrput
		cmd.Stdin = output
		if err := cmd.Run(); err != nil {
			cancel()
			return types.JudgeError{ProcErr: err.Error()}
		}
		cancel()
		var testStatus int64
		fmt.Fscanf(jerrput, "%d", &testStatus)
		if testStatus != 0 {
			return types.ConstructCodeError(testStatus)(bytes.TrimSpace(jerrput.Bytes()))
		}
		return types.JudgeError{ProcErr: "not read test status"}
	}
	sctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	cmd := exec.CommandContext(sctx, "/judger_tools/sos-checker", "presu.txt", "stdin", t.StdOutputPath)
	var jerrout = new(bytes.Buffer)
	cmd.Stderr = jerrout
	cmd.Stdin = output
	if err := cmd.Run(); err != nil {
		cancel()
		return types.ConstructCodeErrorWithTestLib(jerrout)
	}
	cancel()
	return types.ConstructCodeErrorWithTestLib(jerrout)
}
