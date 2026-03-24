package jobs

import (
	"context"
	"testing"
	"time"
)

func TestDailyReset_ResetsSigninState(t *testing.T) {
	deps := newTestDeps(t)
	deps.DailyStates[1001] = DailyState{PlayerID: 1001, SigninStatus: 1, VIPChestClaimed: true, ArenaRewardClaimed: true}

	job := NewDailyResetJob(deps)
	if err := job.Run(context.Background()); err != nil {
		t.Fatalf("expected daily reset success, got %v", err)
	}

	state := deps.DailyStates[1001]
	if state.SigninStatus != 0 || state.VIPChestClaimed || state.ArenaRewardClaimed {
		t.Fatalf("expected daily state reset, got %+v", state)
	}
}

func TestCultivationTick_MarksReadyWhenDue(t *testing.T) {
	now := time.Date(2026, 3, 22, 10, 0, 0, 0, time.UTC)
	deps := newTestDeps(t)
	deps.Cultivations[1001] = CultivationTask{PlayerID: 1001, ReadyAt: now.Add(-time.Minute)}

	job := NewCultivationTickJob(deps, func() time.Time { return now })
	if err := job.Run(context.Background()); err != nil {
		t.Fatalf("expected cultivation tick success, got %v", err)
	}

	if !deps.Cultivations[1001].Claimable {
		t.Fatal("expected cultivation task claimable after tick")
	}
}

func TestManorTick_MarksPlotReady(t *testing.T) {
	now := time.Date(2026, 3, 22, 10, 0, 0, 0, time.UTC)
	deps := newTestDeps(t)
	deps.ManorPlots["plot-1"] = ManorPlot{PlotID: "plot-1", MatureAt: now.Add(-time.Minute), Status: "growing"}

	job := NewManorTickJob(deps, func() time.Time { return now })
	if err := job.Run(context.Background()); err != nil {
		t.Fatalf("expected manor tick success, got %v", err)
	}

	if deps.ManorPlots["plot-1"].Status != "ready" {
		t.Fatalf("expected plot ready, got %s", deps.ManorPlots["plot-1"].Status)
	}
}

func TestAllianceWarPhase_AdvancesBySchedule(t *testing.T) {
	now := time.Date(2026, 3, 22, 20, 0, 0, 0, time.UTC)
	deps := newTestDeps(t)
	deps.AllianceWars["round-1"] = AllianceWarRound{
		RoundID:  "round-1",
		Phase:    "signup",
		LockAt:   now.Add(-time.Minute),
		BattleAt: now.Add(time.Hour),
		ResultAt: now.Add(2 * time.Hour),
	}

	job := NewAllianceWarPhaseJob(deps, func() time.Time { return now })
	if err := job.Run(context.Background()); err != nil {
		t.Fatalf("expected alliance war phase tick success, got %v", err)
	}

	if deps.AllianceWars["round-1"].Phase != "locked" {
		t.Fatalf("expected war phase locked, got %s", deps.AllianceWars["round-1"].Phase)
	}
}
