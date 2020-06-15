package config

import (
	"github.com/dgraph-io/badger/v2"
)

type Cache struct {
	db *badger.DB
}

func NewCache(path string) (*Cache, error) {
	opt := badger.DefaultOptions(path)
	opt.Logger = nil

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	cache := &Cache{
		db: db,
	}
	return cache, nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	var value []byte

	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(nil)
		return err
	})

	return value, err
}

func (c *Cache) Set(key, value []byte) error {
	if err := c.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	}); err != nil {
		return err
	}

	return c.db.Sync()
}

func (c *Cache) Close() error {
	return c.db.Close()
}
