package secrets

import (
	"context"
	"testing"
	"time"
)

type testSecretManager struct {
	value string
	hits  int
}

func (sm *testSecretManager) GetSecretValue(ctx context.Context, name string) (string, error) {
	sm.hits++
	return sm.value, nil
}

func TestCacher_GetSecretValue(t *testing.T) {
	ctx := context.Background()

	sm := &testSecretManager{value: "first"}
	cached, err := WrapCacher(ctx, sm, 250*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	// Read the value once, which should cache it.
	if _, err := cached.GetSecretValue(ctx, "secret"); err != nil {
		t.Fatal(err)
	}

	// Change the value.
	sm.value = "second"

	// Read the value a few more times.
	for i := 0; i < 5; i++ {
		got, err := cached.GetSecretValue(ctx, "secret")
		if err != nil {
			t.Fatal(err)
		}

		if want := "first"; got != want {
			t.Errorf("expected %q to be %q", got, want)
		}

		if sm.hits > 1 {
			t.Errorf("cache miss: %d", sm.hits)
		}
	}

	// Wait for cache to expire.
	time.Sleep(300 * time.Millisecond)

	// Try again - should miss
	got, err := cached.GetSecretValue(ctx, "secret")
	if err != nil {
		t.Fatal(err)
	}

	if want := "second"; got != want {
		t.Errorf("expected %q to be %q", got, want)
	}

	if sm.hits != 2 {
		t.Errorf("expected another hit: %d", sm.hits)
	}
}
