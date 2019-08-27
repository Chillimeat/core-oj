package types

type CodeError interface {
	error
	JudgeError() []byte
	ErrorCode() int
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

func (err CompileError) ErrorCode() int {
	return StatusCompileError
}

type TimeLimitExceed struct {
	ProcErr error
}

func (err TimeLimitExceed) Error() string {
	return err.ProcErr.Error()
}

var tle = []byte("Time limit exceed")

func (err TimeLimitExceed) JudgeError() []byte {
	return tle
}

func (err TimeLimitExceed) ErrorCode() int {
	return StatusTimeLimitExceed
}

type MemoryLimitExceed struct {
	ProcErr error
}

func (err MemoryLimitExceed) Error() string {
	return err.ProcErr.Error()
}

var mle = []byte("Memory limit exceed")

func (err MemoryLimitExceed) JudgeError() []byte {
	return mle
}

func (err MemoryLimitExceed) ErrorCode() int {
	return StatusMemoryLimitExceed
}

type RuntimeError struct {
	ProcErr error
}

func (err RuntimeError) Error() string {
	return err.ProcErr.Error()
}

var re = []byte("Runtime error")

func (err RuntimeError) JudgeError() []byte {
	return re
}

func (err RuntimeError) ErrorCode() int {
	return StatusRuntimeError
}

type JudgeError struct {
	ProcErr error
}

func (err JudgeError) Error() string {
	return err.ProcErr.Error()
}

var je = []byte("Judge error")

func (err JudgeError) JudgeError() []byte {
	return je
}

func (err JudgeError) ErrorCode() int {
	return StatusJudgeError
}

type PresentationError struct {
	Info    []byte
	ProcErr error
}

func (err PresentationError) Error() string {
	return err.ProcErr.Error()
}

func (err PresentationError) JudgeError() []byte {
	return err.Info
}

func (err PresentationError) ErrorCode() int {
	return StatusPresentationError
}

type WrongAnswer struct {
	Info    []byte
	ProcErr error
}

func (err WrongAnswer) Error() string {
	return err.ProcErr.Error()
}

func (err WrongAnswer) JudgeError() []byte {
	return err.Info
}

func (err WrongAnswer) ErrorCode() int {
	return StatusWrongAnswer
}

type UnknownError struct {
	Info    []byte
	ProcErr error
}

func (err UnknownError) Error() string {
	return err.ProcErr.Error()
}

func (err UnknownError) JudgeError() []byte {
	return err.Info
}

func (err UnknownError) ErrorCode() int {
	return StatusUnknownError
}

type OutputLimitExceed struct {
	Info    []byte
	ProcErr error
}

func (err OutputLimitExceed) Error() string {
	return err.ProcErr.Error()
}

func (err OutputLimitExceed) JudgeError() []byte {
	return err.Info
}

func (err OutputLimitExceed) ErrorCode() int {
	return StatusOutputLimitExceed
}

type ExhaustedMatch struct {
	ProcErr error
}

func (err ExhaustedMatch) Error() string {
	return err.ProcErr.Error()
}

var em = []byte("Exhausted match of testlib")

func (err ExhaustedMatch) JudgeError() []byte {
	return em
}

func (err ExhaustedMatch) ErrorCode() int {
	return StatusExhaustedMatch
}

type SystemError struct {
	ProcErr error
}

func (err SystemError) Error() string {
	return err.ProcErr.Error()
}

var se = []byte("System error")

func (err SystemError) JudgeError() []byte {
	return se
}

func (err SystemError) ErrorCode() int {
	return StatusSystemError
}
