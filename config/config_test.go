package config

import "testing"

const (
	testKey   = "hello"
	testValue = "world"
)

func TestOpenMultiDBFiles(t *testing.T) {
	c1, err := NewCache(CachePath)
	if err != nil {
		t.Fatalf("c1: %s\n", err)
	}
	defer c1.Close()

	c2, err := NewCache(CachePath)
	if err != nil {
		t.Fatalf("c2: %s\n", err)
	}
	defer c2.Close()
}
