package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Asmodeus8/go-chat/server"
)

func main() {
	addr := flag.String("addr", ":8080", "TCP listen address")
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil { log.Fatal(err) }
	defer ln.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() { <-sig; _ = ln.Close() }()

	fmt.Printf("go-chat listening on %s\n", *addr)
	if err := server.New().Serve(ln); err != nil && !errors.Is(err, net.ErrClosed) { log.Fatal(err) }
}
