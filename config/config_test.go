package config

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	fmt.Printf("%+v, %+v\n", cw.TelegramBot, cw.Github)
	fmt.Printf("%+v, %+v\n", TelegramBot, Github)
}
