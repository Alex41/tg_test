package main

import (
	"tg_test/bot"
	"tg_test/server"
)

func main() {
	go bot.ListenBot()
	go server.Serve()

	select {}
}
