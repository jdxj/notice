package mock

func New() *Watcher {
	return &Watcher{}
}

type Watcher struct {
}

func (w *Watcher) Name() string {
	return "mock"
}

func (w *Watcher) Watch() (string, bool, error) {
	return "hah", true, nil
}
