package types

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

const (
	StatusWaitingForJudge int = iota
	StatusAccepted
	StatusRunning
	StatusCompiling
	StatusCompileError
	StatusCompileTimeout
	StatusWrongAnswer
	StatusTimeLimitExceed
	StatusMemoryLimitExceed
	StatusOutputLimitExceed
	StatusSystemError
	StatusUnknownError
	StatusPresentationError
	StatusRuntimeError
	StatusJudgeError
	StatusExhaustedMatch
)

// StatusConstructor apply the info to code error
type StatusConstructor func(Info []byte) CodeError

// ConstructCodeError return a constrcutor to generate code error
func ConstructCodeError(status int) StatusConstructor {
	switch status {
	case StatusAccepted:
		return func([]byte) CodeError {
			return nil
		}
	case StatusWrongAnswer:
		return func(info []byte) CodeError {
			return &WrongAnswer{Info: info, ProcErr: errors.New("Wrong answer")}
		}
	case StatusUnknownError:
		return func(info []byte) CodeError {
			return &UnknownError{Info: info, ProcErr: errors.New("Unknown error")}
		}
	case StatusPresentationError:
		return func(info []byte) CodeError {
			return &PresentationError{Info: info, ProcErr: errors.New("Presentation error")}
		}
	case StatusOutputLimitExceed:
		return func(info []byte) CodeError {
			return &OutputLimitExceed{Info: info, ProcErr: errors.New("Output limit exceed")}
		}
	case StatusJudgeError:
		return func([]byte) CodeError {
			return &JudgeError{ProcErr: errors.New("Judge error")}
		}
	case StatusSystemError:
		return func(info []byte) CodeError {
			return &SystemError{ProcErr: errors.New(string(info))}
		}
	default:
		return func([]byte) CodeError {
			return &JudgeError{fmt.Errorf("Unknown status of special judge? %v", status)}
		}
	}
}

func ConstructCodeErrorWithTestLib(errBuf io.Reader) CodeError {
	buf := bufio.NewReader(errBuf)
	status := matchStatus(buf)
	b, _ := ioutil.ReadAll(buf)
	return ConstructCodeError(status)(bytes.TrimSpace(b))
}

const (
	firstState = iota
	secondState
)

func calc(t string) (ret int) {
	for _, b := range t {
		ret = ret*233 + int(b)
	}
	return
}

var (
	wrongHash      = calc("wrong ")
	outputHash     = calc("output ")
	okHash         = calc("ok ")
	failHash       = calc("FAIL ")
	pointsHash     = calc("points ")
	unexpectedHash = calc("unexpected ")
	partiallyHash  = calc("partially ")
)

func matchStatus(buf *bufio.Reader) int {
	var outVal int
	var aut = 0
	for {
		bs, err := buf.ReadBytes(' ')
		if err == io.EOF {
			return StatusExhaustedMatch
		}
		outVal = 0
		for _, b := range bs {
			outVal = outVal*233 + int(b)
		}
		switch aut {
		case firstState:
			if outVal == wrongHash {
				aut = secondState
			} else {
				switch outVal {
				case okHash:
					return StatusAccepted
				case failHash:
					return StatusJudgeError
				case unexpectedHash:
					_, err = buf.ReadBytes(' ')
					if err == io.EOF {
						return StatusExhaustedMatch
					}
					return StatusOutputLimitExceed
				case pointsHash:
					return StatusWrongAnswer
				case partiallyHash:
					return StatusWrongAnswer
				default:
					return StatusSystemError
				}
			}
		case secondState:
			switch outVal {
			case outputHash:
				_, err = buf.ReadBytes(' ')
				if err == io.EOF {
					return StatusExhaustedMatch
				}
				return StatusPresentationError
			default:
				return StatusWrongAnswer
			}
		}
	}
}
