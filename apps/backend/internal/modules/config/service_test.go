package config

import (
	"context"
	"testing"
)

func TestPublishAndRollbackDraft(t *testing.T) {
	svc := NewService()
	ctx := context.Background()

	if _, err := svc.SaveDraft(ctx, "vip_shop", `{"price":100}`); err != nil {
		t.Fatalf("expected save draft success, got %v", err)
	}
	first, err := svc.Publish(ctx, "vip_shop", 9001, "init publish")
	if err != nil {
		t.Fatalf("expected first publish success, got %v", err)
	}
	if first.Version != 1 {
		t.Fatalf("expected first version 1, got %d", first.Version)
	}

	if _, err := svc.SaveDraft(ctx, "vip_shop", `{"price":120}`); err != nil {
		t.Fatalf("expected save second draft success, got %v", err)
	}
	second, err := svc.Publish(ctx, "vip_shop", 9001, "raise price")
	if err != nil {
		t.Fatalf("expected second publish success, got %v", err)
	}
	if second.Version != 2 {
		t.Fatalf("expected second version 2, got %d", second.Version)
	}

	if _, err := svc.Rollback(ctx, "vip_shop", 1, 9001); err != nil {
		t.Fatalf("expected rollback success, got %v", err)
	}
	detail, err := svc.GetDetail(ctx, "vip_shop")
	if err != nil {
		t.Fatalf("expected detail success, got %v", err)
	}
	if detail.PublishedContent != `{"price":100}` {
		t.Fatalf("expected rollback to first content, got %s", detail.PublishedContent)
	}
}

func TestValidateRejectsEmptyDraft(t *testing.T) {
	svc := NewService()
	if _, err := svc.Validate(context.Background(), "vip_shop"); err == nil {
		t.Fatal("expected validate without draft to fail")
	}
}
