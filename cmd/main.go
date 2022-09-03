package main

import (
	"fmt"
	"time"
)

func main() {
	go run()
	select {} // block forever
}

func run() {
	for {
		fmt.Printf("time is%v+\n", time.Now())
		time.Sleep(time.Second)
	}
}
