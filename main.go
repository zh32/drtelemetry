package main

import (
	"drtelemetry/telemetry"
	"os"
	"os/signal"
	"fmt"
	"time"
	"flag"
	"log"
	"drtelemetry/ui"
	"syscall"
)

func main() {
	fmt.Println("**********************************************")
	fmt.Println("**** DiRT Rally Telemetry Overlay by zh32 ****")
	fmt.Println("**********************************************")

	flag.Parse()
	log.SetFlags(0)

	dataChannel, quit := telemetry.RunServer(":10001")
	go ui.ListenAndServe(dataChannel)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c

	fmt.Printf("captured %v, stopping profiler and exiting..", sig)
	close(quit)
	time.Sleep(2 * time.Second)
	os.Exit(1)
}
