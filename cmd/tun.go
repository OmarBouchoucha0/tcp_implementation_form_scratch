package main

import (
	"github.com/songgao/water"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	done := make(chan bool)
	go keyBoardInterrupt(done, sigChan)
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Print("Coudnt make tunnel")
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())
	packetChan := make(chan []byte)

	for {
		go func() {
			err := readBytes(packetChan, ifce)
			if err != nil {
				log.Println("Error reading bytes:", err)
			}
		}()
		select {
		case <-done:
			log.Println("Exiting...")
			err := ifce.Close()
			if err != nil {
				log.Printf("coudnt close tunnel correctly!")
			}
			os.Exit(0)
		case packet := <-packetChan:
			packetPrint(packet)
		}
	}
}
