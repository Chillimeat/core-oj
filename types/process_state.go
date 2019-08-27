package types

import "time"

// ProcState record the cost of process
type ProcState struct {
	Status    int
	CodeError CodeError
	TimeUsed  time.Duration

	// in kilobytes
	MemoryUsed float64
}
