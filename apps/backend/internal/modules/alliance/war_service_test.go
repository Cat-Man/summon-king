package alliance

import (
	"context"
	"errors"
	"testing"
)

func TestWarService_LeaderCanRegisterTarget(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	allianceSvc := NewService(repo)
	warSvc := NewWarService(repo)

	if _, err := allianceSvc.CreateAlliance(ctx, 5101, "赤焰盟"); err != nil {
		t.Fatalf("expected create alliance success, got %v", err)
	}

	index, err := warSvc.GetIndex(ctx, 5101)
	if err != nil {
		t.Fatalf("expected war index success, got %v", err)
	}
	if !index.CanRegister {
		t.Fatal("expected leader can register target")
	}
	if len(index.AvailableTargets) == 0 {
		t.Fatal("expected available war targets")
	}

	updated, err := warSvc.RegisterTarget(ctx, 5101, "赤焰谷")
	if err != nil {
		t.Fatalf("expected register target success, got %v", err)
	}
	if updated.TargetLabel != "赤焰谷" {
		t.Fatalf("expected target 赤焰谷, got %s", updated.TargetLabel)
	}
}

func TestWarService_MemberCannotRegisterTarget(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()
	allianceSvc := NewService(repo)
	warSvc := NewWarService(repo)

	created, err := allianceSvc.CreateAlliance(ctx, 5201, "寒月盟")
	if err != nil {
		t.Fatalf("expected create alliance success, got %v", err)
	}
	if err := allianceSvc.Apply(ctx, 5202, created.Alliance.AllianceID); err != nil {
		t.Fatalf("expected apply success, got %v", err)
	}
	if _, err := allianceSvc.ApproveApplication(ctx, 5201, 5202); err != nil {
		t.Fatalf("expected approve success, got %v", err)
	}

	_, err = warSvc.RegisterTarget(ctx, 5202, "寒月岭")
	if !errors.Is(err, ErrAlliancePermissionDenied) {
		t.Fatalf("expected ErrAlliancePermissionDenied, got %v", err)
	}
}
