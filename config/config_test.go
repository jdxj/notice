package config

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	Init()
	fmt.Printf("%+v\n", defaultConfig)
}
