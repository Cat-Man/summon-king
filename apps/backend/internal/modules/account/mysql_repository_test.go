package account

import (
	"context"
	"errors"
	"testing"

	"github.com/Cat-Man/summon-king/apps/backend/internal/storage/mysqlstore"
)

func TestMySQLRepository_CreateGuestAndLookup(t *testing.T) {
	store := &fakeGuestAccountStore{}
	repo := NewMySQLRepository(store)

	account, err := repo.CreateGuest(context.Background(), "道友甲", "guest-1001")
	if err != nil {
		t.Fatalf("expected create guest success, got %v", err)
	}
	if account.PlayerID != 1001 || account.Token != "guest-1001" || account.Nickname != "道友甲" {
		t.Fatalf("expected created guest account to round trip, got %+v", account)
	}

	loadedByToken, err := repo.GetByToken(context.Background(), "guest-1001")
	if err != nil {
		t.Fatalf("expected token lookup success, got %v", err)
	}
	if loadedByToken.Nickname != "道友甲" {
		t.Fatalf("expected nickname 道友甲, got %s", loadedByToken.Nickname)
	}

	loadedByPlayer, err := repo.GetByPlayerID(context.Background(), 1001)
	if err != nil {
		t.Fatalf("expected player lookup success, got %v", err)
	}
	if loadedByPlayer.Token != "guest-1001" {
		t.Fatalf("expected token guest-1001, got %s", loadedByPlayer.Token)
	}
}

func TestMySQLRepository_MapsStoreMissToDomainError(t *testing.T) {
	store := &fakeGuestAccountStore{getErr: errors.New("missing")}
	repo := NewMySQLRepository(store)

	if _, err := repo.GetByToken(context.Background(), "guest-missing"); !errors.Is(err, ErrGuestAccountNotFound) {
		t.Fatalf("expected guest account not found, got %v", err)
	}
}

type fakeGuestAccountStore struct {
	record mysqlstore.GuestAccountRecord
	getErr error
}

func (s *fakeGuestAccountStore) CreateGuestAccount(_ context.Context, nickname, token string) (mysqlstore.GuestAccountRecord, error) {
	s.record = mysqlstore.GuestAccountRecord{
		PlayerID: 1001,
		Token:    token,
		Nickname: nickname,
	}
	return s.record, nil
}

func (s *fakeGuestAccountStore) GetGuestAccountByToken(_ context.Context, token string) (mysqlstore.GuestAccountRecord, error) {
	if s.getErr != nil {
		return mysqlstore.GuestAccountRecord{}, s.getErr
	}
	return s.record, nil
}

func (s *fakeGuestAccountStore) GetGuestAccountByPlayerID(_ context.Context, playerID int64) (mysqlstore.GuestAccountRecord, error) {
	if s.getErr != nil {
		return mysqlstore.GuestAccountRecord{}, s.getErr
	}
	if playerID != s.record.PlayerID {
		return mysqlstore.GuestAccountRecord{}, errors.New("unexpected player id")
	}
	return s.record, nil
}
