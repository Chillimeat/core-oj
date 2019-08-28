package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	fmt.Println(json.Marshal(WrongAnswer{Info: []byte("wrong answer"), ProcErr: errors.New("wrong answer")}))
}
