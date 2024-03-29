package orm

import (
	"log"

	"github.com/go-xorm/xorm"
)

// must register
var x *xorm.Engine

// RegisterEngine to store objects
func RegisterEngine(y *xorm.Engine) {
	x = y

	if err := x.Sync(new(Code)); err != nil {
		log.Fatal("Syn Error: Code:", err)
	}

	if err := x.Sync(new(User)); err != nil {
		log.Fatal("Syn Error: User:", err)
	}

	if err := x.Sync(new(Problem)); err != nil {
		log.Fatal("Syn Error: Problem:", err)
	}

	if err := x.Sync(new(RuntimeProblem)); err != nil {
		log.Fatal("Syn Error: RuntimeProblem:", err)
	}

}
