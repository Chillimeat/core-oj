package profiler

import (
	"bytes"
	"context"
	"errors"
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

func Check(t *types.TestCase, output io.Reader) types.CodeError {
	if _, err := os.Stat(t.StdOutputPath); err != nil {
		return types.SystemError{ProcErr: err}
	}
	if t.SpecialJudge {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		cmd := exec.CommandContext(ctx, t.StdOutputPath)
		var jerrput = new(bytes.Buffer)
		cmd.Stdout = jerrput
		cmd.Stdin = output
		if err := cmd.Run(); err != nil {
			cancel()
			return types.JudgeError{ProcErr: err}
		}
		cancel()
		var testStatus int
		fmt.Fscanf(jerrput, "%d", &testStatus)
		if testStatus != 0 {
			return types.ConstructCodeError(testStatus)(bytes.TrimSpace(jerrput.Bytes()))
		}
		return types.JudgeError{ProcErr: errors.New("not read test status")}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cmd := exec.CommandContext(ctx, "/sos-checker", "presu.txt", "stdin", t.StdOutputPath)
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
