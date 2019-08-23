package types

type CodeError interface {
	error
	JudgeError() []byte
	ErrorCode() int32
}

type CompileError struct {
	Info    []byte
	ProcErr error
}

func (err CompileError) Error() string {
	return err.ProcErr.Error()
}

func (err CompileError) JudgeError() []byte {
	return err.Info
}

func (err CompileError) ErrorCode() int32 {
	return 3
}
