package orm

import (
	"sort"
	"sync"
	"sync/atomic"

	"github.com/go-xorm/xorm"
)

const (
	StatusWaitingForJudge int = iota
	StatusAccepted
	StatusRunning
	StatusCompiling
	StatusCompileError
	StatusCompileTimeout
	StatusWrongAnswer
	StatusTimeLimitExceed
	StatusMemoryLimitExceed
	StatusSystemError
	StatusUnknownError
	StatusPresentationError
	StatusRuntimeError
)

// Code records the code in online judge
type Code struct {
	ID        int    `xorm:"not null pk autoincr"`
	CodeType  string `xorm:"'code_type'"`
	Hash      []byte `xorm:"'hash'"`
	OwnedUID  int    `xorm:"'owner_uid'"`
	ProblemID int    `xorm:"'problem_id'"`
	Status    int    `xorm:"'status'"`
}

type CodeSlice []Code

func (codeSlice CodeSlice) Len() int {
	return len(codeSlice)
}

func (codeSlice CodeSlice) Swap(i, j int) {
	codeSlice[i], codeSlice[j] = codeSlice[j], codeSlice[i]
}

func (codeSlice CodeSlice) Less(i, j int) bool {
	return codeSlice[i].ID < codeSlice[j].ID
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
	return x.ID(obj.ID).Delete(obj)
}

// Update to Engine
func (obj *Code) Update() (int64, error) {
	return x.ID(obj.ID).Update(obj)
}

// Query from Engine
func (obj *Code) Query() (bool, error) {
	return x.Get(obj)
}

type node struct {
	v  int
	ls *node
	nx *node
}

type list struct {
	begin *node
	end   *node
}

func newList() (l *list) {
	l = new(list)
	l.begin = new(node)
	l.end = new(node)
	l.begin.nx = l.end
	l.end.ls = l.begin
	l.begin.v = 0
}

func (l *list) Len() int {
	return l.begin.v
}

func (l *list) PushFront(v int) {
	n := new(node)
	n.v = v
	n.nx = l.begin.nx
	n.ls = l.begin
	l.begin.nx.ls = n
	l.begin.nx = n
}

func (l *list) PopBack() (v int) {
	if l.end.ls == l.begin {
		panic("nil list")
	}
	v = l.end.ls.v
	l.end.ls.ls.nx = l.end
	l.end.ls = l.ls.ls
	return
}

// Coder Extend the Engine operation
type Coder struct {
	aliveID    int64
	lastID     int64
	mutex      sync.Mutex
	aliveCodes map[int]*Code

	head list

	waitingCodes   chan int
	CompilingCodes chan *Code
	RunningCodes   chan *Code
}

func min(l, r int64) int64 {
	if l < r {
		return l
	}
	return r
}

func max(l, r int64) int64 {
	if l > r {
		return l
	}
	return r
}

func NewCoder() (cr *Coder, err error) {
	cr = new(Coder)
	cr.waitingCodes = make(chan int, 10000)
	cr.CompilingCodes = make(chan *Code, 10000)
	cr.RunningCodes = make(chan *Code, 10000)

	cr.aliveID = 0x7fffffffffffffff

	var cond = &Code{Status: StatusCompiling}
	codes, err := cr.Find(cond)
	if err != nil {
		return nil, err
	}
	if len(codes) != 0 {

		sort.Sort(CodeSlice(codes))
		cr.aliveID = codes[0].ID
		cr.lastID = codes[len(codes)-1].ID
	}

	for _, code := range codes {
		cr.aliveCodes[code.ID] = &code
		cr.CompilingCodes <- &code
	}

	cond.Status = StatusRunning
	codes, err = cr.Find(cond)
	if err != nil {
		return nil, err
	}

	if len(codes) != 0 {

		sort.Sort(CodeSlice(codes))
		cr.aliveID = min(cr.aliveID, codes[0].ID)
		cr.lastID = max(cr.lastID, codes[len(codes)-1].ID)
	} else {
		if cr.aliveID == 0x7fffffffffffffff {
			cr.aliveID = 0
		}
	}

	for _, code := range codes {
		cr.aliveCodes[code.ID] = &code
		cr.RunningCodes <- &code
	}

	cond.Status = StatusWaitingForJudge
	codes, err = cr.Find(cond)
	if err != nil {
		return nil, err
	}
	if len(codes) != 0 {
		sort.Sort(CodeSlice(codes))
	}

	for _, code := range codes {
		cr.waitingCodes <- code.ID
	}

	return cr, nil

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

// QueryStatus return the code with Property property
func (objx *Coder) QueryStatus(property int) (*Code, error) {
	obj := new(Code)
	obj.Status = property
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

// Where provIDes custom query condition.
func (objx *Coder) Where(query interface{}, args ...interface{}) *CoderSession {
	return (*CoderSession)(x.Where(query, args...))
}

// Find retrieve records from table, condiBeans's non-empty fields
// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
// map[int64]*Struct
func (objx *Coder) Find(conds ...interface{}) ([]Code, error) {
	objs := make([]Code, 0)
	err := x.Find(objs, conds...)
	return objs, err
}

// Where provIDes custom query condition.
func (objx *CoderSession) Where(query interface{}, args ...interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).Where(query, args...))
}

// And provIDes custom query condition.
func (objx *CoderSession) And(query interface{}, args ...interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).And(query, args...))
}

// Or provIDes custom query condition.
func (objx *CoderSession) Or(query interface{}, args ...interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).Or(query, args...))
}

// ID provIDes custom query condition.
func (objx *CoderSession) ID(query interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).ID(query))
}

// NotIn provIDes custom query condition.
func (objx *CoderSession) NotIn(query string, args ...interface{}) *CoderSession {
	return (*CoderSession)(((*xorm.Session)(objx)).NotIn(query, args...))
}

// In provIDes custom query condition.
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

func (objx *Coder) PushTask(id int) error {
	objx.waitingCodes <- id
	//todo check task queue
	return nil
}

func (objx *Coder) StartToExecuteTask() (*Code, bool, error) {
	if id, ok := <-objx.waitingCodes; ok {
		code, err := objx.Query(id)
		if err != nil {
			objx.waitingCodes <- id
			return nil, false, err
		}

		objx.mutex.Lock()
		objx.aliveCodes[id] = code
		objx.mutex.Unlock()

		atomic.StoreInt64(&objx.lastID, int64(id))
		return code, true, nil
	}
	return nil, false, nil
}

func (objx *Coder) SettleTask(id int) bool {

	atomic.StoreInt64(&objx.lastID, int64(id))

	objx.mutex.Lock()
	delete(objx.aliveCodes, id)
	objx.mutex.Unlock()
	return true
}
