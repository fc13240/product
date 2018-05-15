package redisCli

import (
	"fmt"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	n := 0
	for range [1000]int{} {
		go func(n int) {
			n = n + 1
			Get("a")
		}(n)
	}
	time.Sleep(time.Second * 3)
	fmt.Println(n)

}
