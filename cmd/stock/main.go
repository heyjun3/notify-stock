package main

import (
	"log"
	"net/http"
	"time"

	notify "github.com/heyjun3/notify-stock"
)

func main() {
	register := notify.InitStockRegister(notify.Cfg.DBDSN, &http.Client{})
	t := time.Unix(0, 0)
	times := []time.Time{t}
	for {
		t = t.AddDate(5, 0, 0)
		if t.After(time.Now().UTC()) {
			t = time.Now().UTC()
			times = append(times, t)
			break
		}
		times = append(times, t)
	}
	for i := 0; i < len(times)-1; i++ {
		if err := register.SaveStock(
			notify.N225, times[i], times[i+1]); err != nil {
			log.Fatal(err)
		}
		time.Sleep(2 * time.Second)
	}
}
