package helpfunc

import (
	"bytes"
	"sync"
)

type Errors []error

func (errs Errors) Error() string {
	var b = new(bytes.Buffer)
	for _, err := range errs {
		b.WriteString(err.Error())
	}
	return b.String()
}

func Walk(n int, work func(int) error) (err error) {
	var wg sync.WaitGroup
	var errs Errors
	for idx := 0; idx < n; idx++ {
		wg.Add(1)
		go func(idx int) {
			err := work(idx)
			if err != nil {
				errs = append(errs, err)
			}
			wg.Done()
		}(idx)
	}
	wg.Wait()

	if len(errs) != 0 {
		err = errs
		return
	}
	return
}
