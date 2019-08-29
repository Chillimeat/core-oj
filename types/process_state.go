package types

import (
	"encoding/json"
	"reflect"
	"time"
	"unsafe"

	"github.com/buger/jsonparser"
)

// ProcState record the cost of process
type ProcState struct {
	Status    int64
	CodeError CodeError
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
	g.CodeError = AcceptedBaseCodeError()
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
