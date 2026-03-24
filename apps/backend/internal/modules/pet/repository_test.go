package pet

import (
	"context"
	"sync"
	"testing"
)

func TestGetPlayerPets_ConcurrentInitialization(t *testing.T) {
	repo := NewMemoryRepository()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = repo.GetPlayerPets(context.Background(), 1001)
		}()
	}
	wg.Wait()
}
