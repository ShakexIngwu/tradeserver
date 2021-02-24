package main

import (
	"sync"

	"github.com/ShakexIngwu/tradeserver/server"
)

func main() {
	var wg sync.WaitGroup
	server.LogInit()
	t, err := server.NewTradeServer()
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go t.Work(&wg)
	wg.Wait()
}