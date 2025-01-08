package main

import (
	"log"
	"time"

	notify "github.com/heyjun3/notify-stock"
)

func main() {
	if err := notify.SaveStock(
		notify.N225, time.Now().AddDate(-1, 0, 0), time.Now()); err != nil {
		log.Fatal(err)
	}
}
