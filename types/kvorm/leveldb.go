package korm

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"

	"github.com/syndtr/goleveldb/leveldb"
)

type Buffer interface {
	io.ReadWriter
	Bytes() []byte
}

type BufferPool interface {
	Get() Buffer
	Put(Buffer)
}

type BaseKVDB interface {
	Get([]byte) ([]byte, error)
	Put([]byte, []byte) error
	Delete([]byte) error
}

type BaseKVDBCloser interface {
	BaseKVDB
	io.Closer
}

type BaseSession interface {
	BaseKVDB
	Extract(interface{}) error
	Push(interface{}) error
	Clear() error

	Release()

	Table([]byte) (KVDB, error)
	TableString(string) (KVDB, error)
}

type BaseJustSession interface {
	ID(int) (Session, error)
	ID64(int64) (Session, error)
	GUID([]byte) (Session, error)
	Name(string) (Session, error)
	ComposeKey(interface{}) (Session, error)
}

type BaseMaybeSession interface {
	IDThen(int) MaybeSession
	ID64Then(int64) MaybeSession
	GUIDThen([]byte) MaybeSession
	NameThen(string) MaybeSession
	ComposeKeyThen(interface{}) MaybeSession
}

type Session interface {
	BaseSession
	BaseJustSession
	BaseMaybeSession
}

type MaybeSession interface {
	error
	BaseSession
	BaseMaybeSession
}

type KVDB interface {
	Session
}

type TopKVDB interface {
	KVDB
	LinkPool(BufferPool)
	io.Closer
}

type commonSpace struct {
	BaseKVDBCloser
	bp BufferPool
}

type GobDB struct {
	slotbuf Buffer
	*commonSpace
}

type MaybeGobDB struct {
	error
	db *GobDB
}

type dbp struct{}

func (dbp *dbp) Get() Buffer {
	return new(bytes.Buffer)
}

func (dbp *dbp) Put(Buffer) {
	return
}

var defaultBufferPool = new(dbp)

type LevelKVDB struct {
	db *leveldb.DB
}

func NewLevelKVDB(db *leveldb.DB) *LevelKVDB {
	return &LevelKVDB{db}
}

func (objx *LevelKVDB) Put(k, v []byte) error {
	return objx.db.Put(k, v, nil)
}
func (objx *LevelKVDB) Delete(b []byte) error {
	return objx.db.Delete(b, nil)
}
func (objx *LevelKVDB) Get(b []byte) ([]byte, error) {
	return objx.db.Get(b, nil)
}
func (objx *LevelKVDB) Close() error {
	return objx.db.Close()
}

func NewGobDB(db *leveldb.DB) *GobDB {
	return &GobDB{
		slotbuf: defaultBufferPool.Get(),
		commonSpace: &commonSpace{
			BaseKVDBCloser: NewLevelKVDB(db),
			bp:             defaultBufferPool,
		},
	}
}

func NewGobLevelDB(path string) (*GobDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return NewGobDB(db), nil
}

func (objx *GobDB) Clone(src Buffer) Buffer {
	var buf = objx.commonSpace.bp.Get()
	_, err := buf.Write(src.Bytes())
	if err != nil {
		panic("can not clone...")
	}

	return buf
}

func (objx *GobDB) Grow(values ...interface{}) error {
	for _, value := range values {
		if err := binary.Write(objx.slotbuf, binary.BigEndian, value); err != nil {
			return err
		}
	}
	return nil
}

func (objx *GobDB) LinkPool(bp BufferPool) {
	objx.commonSpace.bp = bp
}

var sep = ' '
var septable = '$'

func (objx *GobDB) Table(id []byte) (KVDB, error) {
	return objx, objx.Grow(septable, id)
}

func (objx *GobDB) TableString(id string) (KVDB, error) {
	return objx, objx.Grow(septable, []byte(id))
}

func (objx *GobDB) Name(id string) (Session, error) {
	newobjx := &GobDB{
		commonSpace: objx.commonSpace,
		slotbuf:     objx.Clone(objx.slotbuf),
	}
	return newobjx, newobjx.Grow(sep, []byte(id))
}

func (objx *GobDB) NameThen(id string) MaybeSession {
	newobjx := &GobDB{
		commonSpace: objx.commonSpace,
		slotbuf:     objx.Clone(objx.slotbuf),
	}
	return &MaybeGobDB{
		error: newobjx.Grow(sep, id),
		db:    newobjx,
	}
}

func (objx *GobDB) ID64(id int64) (Session, error) {
	newobjx := &GobDB{
		commonSpace: objx.commonSpace,
		slotbuf:     objx.Clone(objx.slotbuf),
	}
	return newobjx, newobjx.Grow(sep, id)
}

func (objx *GobDB) ID64Then(id int64) MaybeSession {
	newobjx := &GobDB{
		commonSpace: objx.commonSpace,
		slotbuf:     objx.Clone(objx.slotbuf),
	}
	return &MaybeGobDB{
		error: newobjx.Grow(sep, id),
		db:    newobjx,
	}
}

