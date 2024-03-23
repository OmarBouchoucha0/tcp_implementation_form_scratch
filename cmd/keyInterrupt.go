package main

import (
	"fmt"
	"log"
	"os"
)

func keyBoardInterrupt(done chan bool, sigChan chan os.Signal) {
	<-sigChan
	fmt.Println("")
	log.Println("Received interrupt signal. Cleaning up...")
	close(done)
}
