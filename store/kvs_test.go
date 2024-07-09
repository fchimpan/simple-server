package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fchimpan/simple-server/entity"
	"github.com/fchimpan/simple-server/testutil"
)

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	c := testutil.OpenRedisForTest(t)
	s := &KVS{Cli: c}

	t.Run("save", func(t *testing.T) {
		t.Parallel()

		key := "test-key"
		uid := entity.UserID(1)
		ctx := context.Background()
		t.Cleanup(func() {
			c.Del(ctx, key)
		})

		if err := s.Save(ctx, key, uid); err != nil {
			t.Fatalf("failed to load: %v", err)
		}
	})
}

func TestKVS_Load(t *testing.T) {
	t.Parallel()

	c := testutil.OpenRedisForTest(t)
	s := &KVS{Cli: c}

	t.Run("load", func(t *testing.T) {
		t.Parallel()

		key := "load-key"
		uid := entity.UserID(2)
		ctx := context.Background()

		c.Set(ctx, key, int64(uid), 30*time.Minute)
		t.Cleanup(func() {
			c.Del(ctx, key)
		})

		got, err := s.Load(ctx, key)
		if err != nil {
			t.Fatalf("failed to load: %v", err)
		}
		if got != uid {
			t.Errorf("want %d, got %d", uid, got)
		}

	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		key := "load-key-not-found"
		ctx := context.Background()

		got, err := s.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Fatalf("want error, got nil")
		}
		if got != 0 {
			t.Errorf("want 0, got %d", got)
		}
	})
}
