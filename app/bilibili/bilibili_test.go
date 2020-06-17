package modules

import (
	"fmt"
	"math"
	"testing"

	"github.com/jdxj/notice/config"
)

func TestBiliBili_CollectCoins(t *testing.T) {
	cookie := ""
	emailCfg := &config.Email{
		Addr:  "985759262@qq.com",
		Token: "",
	}
	bili := NewBiliBili("985759262@qq.com", cookie, emailCfg)
	bili.CollectCoins()
}

func TestMathPow(t *testing.T) {
	fmt.Printf("%f\n", math.Pow(2, 0))
	fmt.Printf("%f\n", math.Pow(2, 5))
}
