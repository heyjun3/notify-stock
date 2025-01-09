package main

import (
	"log"
	"net/http"
	"time"

	notify "github.com/heyjun3/notify-stock"
)

func main() {
	client := notify.NewFinanceClient(&http.Client{})
	repo := notify.NewStockRepository(notify.DB)
	register := notify.NewStockRegister(client, repo)
	if err := register.SaveStock(
		notify.N225, time.Now().AddDate(-1, 0, 0), time.Now()); err != nil {
		log.Fatal(err)
	}
}
