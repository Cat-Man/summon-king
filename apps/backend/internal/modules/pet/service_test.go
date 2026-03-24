package pet_test

import (
	"context"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/player"
	"github.com/gin-gonic/gin"
)

func TestSaveTeam_RejectsDuplicatePetInMultipleSlots(t *testing.T) {
	svc := pet.NewService(pet.NewMemoryRepository())

	err := svc.SaveTeam(context.Background(), 1001, []int64{1, 1, 2})
	if err == nil {
		t.Fatal("expected duplicate pet assignment error, got nil")
	}
	if err != pet.ErrDuplicatePetInTeam {
		t.Fatalf("expected pet.ErrDuplicatePetInTeam, got %v", err)
	}
}

func TestCatalogAndDetail_AreAvailable(t *testing.T) {
	svc := pet.NewService(pet.NewMemoryRepository())

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

func TestGetTeam_ReturnsDefaultTeamOrderedByPower(t *testing.T) {
	svc := pet.NewService(pet.NewMemoryRepository())

	team, err := svc.GetTeam(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int64{3, 2, 1}
	if !reflect.DeepEqual(team.PetIDs, expected) {
		t.Fatalf("expected default team %v, got %v", expected, team.PetIDs)
	}
}

func TestGetTeam_ReturnsSavedTeam(t *testing.T) {
	svc := pet.NewService(pet.NewMemoryRepository())
	saved := []int64{2, 1}

	if err := svc.SaveTeam(context.Background(), 1001, saved); err != nil {
		t.Fatalf("expected save team success, got %v", err)
	}

	team, err := svc.GetTeam(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(team.PetIDs, saved) {
		t.Fatalf("expected saved team %v, got %v", saved, team.PetIDs)
	}
}

func TestHandler_TeamReturnsCurrentTeam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := pet.NewService(pet.NewMemoryRepository())
	router := gin.New()
	router.Use(middleware.TraceID())
	pet.NewHandler(svc).RegisterRoutes(router.Group("/api/v1/player/pets"))

	req := httptest.NewRequest(stdhttp.MethodGet, "/api/v1/player/pets/team?player_id=1001", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != stdhttp.StatusOK {
		t.Fatalf("expected http 200, got %d", rec.Code)
	}

	var envelope struct {
		Code int         `json:"code"`
		Data pet.PetTeam `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if envelope.Code != 0 {
		t.Fatalf("expected business code 0, got %d", envelope.Code)
	}
	expected := []int64{3, 2, 1}
	if !reflect.DeepEqual(envelope.Data.PetIDs, expected) {
		t.Fatalf("expected team %v, got %v", expected, envelope.Data.PetIDs)
	}
}

func TestHomeIndex_PowerUsesSavedTeamFirst(t *testing.T) {
	ctx := context.Background()
	accountRepo := account.NewMemoryRepository()
	assetRepo := asset.NewMemoryRepository()
	accountSvc := account.NewService(accountRepo)

	playerEntity, err := accountSvc.CreateGuestPlayer(ctx, "web")
	if err != nil {
		t.Fatalf("expected create guest success, got %v", err)
	}

	petSvc := pet.NewService(pet.NewMemoryRepository())
	if err := petSvc.SaveTeam(ctx, playerEntity.PlayerID, []int64{1, 2}); err != nil {
		t.Fatalf("expected save team success, got %v", err)
	}

	playerSvc := player.NewService(
		player.NewRepository(accountRepo, assetRepo),
		player.WithPetReader(petSvc),
	)

	home, err := playerSvc.GetHomeIndex(ctx, playerEntity.PlayerID)
	if err != nil {
		t.Fatalf("expected home index success, got %v", err)
	}

	if got := findResourceValue(home.Resources, "战力"); got != "270" {
		t.Fatalf("expected team power 270, got %s", got)
	}
}

func findResourceValue(resources []player.HomeResource, label string) string {
	for _, item := range resources {
		if item.Label == label {
			return item.Value
		}
	}
	return ""
}
