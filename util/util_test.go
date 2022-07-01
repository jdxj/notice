package util

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	m := map[string]string{}

	go func() {
		for {
			for k, v := range m {
				fmt.Printf("k: %s, v: %s\n", k, v)
				time.Sleep(time.Second)
			}
			fmt.Println("ok")
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			m[strconv.Itoa(i)] = strconv.Itoa(i)
			time.Sleep(time.Second)
		}
	}()
	time.Sleep(time.Hour)
}

func TestSplit(t *testing.T) {
	str := "/"
	res := strings.Split(str, "/")
	fmt.Printf("len: %d, res: %v\n", len(res), res)

	str = "a/"
	res = strings.Split(str, "/")
	fmt.Printf("len: %d, res: %v\n", len(res), res)

	str = ""
	res = strings.Split(str, "/")
	fmt.Printf("len: %d, res: %v\n", len(res), res)
}
