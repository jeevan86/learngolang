package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitForSig() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	// signal.Notify(c)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()
	<-done
}
