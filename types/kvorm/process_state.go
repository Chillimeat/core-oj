package korm

import (
	"encoding/json"
	"reflect"
	"time"
	"unsafe"

	"github.com/Myriad-Dreamin/core-oj/types"
	"github.com/buger/jsonparser"
)

// ProcState record the cost of process
type ProcState struct {
	Status    int64
	CodeError types.CodeError
	TimeUsed  time.Duration

	// in kilobytes
	MemoryUsed float64
}

func stringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func (g *ProcState) UnmarshalJSON(b []byte) (err error) {
	g.CodeError = types.AcceptedBaseCodeError()
	c, err := jsonparser.GetUnsafeString(b, "CodeError")
	if err != nil {
		return
	}
	err = json.Unmarshal(stringToBytes(c), g.CodeError)
	if err != nil {
		return
	}
	g.Status, err = jsonparser.GetInt(b, "TimeUsed")
	if err != nil {
		return
	}
	g.TimeUsed = time.Duration(g.Status)
	g.Status, err = jsonparser.GetInt(b, "Status")
	if err != nil {
		return
	}
	g.MemoryUsed, err = jsonparser.GetFloat(b, "MemoryUsed")
	if err != nil {
		return
	}
	return nil
}

type ProcStater struct {
	db KVDB
}

var procStateSlot = []byte("procstate")

func NewProcStater(db KVDB) (*ProcStater, error) {
	db, err := db.Table(procStateSlot)
	if err != nil {
		return nil, err
	}
	return &ProcStater{db}, nil
}

func (objx *ProcStater) Query(id int) ([]ProcState, error) {
	var procStates []ProcState
	err := objx.db.IDThen(id).Extract(&procStates)
	if err != nil {
		return nil, err
	}
	return procStates, nil
}

func (objx *ProcStater) Insert(id int, obj []ProcState) error {
	return objx.db.IDThen(id).Push(obj)
}

func (objx *ProcStater) Delete(id int) error {
	return objx.db.IDThen(id).Clear()
}
