package profiler

import (
	"context"
	"io"
	"os/exec"
	"syscall"
	"time"

	types "github.com/Myriad-Dreamin/core-oj/types"
)

func itoa(val int) string { // do it here rather than with fmt to avoid dependency
	if val < 0 {
		return "-" + uitoa(uint(-val))
	}
	return uitoa(uint(val))
}

func uitoa(val uint) string {
	var buf [32]byte // big enough for int64
	i := len(buf) - 1
	for val >= 10 {
		buf[i] = byte(val%10 + '0')
		i--
		val /= 10
	}
	buf[i] = byte(val + '0')
	return string(buf[i:])
}

// BSDMaxrss is in kilobytes
const BSDMaxrss = 1024

// Profile execute a command and profile it
func Profile(ctx context.Context, testCase *types.TestCase, input io.Reader, output io.Writer) *types.ProcState {
	sctx, cancel := context.WithTimeout(ctx, testCase.TimeLimit)
	cmd := exec.CommandContext(sctx, testCase.TestPath)
	cmd.Stdout = output
	cmd.Stdin = input
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(23333), Gid: uint32(23333)}
	if err := cmd.Run(); err != nil {
		cancel()
		if errr, ok := err.(*exec.ExitError); ok {

			rusage := errr.SysUsage().(*syscall.Rusage)
			status := errr.Sys().(syscall.WaitStatus)
			res := ""
			switch {
			case status.Exited():
				return &types.ProcState{
					int64(status.ExitStatus()),
					types.RuntimeError{ProcErr: "exit status " + itoa(status.ExitStatus())},
					time.Microsecond*time.Duration(rusage.Utime.Usec) + time.Second*time.Duration(rusage.Utime.Sec) + time.Microsecond*time.Duration(rusage.Stime.Usec) + time.Second*time.Duration(rusage.Stime.Sec), float64(rusage.Maxrss*BSDMaxrss) / 1024.0,
				}
			case status.Signaled():
				res = "signal: " + status.Signal().String()
			case status.Stopped():
				res = "stop signal: " + status.StopSignal().String()
				if status.StopSignal() == syscall.SIGTRAP && status.TrapCause() != 0 {
					res += " (trap " + itoa(status.TrapCause()) + ")"
				}
			case status.Continued():
				res = "continued"
			}

			return &types.ProcState{
				int64(cmd.ProcessState.ExitCode()),
				types.TimeLimitExceed{ProcErr: res},
				time.Microsecond*time.Duration(rusage.Utime.Usec) + time.Second*time.Duration(rusage.Utime.Sec) + time.Microsecond*time.Duration(rusage.Stime.Usec) + time.Second*time.Duration(rusage.Stime.Sec), float64(rusage.Maxrss*BSDMaxrss) / 1024.0,
			}
		}

		return &types.ProcState{
			int64(23333),
			types.JudgeError{ProcErr: err.Error()},
			0, 0,
		}
	}
	cancel()

	rusage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	var timeUsed = time.Microsecond*time.Duration(rusage.Utime.Usec) + time.Second*time.Duration(rusage.Utime.Sec) + time.Microsecond*time.Duration(rusage.Stime.Usec) + time.Second*time.Duration(rusage.Stime.Sec)

	if timeUsed > testCase.TimeLimit {
		return &types.ProcState{
			int64(cmd.ProcessState.ExitCode()), types.TimeLimitExceed{ProcErr: "Time limit exceed"},
			timeUsed,
			float64(rusage.Maxrss*BSDMaxrss) / 1024.0,
		}
	}

	var MemoryUsed = float64(rusage.Maxrss*BSDMaxrss) / 1024.0
	if MemoryUsed > float64(testCase.MemoryLimit) {
		return &types.ProcState{
			int64(cmd.ProcessState.ExitCode()), types.MemoryLimitExceed{ProcErr: "Memory limit exceed"},
			timeUsed,
			float64(rusage.Maxrss*BSDMaxrss) / 1024.0,
		}
	}

	return &types.ProcState{
		int64(cmd.ProcessState.ExitCode()), nil,
		timeUsed,
		float64(rusage.Maxrss*BSDMaxrss) / 1024.0,
	}
}
