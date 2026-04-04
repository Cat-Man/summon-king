package tower

import (
	"context"
	"testing"
)

func TestService_TracksProgressPerTower(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())
	playerID := int64(1001)

	pagodaResult, err := svc.StartChallenge(ctx, playerID, "pagoda")
	if err != nil {
		t.Fatalf("expected pagoda challenge success, got %v", err)
	}
	spiritResult, err := svc.StartChallenge(ctx, playerID, "spirit")
	if err != nil {
		t.Fatalf("expected spirit challenge success, got %v", err)
	}

	if pagodaResult.Floor != 1 {
		t.Fatalf("expected pagoda first floor 1, got %d", pagodaResult.Floor)
	}
	if spiritResult.Floor != 1 {
		t.Fatalf("expected spirit first floor 1, got %d", spiritResult.Floor)
	}
}
