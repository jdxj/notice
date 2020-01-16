package medeming

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"notice/module"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

// http://idea.medeming.com/jetbrains/

const (
	ActivationCodePage = "http://idea.medeming.com/jetbrains/images/jihuoma.txt"
)

func NewJetbrains(ec *module.EmailConfig) (*Jetbrains, error) {
	if ec == nil {
		return nil, fmt.Errorf("invalid EmailSender config")
	}

	es, err := module.NewEmailSender(ec)
	if err != nil {
		return nil, err
	}

	jet := &Jetbrains{
		url:         "http://idea.medeming.com/jetbrains/images/jihuoma.txt",
		emailSender: es,
		wg:          &sync.WaitGroup{},
		client:      &http.Client{},
		stop:        make(chan struct{}),
	}
	return jet, nil
}

type Jetbrains struct {
	url         string
	emailSender *module.EmailSender

	wg     *sync.WaitGroup
	client *http.Client

	stop chan struct{}
}

func (jet *Jetbrains) Start() {
	wg := jet.wg
	wg.Add(1)
	go func() {
		jet.NotifyFoundNewActivationCode()
		wg.Done()
	}()

	logs.Info("jetbrains started")
}

func (jet *Jetbrains) Stop() {
	close(jet.stop)
	jet.wg.Wait()

	logs.Info("jetbrains stopped")
}

func (jet *Jetbrains) NotifyFoundNewActivationCode() {
	var lastActivationCode string

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-jet.stop:
			logs.Info("stop notify found new activation code")
			return

		case <-ticker.C:
			ac, err := jet.getActivationCode()
			if err != nil {
				logs.Error("%s", err)
				continue
			}
			if ac == lastActivationCode {
				continue
			}
			lastActivationCode = ac

			subject := "发现新的 jetbrains 激活码"
			if err := jet.emailSender.SendMsgString(subject, ac); err != nil {
				logs.Error("%s", err)
			}
		}
	}
}

func (jet *Jetbrains) getActivationCode() (string, error) {
	resp, err := jet.client.Get(jet.url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", data), nil
}
