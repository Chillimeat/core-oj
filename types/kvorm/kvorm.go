package korm

var engine KVDB

func RegisterEngine(db KVDB) {
	engine = db
}
