package orm

import (
	"github.com/go-xorm/xorm"
)

// RuntimeProblem example
type RuntimeProblem struct {
	ID int `xorm:"not null pk autoincr"`
}

// TableName return the table name
func (obj *RuntimeProblem) TableName() string {
	return "runtime_problem"
}

// GetSliceWithPredict return the slice of RuntimeProblem with reserving the space of n RuntimeProblem
func (obj *RuntimeProblem) GetSliceWithPredict(n int) interface{} {
	return make([]RuntimeProblem, 0, n)
}

// GetSlice return the slice of RuntimeProblem
func (obj *RuntimeProblem) GetSlice() interface{} {
	return make([]RuntimeProblem, 0)
}

// Insert into Engine
func (obj *RuntimeProblem) Insert() (int64, error) {
	return x.Insert(obj)
}

// Delete from Engine
func (obj *RuntimeProblem) Delete() (int64, error) {
	return x.ID(obj.ID).Delete(obj)
}

// Update to Engine
func (obj *RuntimeProblem) Update() (int64, error) {
	return x.ID(obj.ID).Update(obj)
}

// Query from Engine
func (obj *RuntimeProblem) Query() (bool, error) {
	return x.Get(obj)
}

// RuntimeProblemer Extend the Engine operation
type RuntimeProblemer struct {
}

// RuntimeProblemSession Extend the Engine operation
type RuntimeProblemSession xorm.Session

// Query return the runtime problem with Property property
func (objx *RuntimeProblemer) Query(property int) (*RuntimeProblem, error) {
	obj := new(RuntimeProblem)
	obj.ID = property
	has, err := x.Get(obj)
	if has {
		return obj, nil
	}
	return nil, err
}

// QueryName return the runtime problem with Property property
// func (objx *RuntimeProblemer) QueryName(property string) (*RuntimeProblem, error) {
// 	obj := new(RuntimeProblem)
// 	obj.Name = property
// 	has, err := x.Get(obj)
// 	if has {
// 		return obj, nil
// 	}
// 	return nil, err
// }

// Inserts many Users
func (objx *RuntimeProblemer) Inserts(objs []RuntimeProblem) (int64, error) {
	return x.Insert(objs)
}

// Querys with conditions
func (objx *RuntimeProblemer) Querys(objs []RuntimeProblem, conds ...interface{}) error {
	return x.Find(objs, conds...)
}

// ColsQuerys with conditions with specifying columns
func (objx *RuntimeProblemer) ColsQuerys(objs []RuntimeProblem, cols ...string) error {
	return x.Cols(cols...).Find(objs)
}

// Where provIDes custom query condition.
func (objx *RuntimeProblemer) Where(query interface{}, args ...interface{}) *RuntimeProblemSession {
	return (*RuntimeProblemSession)(x.Where(query, args...))
}

// Where provIDes custom query condition.
func (objx *RuntimeProblemSession) Where(query interface{}, args ...interface{}) *RuntimeProblemSession {
	return (*RuntimeProblemSession)(((*xorm.Session)(objx)).Where(query, args...))
}

// And provIDes custom query condition.
func (objx *RuntimeProblemSession) And(query interface{}, args ...interface{}) *RuntimeProblemSession {
	return (*RuntimeProblemSession)(((*xorm.Session)(objx)).And(query, args...))
}

// Or provIDes custom query condition.
func (objx *RuntimeProblemSession) Or(query interface{}, args ...interface{}) *RuntimeProblemSession {
	return (*RuntimeProblemSession)(((*xorm.Session)(objx)).Or(query, args...))
}

// ID provIDes custom query condition.
func (objx *RuntimeProblemSession) ID(query interface{}) *RuntimeProblemSession {
	return (*RuntimeProblemSession)(((*xorm.Session)(objx)).ID(query))
}

// NotIn provIDes custom query condition.
func (objx *RuntimeProblemSession) NotIn(query string, args ...interface{}) *RuntimeProblemSession {
	return (*RuntimeProblemSession)(((*xorm.Session)(objx)).NotIn(query, args...))
}

// In provIDes custom query condition.
func (objx *RuntimeProblemSession) In(query string, args ...interface{}) *RuntimeProblemSession {
	return (*RuntimeProblemSession)(((*xorm.Session)(objx)).In(query, args...))
}

// Find retrieve records from table, condiBeans's non-empty fields
// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
// map[int64]*Struct
func (objx *RuntimeProblemSession) Find(conds ...interface{}) ([]RuntimeProblem, error) {
	objs := make([]RuntimeProblem, 0)
	err := ((*xorm.Session)(objx)).Find(objs, conds...)
	return objs, err
}
