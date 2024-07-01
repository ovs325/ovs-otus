package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Connection timeout")
	flag.Parse()

	if flag.NArg() != 2 {
		log.Fatalf("Usage: %s --timeout=10s host port", os.Args[0])
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	ctx, cancel := context.WithCancel(context.Background())
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer client.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := client.Receive()
		if err != nil {
			log.Printf("Error receiving data: %v", err)
		}
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := client.Send()
		if err != nil {
			log.Printf("Error sending data: %v", err)
		}
	}()

	<-ctx.Done()
	wg.Wait()
}
