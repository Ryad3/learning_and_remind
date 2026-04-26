package main

import (
	"fmt"
)

func main() {
	for i := 1; i <= 10; i++ {
		parityCheck(i)
	}
}

func parityCheck(num int) {
	if num%2 == 0 {
		fmt.Println(num, "is even")
	} else {
		fmt.Println(num, "is odd")
	}
}