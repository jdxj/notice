package config

import (
	"github.com/astaxie/beego/logs"
	"github.com/dgraph-io/badger/v2"
)

const (
	CachePath = "cache.db"
)

var (
	badgerDB *badger.DB
)

func init() {
	opt := badger.DefaultOptions(CachePath)
	opt.Logger = nil

	db, err := badger.Open(opt)
	if err != nil {
		panic(err)
	}
	badgerDB = db
}

func get(key []byte) (value []byte, err error) {
	err = badgerDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(nil)
		return err
	})
	return value, err
}

func set(key, value []byte) error {
	return badgerDB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func Close() error {
	if err := badgerDB.Sync(); err != nil {
		logs.Error("sync cache failed: %s", err)
	}
	return badgerDB.Close()
}
