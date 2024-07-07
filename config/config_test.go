package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	wantPort := 1234
	wantEnv := "dev"
	t.Setenv("PORT", fmt.Sprint(wantPort))
	t.Setenv("ENV", wantEnv)

	got, err := New()
	if err != nil {
		t.Fatalf("failed to parse env: %v", err)
	}
	if got.Port != wantPort {
		t.Errorf("want %v, got %v", wantPort, got.Port)
	}
	if got.Env != wantEnv {
		t.Errorf("want %v, got %v", wantEnv, got.Env)
	}

}
