package alliance

import (
	"context"
	"errors"
	"testing"
)

func TestService_CreateAllianceBuildsLeaderIndex(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	index, err := svc.CreateAlliance(ctx, 3001, "青云盟")
	if err != nil {
		t.Fatalf("expected create alliance success, got %v", err)
	}

	if !index.HasAlliance {
		t.Fatal("expected created player to join alliance")
	}
	if index.Alliance == nil {
		t.Fatal("expected alliance summary")
	}
	if index.Alliance.Name != "青云盟" {
		t.Fatalf("expected alliance name 青云盟, got %s", index.Alliance.Name)
	}
	if index.CurrentRole != "leader" {
		t.Fatalf("expected role leader, got %s", index.CurrentRole)
	}
	if len(index.Alliance.Members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(index.Alliance.Members))
	}
	if len(index.Alliance.Buildings) != 2 {
		t.Fatalf("expected 2 default buildings, got %d", len(index.Alliance.Buildings))
	}
	if index.FireTraining == nil {
		t.Fatal("expected fire training summary")
	}
	if index.FireTraining.FurnaceLevel != 1 {
		t.Fatalf("expected furnace level 1, got %d", index.FireTraining.FurnaceLevel)
	}
	if index.FireTraining.RewardPreview != "焚火晶 x6 / 2小时" {
		t.Fatalf("expected reward preview 焚火晶 x6 / 2小时, got %s", index.FireTraining.RewardPreview)
	}
	if index.War == nil {
		t.Fatal("expected war summary")
	}
	if !index.War.CanRegister {
		t.Fatal("expected leader can register alliance war")
	}
}

func TestService_ApplyAndApprovePromotesApplicantIntoAlliance(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	created, err := svc.CreateAlliance(ctx, 3101, "赤炎盟")
	if err != nil {
		t.Fatalf("expected create alliance success, got %v", err)
	}

	if err := svc.Apply(ctx, 3102, created.Alliance.AllianceID); err != nil {
		t.Fatalf("expected apply success, got %v", err)
	}

	applies, err := svc.ListApplications(ctx, 3101)
	if err != nil {
		t.Fatalf("expected apply list success, got %v", err)
	}
	if len(applies) != 1 {
		t.Fatalf("expected 1 pending application, got %d", len(applies))
	}

	applicantIndex, err := svc.ApproveApplication(ctx, 3101, 3102)
	if err != nil {
		t.Fatalf("expected approve success, got %v", err)
	}

	if !applicantIndex.HasAlliance {
		t.Fatal("expected applicant to join alliance after approval")
	}
	if applicantIndex.CurrentRole != "member" {
		t.Fatalf("expected approved applicant role member, got %s", applicantIndex.CurrentRole)
	}
	if applicantIndex.Alliance == nil || len(applicantIndex.Alliance.Members) != 2 {
		t.Fatalf("expected 2 alliance members after approval, got %#v", applicantIndex.Alliance)
	}
	if applicantIndex.War == nil {
		t.Fatal("expected war summary for member")
	}
	if applicantIndex.War.CanRegister {
		t.Fatal("expected member cannot register alliance war")
	}
}

func TestService_RejectsApplyingToSecondAllianceWhilePending(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	first, err := svc.CreateAlliance(ctx, 3201, "青霄盟")
	if err != nil {
		t.Fatalf("expected first alliance create success, got %v", err)
	}
	second, err := svc.CreateAlliance(ctx, 3202, "玄岳盟")
	if err != nil {
		t.Fatalf("expected second alliance create success, got %v", err)
	}

	if err := svc.Apply(ctx, 3203, first.Alliance.AllianceID); err != nil {
		t.Fatalf("expected first apply success, got %v", err)
	}

	err = svc.Apply(ctx, 3203, second.Alliance.AllianceID)
	if !errors.Is(err, ErrAlreadyApplied) {
		t.Fatalf("expected ErrAlreadyApplied for second pending application, got %v", err)
	}
}
