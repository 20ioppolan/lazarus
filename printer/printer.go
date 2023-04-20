package main

import (
	"fmt"
	"time"
)

func main() {
	i := 1
	for i <= 1000 {
		fmt.Println(i)
		i = i + 1
		time.Sleep(time.Millisecond * 1000)
	}
}
