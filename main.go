package main

import (
	"github.com/yokaracho/cliCat/cmd/root"
	"os"
	"os/signal"
)

func main() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		
	}()

	root.Execute()
}
