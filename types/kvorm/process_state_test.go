package korm

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Myriad-Dreamin/core-oj/types"
)

func TestDecodeProcState(t *testing.T) {
	var g ProcState

	b, _ := json.Marshal(ProcState{
		0,
		types.WrongAnswer{Info: []byte("wrong answer"), ProcErr: "wrong answer"},
		100000000, 1555,
	})

	fmt.Println(json.Unmarshal(b, &g))
	fmt.Println(g.CodeError.ErrorCode(), string(g.CodeError.JudgeError()), g.CodeError)
}

func TestProcStatesStorage(t *testing.T) {
	var g []ProcState
	g = append(g, ProcState{
		0,
		types.WrongAnswer{Info: []byte("wrong answer"), ProcErr: "wrong answer"},
		100000000, 1555,
	})
	g = append(g, ProcState{
		0,
		types.AcceptedBaseCodeError(),
		10002000, 1755,
	})

	db, err := NewGobLevelDB("./test.db")
	if err != nil {
		t.Error(err)
		return
	}

	RegisterEngine(db)

	defer db.Close()
	var proc *ProcStater
	proc, err = NewProcStater()
	err = proc.Insert(7, g)
	if err != nil {
		t.Error(err)
		return
	}
	ss, err := proc.Query(7)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(ss)
	err = proc.Delete(7)
	if err != nil {
		t.Error(err)
		return
	}
}
