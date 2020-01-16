package medeming

import (
	"fmt"
	"notice/module"
	"testing"
)

func readConfig() *module.Config {
	config, err := module.ReadConfig()
	if err != nil {
		panic(err)
	}
	return config
}

func TestNewJetbrains(t *testing.T) {
	config := readConfig()

	jet, err := NewJetbrains(config.Email)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	code, err := jet.getActivationCode()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Println(code)
}
