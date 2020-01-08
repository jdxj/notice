package neoproxy

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	str := time.Now().Add(time.Hour * 24 * 30).Format("Jan 02")
	fmt.Println(str)
}

func TestWaitGroup(t *testing.T) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		time.Sleep(5 * time.Second)
		wg.Done()
	}()

	//wg.Wait()
}

func TestMemorySize(t *testing.T) {

}
