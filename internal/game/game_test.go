package game

import (
	"testing"
)

func TestGameInitialization(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf("failed to initialize game: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil game instance")
	}

	// Verify initial bicycle position
	if g.bicycle.PosX <= 0 || g.bicycle.PosY <= 0 {
		t.Errorf("expected positive initial coordinates, got (%f, %f)", g.bicycle.PosX, g.bicycle.PosY)
	}

	// Verify starting wallet
	if g.wallet.RegularMarbles != 1 {
		t.Errorf("expected 1 starting regular marble, got %d", g.wallet.RegularMarbles)
	}

	// Run 60 ticks of Update
	for i := 0; i < 60; i++ {
		if err := g.Update(); err != nil {
			t.Fatalf("unexpected error during Update tick %d: %v", i, err)
		}
	}
}

func TestGameUpdateZeroHeapAllocations(t *testing.T) {
	g, err := NewGame()
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// Warm up
	for i := 0; i < 20; i++ {
		_ = g.Update()
	}

	allocs := testing.AllocsPerRun(50, func() {
		_ = g.Update()
	})

	if allocs != 0 {
		t.Errorf("expected 0 heap allocations per Update() tick, got %f", allocs)
	}
}

