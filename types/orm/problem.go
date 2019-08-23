package orm

import (
	"github.com/go-xorm/xorm"
)

// Problem example
type Problem struct {
	ID   int `xorm:"not null pk autoincr"`
	Name string

	TimeLimit   int // in milliseconds
	MemoryLimit int // in bytes

	TestPath string // config path
}

// TableName return the table name
func (obj *Problem) TableName() string {
	return "problems"
}

// GetSliceWithPredict return the slice of object with reserving the space of n Problem
func (obj *Problem) GetSliceWithPredict(n int) interface{} {
	return make([]Problem, 0, n)
}

// GetSlice return the slice of object
func (obj *Problem) GetSlice() interface{} {
	return make([]Problem, 0)
}

// Insert into Engine
func (obj *Problem) Insert() (int64, error) {
	return x.Insert(obj)
}

// Delete from Engine
func (obj *Problem) Delete() (int64, error) {
	return x.Id(obj.ID).Delete(obj)
}

// Update to Engine
func (obj *Problem) Update() (int64, error) {
	return x.Id(obj.ID).Update(obj)
}

// Query from Engine
func (obj *Problem) Query() (bool, error) {
	return x.Get(obj)
}

// Problemer Extend the Engine operation
type Problemer struct {
}

// ProblemSession Extend the Engine operation
type ProblemSession xorm.Session

// Query return the code with Property property
func (objx *Problemer) Query(property int) (*Problem, error) {
	obj := new(Problem)
	obj.ID = property
	has, err := x.Get(obj)
	if has {
		return obj, nil
	}
	return nil, err
}

// QueryName return the code with Property property
func (objx *Problemer) QueryName(property string) (*Problem, error) {
	obj := new(Problem)
	obj.Name = property
	has, err := x.Get(obj)
	if has {
		return obj, nil
	}
	return nil, err
}

// Inserts many objects
func (objx *Problemer) Inserts(objs []Problem) (int64, error) {
	return x.Insert(objs)
}

// Querys with conditions
func (objx *Problemer) Querys(objs []Problem, conds ...interface{}) error {
	return x.Find(objs, conds...)
}

// ColsQuerys with conditions with specifying columns
func (objx *Problemer) ColsQuerys(objs []Problem, cols ...string) error {
	return x.Cols(cols...).Find(objs)
}

// Where provides custom query condition.
func (objx *Problemer) Where(query interface{}, args ...interface{}) *ProblemSession {
	return (*ProblemSession)(x.Where(query, args...))
}

// Where provides custom query condition.
func (objx *ProblemSession) Where(query interface{}, args ...interface{}) *ProblemSession {
	return (*ProblemSession)(((*xorm.Session)(objx)).Where(query, args...))
}

// And provides custom query condition.
func (objx *ProblemSession) And(query interface{}, args ...interface{}) *ProblemSession {
	return (*ProblemSession)(((*xorm.Session)(objx)).And(query, args...))
}

// Or provides custom query condition.
func (objx *ProblemSession) Or(query interface{}, args ...interface{}) *ProblemSession {
	return (*ProblemSession)(((*xorm.Session)(objx)).Or(query, args...))
}

// ID provides custom query condition.
func (objx *ProblemSession) ID(query interface{}) *ProblemSession {
	return (*ProblemSession)(((*xorm.Session)(objx)).ID(query))
}

// NotIn provides custom query condition.
func (objx *ProblemSession) NotIn(query string, args ...interface{}) *ProblemSession {
	return (*ProblemSession)(((*xorm.Session)(objx)).NotIn(query, args...))
}

// In provides custom query condition.
func (objx *ProblemSession) In(query string, args ...interface{}) *ProblemSession {
	return (*ProblemSession)(((*xorm.Session)(objx)).In(query, args...))
}

// Find retrieve records from table, condiBeans's non-empty fields
// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
// map[int64]*Struct
func (objx *ProblemSession) Find(conds ...interface{}) ([]Problem, error) {
	objs := make([]Problem, 0)
	err := ((*xorm.Session)(objx)).Find(objs, conds...)
	return objs, err
}