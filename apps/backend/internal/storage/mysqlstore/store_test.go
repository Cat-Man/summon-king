package mysqlstore

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStore_EnsureSchemaAndModuleStateRoundTrip(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("expected sqlmock db, got %v", err)
	}
	defer db.Close()

	store := &Store{db: db}

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS guest_accounts").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS module_states").WillReturnResult(sqlmock.NewResult(0, 0))

	if err := store.ensureSchema(context.Background()); err != nil {
		t.Fatalf("expected schema creation success, got %v", err)
	}

	payload := map[string]any{
		"current_streak": 3,
		"last_win":       true,
	}
	mock.ExpectExec("INSERT INTO module_states").
		WithArgs(int64(1001), "arena", jsonArg(t, payload)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	if err := store.SaveModuleState(context.Background(), 1001, "arena", payload); err != nil {
		t.Fatalf("expected save module state success, got %v", err)
	}

	rows := sqlmock.NewRows([]string{"state_json"}).AddRow(`{"current_streak":3,"last_win":true}`)
	mock.ExpectQuery("SELECT state_json FROM module_states").
		WithArgs(int64(1001), "arena").
		WillReturnRows(rows)

	var loaded map[string]any
	found, err := store.LoadModuleState(context.Background(), 1001, "arena", &loaded)
	if err != nil {
		t.Fatalf("expected load module state success, got %v", err)
	}
	if !found {
		t.Fatal("expected saved module state to be found")
	}
	if loaded["current_streak"] != float64(3) {
		t.Fatalf("expected current_streak 3, got %#v", loaded["current_streak"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestStore_LoadModuleStateReturnsFalseWhenMissing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("expected sqlmock db, got %v", err)
	}
	defer db.Close()

	store := &Store{db: db}

	mock.ExpectQuery("SELECT state_json FROM module_states").
		WithArgs(int64(1002), "growth").
		WillReturnError(sql.ErrNoRows)

	var loaded map[string]any
	found, err := store.LoadModuleState(context.Background(), 1002, "growth", &loaded)
	if err != nil {
		t.Fatalf("expected missing row to return nil error, got %v", err)
	}
	if found {
		t.Fatal("expected found=false when row is missing")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func jsonArg(t *testing.T, value any) driver.Value {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("expected JSON marshal success, got %v", err)
	}
	return raw
}
