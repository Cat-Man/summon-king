package alliance

import (
	"context"
	"testing"
)

func TestCreateAlliance_RequiresLevel30(t *testing.T) {
	svc := newTestAllianceService(t)
	err := svc.CreateAlliance(context.Background(), 1001, "测试联盟")
	if err == nil {
		t.Fatal("expected low level create alliance to fail")
	}
}

func TestClaimFireTraining_RequiresFinished(t *testing.T) {
	svc := newTestAllianceService(t)
	ctx := context.Background()
	allianceID, err := svc.CreateAllianceAsGM(ctx, 3001, "烈焰盟")
	if err != nil {
		t.Fatalf("expected create alliance success, got %v", err)
	}
	record, err := svc.StartFireTraining(ctx, allianceID, 3001, "room-1")
	if err != nil {
		t.Fatalf("expected start fire training success, got %v", err)
	}
	if _, err := svc.ClaimFireTraining(ctx, allianceID, 3001, record.RecordID); err == nil {
		t.Fatal("expected unfinished fire training claim to fail")
	}
}

func newTestAllianceService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewMemoryRepository())
}
