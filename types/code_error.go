package types

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"strings"
)

type CodeError interface {
	error
	JudgeError() []byte
	ErrorCode() int64
}

type BaseCodeError struct {
	Err  string `json:"er"`
	Code int64  `json:"ec"`
	JErr string `json:"je"`
}

func (err BaseCodeError) Error() string {
	return err.Err
}

func (err BaseCodeError) ErrorCode() int64 {
	return err.Code
}

func (err BaseCodeError) JudgeError() []byte {
	return []byte(err.JErr)
}

func AcceptedBaseCodeError() CodeError {
	return &BaseCodeError{Code: StatusAccepted}
}

var ba = []byte(`"`)
var bb = []byte(`\"`)

func marshalCodeError(c CodeError) []byte {
	var buf = bytes.NewBuffer(make([]byte, 70))
	buf.Reset()
	buf.WriteByte('{')

	buf.WriteString(`"er":"`)
	buf.WriteString(strings.Replace(c.Error(), `"`, `\"`, -1))

	if b := c.JudgeError(); b != nil {
		buf.WriteString(`","je":"`)
		buf.Write(bytes.Replace(b, ba, bb, -1))
	}
	buf.WriteString(`","ec":`)
	buf.WriteString(strconv.FormatInt(c.ErrorCode(), 10))

	buf.WriteByte('}')

	// fmt.Println(buf.String(), len(buf.String()), buf.Len())

	return buf.Bytes()
}

type CompileError struct {
	Info    []byte
	ProcErr string
}

func (err CompileError) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	return
}

func (err CompileError) Error() string {
	return err.ProcErr
}

func (err CompileError) JudgeError() []byte {
	return err.Info
}

func (err CompileError) ErrorCode() int64 {
	return StatusCompileError
}

type TimeLimitExceed struct {
	ProcErr string
}

func (err TimeLimitExceed) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	return
}

func (err TimeLimitExceed) Error() string {
	return err.ProcErr
}

var tle = []byte("Time limit exceed")

func (err TimeLimitExceed) JudgeError() []byte {
	return tle
}

func (err TimeLimitExceed) ErrorCode() int64 {
	return StatusTimeLimitExceed
}

type MemoryLimitExceed struct {
	ProcErr string
}

func (err MemoryLimitExceed) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	return
}

func (err MemoryLimitExceed) Error() string {
	return err.ProcErr
}

var mle = []byte("Memory limit exceed")

func (err MemoryLimitExceed) JudgeError() []byte {
	return mle
}

func (err MemoryLimitExceed) ErrorCode() int64 {
	return StatusMemoryLimitExceed
}

type RuntimeError struct {
	ProcErr string
}

func (err RuntimeError) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	return
}

func (err RuntimeError) Error() string {
	return err.ProcErr
}

var re = []byte("Runtime error")

func (err RuntimeError) JudgeError() []byte {
	return re
}

func (err RuntimeError) ErrorCode() int64 {
	return StatusRuntimeError
}

type JudgeError struct {
	ProcErr string
}

func (err JudgeError) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	return
}

func (err JudgeError) Error() string {
	return err.ProcErr
}

var je = []byte("Judge error")

func (err JudgeError) JudgeError() []byte {
	return je
}

func (err JudgeError) ErrorCode() int64 {
	return StatusJudgeError
}

type PresentationError struct {
	Info    []byte
	ProcErr string
}

func (err PresentationError) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	return
}

func (err PresentationError) Error() string {
	return err.ProcErr
}

func (err PresentationError) JudgeError() []byte {
	return err.Info
}

func (err PresentationError) ErrorCode() int64 {
	return StatusPresentationError
}

type WrongAnswer struct {
	Info    []byte
	ProcErr string
}

func (err WrongAnswer) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	// fmt.Println(string(b))
	return
}

func (err WrongAnswer) Error() string {
	return err.ProcErr
}

func (err WrongAnswer) JudgeError() []byte {
	return err.Info
}

func (err WrongAnswer) ErrorCode() int64 {
	return StatusWrongAnswer
}

type UnknownError struct {
	Info    []byte
	ProcErr string
}

func (err UnknownError) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	// fmt.Println(string(b))
	return
}

func (err UnknownError) Error() string {
	return err.ProcErr
}

func (err UnknownError) JudgeError() []byte {
	return err.Info
}

func (err UnknownError) ErrorCode() int64 {
	return StatusUnknownError
}

type OutputLimitExceed struct {
	Info    []byte
	ProcErr string
}

func (err OutputLimitExceed) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	// fmt.Println(string(b))
	return
}

func (err OutputLimitExceed) Error() string {
	return err.ProcErr
}

func (err OutputLimitExceed) JudgeError() []byte {
	return err.Info
}

func (err OutputLimitExceed) ErrorCode() int64 {
	return StatusOutputLimitExceed
}

type ExhaustedMatch struct {
	ProcErr string
}

func (err ExhaustedMatch) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	// fmt.Println(string(b))
	return
}

func (err ExhaustedMatch) Error() string {
	return err.ProcErr
}

var em = []byte("Exhausted match of testlib")

func (err ExhaustedMatch) JudgeError() []byte {
	return em
}

func (err ExhaustedMatch) ErrorCode() int64 {
	return StatusExhaustedMatch
}

type SystemError struct {
	ProcErr string
}

func (err SystemError) Error() string {
	return err.ProcErr
}

func (err SystemError) MarshalJSON() (b []byte, errr error) {
	b = marshalCodeError(err)
	// fmt.Println(string(b))
	return
}

var se = []byte("System error")

func (err SystemError) JudgeError() []byte {
	return se
}

func (err SystemError) ErrorCode() int64 {
	return StatusSystemError
}

func init() {
	gob.Register(BaseCodeError{})
	gob.Register(CompileError{})
	gob.Register(TimeLimitExceed{})
	gob.Register(WrongAnswer{})
	gob.Register(MemoryLimitExceed{})
	gob.Register(SystemError{})
	gob.Register(OutputLimitExceed{})
	gob.Register(PresentationError{})
	gob.Register(RuntimeError{})
	gob.Register(UnknownError{})
}
