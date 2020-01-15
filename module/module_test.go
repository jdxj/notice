package module

import (
	"testing"

	"github.com/astaxie/beego/logs"
)

func TestReadConfig(t *testing.T) {
	config, err := ReadConfig()
	if err != nil {
		t.Fatalf("%s", err)
	}
	logs.Info("%+v", *config.NeoProxy)
}

func getConfig() *Config {
	config, err := ReadConfig()
	if err != nil {
		panic(err)
	}
	return config
}

func TestEmailSender_SendMsg(t *testing.T) {
	ec := getConfig().Email

	es, err := NewEmailSender(ec)
	if err != nil {
		t.Fatalf("%s", err)
	}

	es.SendMsgString("test", "hello")
}

func TestWriteConfig(t *testing.T) {
	config, err := ReadConfig()
	if err != nil {
		t.Fatalf("%s", err)
	}
	err = WriteConfig(config)
	if err != nil {
		t.Fatalf("%s", err)
	}
}
