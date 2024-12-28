package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("now local", time.Unix(time.Now().Unix(), 0))
	fmt.Println("now utc", time.Unix(time.Now().Unix(), 0).UTC())
}
