package main

import (
	"fmt"
	"log"
	"time"

	notifystock "github.com/heyjun3/notify-stock"
)

func main() {
	fmt.Println("now local", time.Unix(time.Now().Unix(), 0))
	fmt.Println("now utc", time.Unix(time.Now().Unix(), 0).UTC())
	if err := notifystock.RefreshToken(); err != nil {
		log.Fatal(err)
	}
}
