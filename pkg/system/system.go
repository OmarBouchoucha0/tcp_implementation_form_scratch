package system

import (
	"fmt"
	"log"
	"os"
)

func KeyBoardInterrupt(done chan bool, sigChan chan os.Signal) {
	<-sigChan
	fmt.Println("")
	log.Println("Received interrupt signal. Cleaning up...")
	close(done)
}
