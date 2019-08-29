package orm

import (
	"github.com/go-xorm/xorm"
)

// User example
type User struct {
	ID int `xorm:"not null pk autoincr 'id'"`

	Name           string `xorm:"'name'"`
	Password       []byte `xorm:"'password'"`
	Exp            int    `xorm:"'exp'"`
	SolvedProblems int    `xorm:"'solved_problems'"`
}

// TableName return the table name
func (obj *User) TableName() string {
	return "users"
}

// GetSliceWithPredict return the slice of User with reserving the space of n User
func (obj *User) GetSliceWithPredict(n int) interface{} {
	return make([]User, 0, n)
}

// GetSlice return the slice of User
func (obj *User) GetSlice() interface{} {
	return make([]User, 0)
}

// Insert into Engine
func (obj *User) Insert() (int64, error) {
	return x.Insert(obj)
}

// Delete from Engine
func (obj *User) Delete() (int64, error) {
	return x.ID(obj.ID).Delete(obj)
}

// Update to Engine
func (obj *User) Update() (int64, error) {
	return x.ID(obj.ID).Update(obj)
}

// UpdateAll to Engine
func (obj *User) UpdateAll() (int64, error) {
	return x.ID(obj.ID).AllCols().Update(obj)
}

// Query from Engine
func (obj *User) Query() (bool, error) {
	return x.Get(obj)
}

// Userer Extend the Engine operation
type Userer struct {
}

// UsererSession Extend the Engine operation
type UsererSession xorm.Session

// Query return the user with Property property
func (objx *Userer) Query(property int) (*User, error) {
	obj := new(User)
	obj.ID = property
	has, err := x.Get(obj)
	if has {
		return obj, nil
	}
	return nil, err
}

// QueryName return the user with Property property
func (objx *Userer) QueryName(property string) (*User, error) {
	obj := new(User)
	obj.Name = property
	has, err := x.Get(obj)
	if has {
		return obj, nil
	}
	return nil, err
}

// Inserts many Users
func (objx *Userer) Inserts(objs []User) (int64, error) {
	return x.Insert(objs)
}

// Querys with conditions
func (objx *Userer) Querys(objs []User, conds ...interface{}) error {
	return x.Find(&objs, conds...)
}

// ColsQuerys with conditions with specifying columns
func (objx *Userer) ColsQuerys(objs []User, cols ...string) error {
	return x.Cols(cols...).Find(&objs)
}

// Where provIDes custom query condition.
func (objx *Userer) Where(query interface{}, args ...interface{}) *UsererSession {
	return (*UsererSession)(x.Where(query, args...))
}

// Where provIDes custom query condition.
func (objx *UsererSession) Where(query interface{}, args ...interface{}) *UsererSession {
	return (*UsererSession)(((*xorm.Session)(objx)).Where(query, args...))
}

// And provIDes custom query condition.
func (objx *UsererSession) And(query interface{}, args ...interface{}) *UsererSession {
	return (*UsererSession)(((*xorm.Session)(objx)).And(query, args...))
}

// Or provIDes custom query condition.
func (objx *UsererSession) Or(query interface{}, args ...interface{}) *UsererSession {
	return (*UsererSession)(((*xorm.Session)(objx)).Or(query, args...))
}

// ID provIDes custom query condition.
func (objx *UsererSession) ID(query interface{}) *UsererSession {
	return (*UsererSession)(((*xorm.Session)(objx)).ID(query))
}

// NotIn provIDes custom query condition.
func (objx *UsererSession) NotIn(query string, args ...interface{}) *UsererSession {
	return (*UsererSession)(((*xorm.Session)(objx)).NotIn(query, args...))
}

// In provIDes custom query condition.
func (objx *UsererSession) In(query string, args ...interface{}) *UsererSession {
	return (*UsererSession)(((*xorm.Session)(objx)).In(query, args...))
}

// Find retrieve records from table, condiBeans's non-empty fields
// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
// map[int64]*Struct
func (objx *UsererSession) Find(conds ...interface{}) ([]User, error) {
	objs := make([]User, 0)
	err := ((*xorm.Session)(objx)).Find(objs, conds...)
	return objs, err
}
