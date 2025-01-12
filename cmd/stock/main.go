package main

import (
	"log"
	"net/http"
	"time"

	notify "github.com/heyjun3/notify-stock"
)

func main() {
	register := notify.InitStockRegister(notify.Cfg.DBDSN, &http.Client{})
	if err := register.SaveStock(
		notify.N225, time.Now().AddDate(-1, 0, 0), time.Now()); err != nil {
		log.Fatal(err)
	}
}
