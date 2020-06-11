package config

import (
	"testing"
)

const (
	testKey   = "hello"
	testValue = "world"
)

func TestGet(t *testing.T) {
	value, err := cache.Get([]byte(testKey))
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	if string(value) != testValue {
		t.Fatalf("value not equal: %s\n", value)
	}

	cache.db.Sync()
	cache.db.Close()
}

func TestSet(t *testing.T) {
	err := cache.Set([]byte(testKey), []byte(testValue))
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	cache.db.Sync()
	cache.db.Close()
}
