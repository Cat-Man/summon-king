package pet

import (
	"context"
	"testing"
)

func TestSaveTeam_RejectsDuplicatePetInMultipleSlots(t *testing.T) {
	svc := NewService(NewMemoryRepository())

	err := svc.SaveTeam(context.Background(), 1001, []int64{1, 1, 2})
	if err == nil {
		t.Fatal("expected duplicate pet assignment error, got nil")
	}
	if err != ErrDuplicatePetInTeam {
		t.Fatalf("expected ErrDuplicatePetInTeam, got %v", err)
	}
}

func TestCatalogAndDetail_AreAvailable(t *testing.T) {
	svc := NewService(NewMemoryRepository())

	catalog, err := svc.GetCatalog(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(catalog) == 0 {
		t.Fatal("expected non-empty catalog")
	}

	detail, err := svc.GetDetail(context.Background(), 1001, catalog[0].PetID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if detail.PetID == 0 {
		t.Fatal("expected valid pet detail")
	}
}
