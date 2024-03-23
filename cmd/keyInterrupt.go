package main

import (
	"log"
	"os"
)

func keyBoardInterrupt(done chan bool, sigChan chan os.Signal) {
	<-sigChan
	log.Println("\nReceived interrupt signal. Cleaning up...")
	close(done)
}