func (objx *GobDB) GUID(id []byte) (Session, error) {
	newobjx := &GobDB{
		commonSpace: objx.commonSpace,
		slotbuf:     objx.Clone(objx.slotbuf),
	}
	return newobjx, newobjx.Grow(sep, id)
}

func (objx *GobDB) GUIDThen(id []byte) MaybeSession {
	newobjx := &GobDB{
		commonSpace: objx.commonSpace,
		slotbuf:     objx.Clone(objx.slotbuf),
	}
	return &MaybeGobDB{
		error: newobjx.Grow(sep, id),
		db:    newobjx,
	}
}

func (objx *GobDB) ComposeKey(id interface{}) (Session, error) {
	newobjx := &GobDB{
		commonSpace: objx.commonSpace,
		slotbuf:     objx.Clone(objx.slotbuf),
	}
	return newobjx, newobjx.Grow(sep, id)
}

func (objx *GobDB) ComposeKeyThen(id interface{}) MaybeSession {
	newobjx := &GobDB{
		commonSpace: objx.commonSpace,
		slotbuf:     objx.Clone(objx.slotbuf),
	}
	return &MaybeGobDB{
		error: newobjx.Grow(sep, id),
		db:    newobjx,
	}
}

func (objx *GobDB) ID(id int) (Session, error) {
	newobjx := &GobDB{
		commonSpace: objx.commonSpace,
		slotbuf:     objx.Clone(objx.slotbuf),
	}
	return newobjx, newobjx.Grow(sep, int64(id))
}

func (objx *GobDB) IDThen(id int) MaybeSession {
	newobjx := &GobDB{
		commonSpace: objx.commonSpace,
		slotbuf:     objx.Clone(objx.slotbuf),
	}
	return &MaybeGobDB{
		error: newobjx.Grow(sep, int64(id)),
		db:    newobjx,
	}
}

func (objx *GobDB) Push(v interface{}) error {
	var b = objx.commonSpace.bp.Get()
	err := gob.NewEncoder(b).Encode(v)
	if err != nil {
		return err
	}
	err = objx.Put(objx.slotbuf.Bytes(), b.Bytes())
	if err != nil {
		return err
	}
	objx.commonSpace.bp.Put(b)
	return nil
}

func (objx *GobDB) Extract(v interface{}) error {
	b, err := objx.Get(objx.slotbuf.Bytes())
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(b)).Decode(v)
}

func (objx *GobDB) Clear() error {
	return objx.Delete(objx.slotbuf.Bytes())
}

func (objx *GobDB) Release() {
	objx.commonSpace.bp.Put(objx.slotbuf)
	objx.slotbuf = nil
}
func (objx *MaybeGobDB) NameThen(id string) MaybeSession {
	if objx.error != nil {
		return objx
	}
	objx.error = objx.db.NameThen(id)
	return objx
}

func (objx *MaybeGobDB) ComposeKeyThen(id interface{}) MaybeSession {
	if objx.error != nil {
		return objx
	}
	objx.error = objx.db.ComposeKeyThen(id)
	return objx
}

func (objx *MaybeGobDB) GUIDThen(id []byte) MaybeSession {
	if objx.error != nil {
		return objx
	}
	objx.error = objx.db.GUIDThen(id)
	return objx
}

func (objx *MaybeGobDB) IDThen(id int) MaybeSession {
	if objx.error != nil {
		return objx
	}
	objx.error = objx.db.IDThen(id)
	return objx
}

func (objx *MaybeGobDB) ID64Then(id int64) MaybeSession {
	if objx.error != nil {
		return objx
	}
	objx.error = objx.db.ID64Then(id)
	return objx
}

func (objx *MaybeGobDB) Table(name []byte) (KVDB, error) {
	if objx.error != nil {
		return nil, objx
	}
	return objx.db.Table(name)
}

func (objx *MaybeGobDB) TableString(name string) (KVDB, error) {
	if objx.error != nil {
		return nil, objx
	}
	return objx.db.Table([]byte(name))
}

func (objx *MaybeGobDB) Get(v []byte) ([]byte, error) {
	if objx.error != nil {
		return nil, objx
	}
	return objx.db.Get(v)
}

func (objx *MaybeGobDB) Put(k, v []byte) error {
	if objx.error != nil {
		return objx
	}
	return objx.db.Put(k, v)
}

func (objx *MaybeGobDB) Delete(v []byte) error {
	if objx.error != nil {
		return objx
	}
	return objx.db.Delete(v)
}

func (objx *MaybeGobDB) Push(v interface{}) error {
	if objx.error != nil {
		return objx.error
	}
	return objx.db.Push(v)
}

func (objx *MaybeGobDB) Extract(v interface{}) error {
	if objx.error != nil {
		return objx.error
	}
	return objx.db.Extract(v)
}

func (objx *MaybeGobDB) Clear() error {
	if objx.error != nil {
		return objx.error
	}
	return objx.db.Clear()
}

func (objx *MaybeGobDB) Release() {
	objx.db.Release()
	return
}
