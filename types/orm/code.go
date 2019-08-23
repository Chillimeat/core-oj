package orm

import "github.com/go-xorm/xorm"

// Code records the code in online judge
type Code struct {
	ID        int `xorm:"not null pk autoincr"`
	Hash      []byte
	OwnedUID  int
	ProblemID int
	RuntimeID int
}

// TableName return the table name
func (obj *Code) TableName() string {
	return "codes_records"
}

// GetSliceWithPredict return the slice of Code with reserving the space of n Code
func (obj *Code) GetSliceWithPredict(n int) interface{} {
	return make([]Code, 0, n)
}

// GetSlice return the slice of Code
func (obj *Code) GetSlice() interface{} {
	return new([]Code)
}

// Insert into Engine
func (obj *Code) Insert() (int64, error) {
	return x.Insert(obj)
}

// Delete from Engine
func (obj *Code) Delete() (int64, error) {
	return x.Id(obj.ID).Delete(obj)
}

// Update to Engine
func (obj *Code) Update() (int64, error) {
	return x.Id(obj.ID).Update(obj)
}

// Query from Engine
func (obj *Code) Query() (bool, error) {
	return x.Get(obj)
}

// Coder Extend the Engine operation
type Coder struct {
}

// CoderSession Extend the Engine operation
type CoderSession xorm.Session

// Query return the code with Property property
func (objx *Coder) Query(property int) (*Code, error) {
	obj := new(Code)
	obj.ID = property
	has, err := x.Get(obj)
	if has {
		return obj, nil
	}
	return nil, err
}

// QueryHash return the code with Property property
func (objx *Coder) QueryHash(property []byte) (*Code, error) {
	obj := new(Code)
	obj.Hash = property
	has, err := x.Get(obj)
	if has {
		return obj, nil
	}
	return nil, err
}

// QueryOwnedUID return the code with Property property
func (objx *Coder) QueryOwnedUID(property int) (*Code, error) {
	obj := new(Code)
	obj.OwnedUID = property
	has, err := x.Get(obj)
	if has {
		return obj, nil
	}
	return nil, err
}

// QueryProblemID return the code with Property property
func (objx *Coder) QueryProblemID(property int) (*Code, error) {
	obj := new(Code)
	obj.ProblemID = property
	has, err := x.Get(obj)
	if has {
		return obj, nil
	}
	return nil, err
}

// QueryRuntimeID return the code with Property property
func (objx *Coder) QueryRuntimeID(property int) (*Code, error) {
	obj := new(Code)
	obj.RuntimeID = property
	has, err := x.Get(obj)
	if has {
		return obj, nil
	}
	return nil, err
}

// Inserts many objects
func (objx *Coder) Inserts(objs []Code) (int64, error) {
	return x.Insert(objs)
}

// Querys with conditions
func (objx *Coder) Querys(objs []Code, conds ...interface{}) error {
	return x.Find(objs, conds...)
}

// ColsQuerys with conditions with specifying columns
func (objx *Coder) ColsQuerys(objs []Code, cols ...string) error {
	return x.Cols(cols...).Find(objs)
}

// Where provides custom query condition.
func (objx *Coder) Where(query interface{}, args ...interface{}) *CoderSession {
	return (*CoderSession)(x.Where(query, args...))
}

// Where provides custom query condition.
func (objx *CoderSession) Where(query interface{}, args ...interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).Where(query, args...))
}

// And provides custom query condition.
func (objx *CoderSession) And(query interface{}, args ...interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).And(query, args...))
}

// Or provides custom query condition.
func (objx *CoderSession) Or(query interface{}, args ...interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).Or(query, args...))
}

// ID provides custom query condition.
func (objx *CoderSession) ID(query interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).ID(query))
}

// NotIn provides custom query condition.
func (objx *CoderSession) NotIn(query string, args ...interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).NotIn(query, args...))
}

// In provides custom query condition.
func (objx *CoderSession) In(query string, args ...interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).In(query, args...))
}

// Find retrieve records from table, condiBeans's non-empty fields
// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
// map[int64]*Struct
func (objx *CoderSession) Find(conds ...interface{}) ([]Code, error) {
	objs := make([]Code, 0)
	err := ((*xorm.Session)(objx)).Find(objs, conds...)
	return objs, err
}
