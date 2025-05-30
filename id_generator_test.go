package wildcat

import (
	"sync"
	"testing"
)

func TestNewIDGenerator(t *testing.T) {
	g := newIDGenerator()
	if g == nil {
		t.Fatal("NewIDGenerator returned nil")
	}
	if g.lastID != 0 {
		t.Fatal("lastID was not initialized")
	}
}

func TestNextID_Unique(t *testing.T) {
	g := newIDGenerator()
	id1 := g.nextID()
	id2 := g.nextID()

	if id1 == id2 {
		t.Fatal("NextID did not generate unique IDs")
	}
}

func TestNextID_Monotonic(t *testing.T) {
	g := newIDGenerator()
	id1 := g.nextID()
	id2 := g.nextID()

	if id2 <= id1 {
		t.Fatalf("NextID did not ensure monotonicity: id1=%d, id2=%d", id1, id2)
	}
}

func TestNextID_ThreadSafety(t *testing.T) {
	g := newIDGenerator()
	const numGoroutines = 100
	const idsPerGoroutine = 100

	var wg sync.WaitGroup
	ids := make(chan int64, numGoroutines*idsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				ids <- g.nextID()
			}
		}()
	}

	wg.Wait()
	close(ids)

	// Check for uniqueness
	idSet := make(map[int64]struct{})
	for id := range ids {
		if _, exists := idSet[id]; exists {
			t.Fatalf("Duplicate ID detected: %d", id)
		}
		idSet[id] = struct{}{}
	}
}
