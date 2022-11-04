package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"
)

var (
	name = flag.String("name", "foo", "The greeting object.")
)

func main() {
	flag.Parse()

	log.Println("hello", *name)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGABRT)
	defer stop()

	<-ctx.Done()

	log.Println("stop", *name)

	time.Sleep(time.Second)
	log.Println("goodbye", *name)
}
