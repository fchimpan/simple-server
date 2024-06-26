package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})

	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("url: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}
	defer resp.Body.Close()
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}
	want := "Hello, /" + in
	if string(got) != want {
		t.Errorf("want %q, got %q", want, got)
	}

	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatalf("failed to wait: %v", err)
	}

}
