package config

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	fmt.Printf("%+v\n", defaultConfig)
}