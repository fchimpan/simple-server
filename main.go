package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/fchimpan/simple-server/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux, cleanup, err := NewMux(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	s := NewServer(l, mux)
	return s.Run(ctx)
}
