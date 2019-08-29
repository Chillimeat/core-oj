package korm

import (
	"bytes"
	"fmt"
	"testing"
)

var _ KVDB = new(GobDB)

func TestGobLevelDB(t *testing.T) {
	const notFound = "leveldb: not found"

	db, err := NewGobLevelDB("./test.db")
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()
	var value = []byte("value")
	err = db.Put([]byte("key"), value)
	if err != nil {
		t.Error(err)
		return
	}
	tvalue, err := db.Get([]byte("key"))
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(tvalue, value) {
		t.Error("not eqal")
		return
	}

	err = db.Delete([]byte("key"))
	if err != nil {
		t.Error(err)
		return
	}
	_, err = db.Get([]byte("key"))
	if err.Error() != notFound {
		t.Error(err)
		return
	}

	subdb, err := db.TableString("mytable")
	if err != nil {
		t.Error(err)
		return
	}
	subdbx := subdb.(*GobDB)
	fmt.Println(string(subdbx.slotbuf.Bytes()))
	subdb, err = db.Table([]byte("mysubtable"))
	if err != nil {
		t.Error(err)
		return
	}
	subdbx = subdb.(*GobDB)
	fmt.Println(string(subdbx.slotbuf.Bytes()))
	subses, err := db.ID('1')
	if err != nil {
		t.Error(err)
		return
	}
	subsesx := subses.(*GobDB)
	if subsesx == subdbx || subsesx == db || subdbx != db {
		t.Error("split failed")
		return
	}
	pre := subsesx

	fmt.Println(string(subsesx.slotbuf.Bytes()))
	subses, err = subses.ID64('1')
	if err != nil {
		t.Error(err)
		return
	}
	subsesx = subses.(*GobDB)
	fmt.Println(string(subsesx.slotbuf.Bytes()))
	if pre == subsesx || subsesx == subdbx || subsesx == db || subdbx != db {
		t.Error("split failed")
		return
	}
	subses, err = subses.GUID([]byte("233QAQ"))
	if err != nil {
		t.Error(err)
		return
	}
	subsesx = subses.(*GobDB)
	fmt.Println(string(subsesx.slotbuf.Bytes()))
	subses, err = subses.Name("233QAQQAQ")
	if err != nil {
		t.Error(err)
		return
	}
	subsesx = subses.(*GobDB)
	fmt.Println(string(subsesx.slotbuf.Bytes()))
	subses, err = subses.ComposeKey(struct {
		A int64
		B int64
		C [5]byte
	}{'1', '2', [5]byte{'3', 'Q', 'A', 'Q'}})
	if err != nil {
		t.Error(err)
		return
	}
	subsesx = subses.(*GobDB)
	fmt.Println(string(subsesx.slotbuf.Bytes()))
	err = subses.Push(value)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(subsesx.slotbuf.Bytes()))
	tvalue = nil
	err = subses.Extract(&tvalue)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(tvalue, value) {
		t.Error("not eqal")
		return
	}
	fmt.Println(string(subsesx.slotbuf.Bytes()))

	err = subses.Clear()
	if err != nil {
		t.Error(err)
		return
	}
	err = subses.Extract(&tvalue)
	if err.Error() != notFound {
		t.Error(err)
		return
	}
	fmt.Println(string(subsesx.slotbuf.Bytes()))
	subses.Release()
	if subsesx.slotbuf != nil {
		t.Error("not release")
		return
	}
}

//
// Extract(interface{}) error
// Push(interface{}) error
// Clear() error
