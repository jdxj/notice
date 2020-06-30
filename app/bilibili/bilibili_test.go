package bilibili

import (
	"fmt"
	"math"
	"testing"
)

func TestBiliBili_CollectCoins(t *testing.T) {
	bili := NewBiliBili()
	bili.KeepCollect()
}

func TestMathPow(t *testing.T) {
	fmt.Printf("%f\n", math.Pow(2, 0))
	fmt.Printf("%f\n", math.Pow(2, 5))
}
