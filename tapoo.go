package main

import (
	"fmt"
	"time"
)

func main() {
	var (
		m = 3
	)

	for i := 1; i < 20; i++ {
		fmt.Printf("\t# Random No. %+v \n", time.Now().UnixNano()%int64(m))
	}

}
