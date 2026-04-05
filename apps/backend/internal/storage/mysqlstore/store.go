package mysqlstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type GuestAccountRecord struct {
	PlayerID int64  `json:"player_id"`
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
}

type GuestAccountStore interface {
	CreateGuestAccount(ctx context.Context, nickname, token string) (GuestAccountRecord, error)
	GetGuestAccountByToken(ctx context.Context, token string) (GuestAccountRecord, error)
	GetGuestAccountByPlayerID(ctx context.Context, playerID int64) (GuestAccountRecord, error)
}

type ModuleStateStore interface {
	LoadModuleState(ctx context.Context, playerID int64, module string, target any) (bool, error)
	SaveModuleState(ctx context.Context, playerID int64, module string, value any) error
}

type Store struct {
	db *sql.DB
}

func Open(dsn string) (*Store, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	store := &Store{db: db}
	if err := store.ensureSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) CreateGuestAccount(ctx context.Context, nickname, token string) (GuestAccountRecord, error) {
	result, err := s.db.ExecContext(
		ctx,
		`INSERT INTO guest_accounts (token, nickname) VALUES (?, ?)`,
		token,
		nickname,
	)
	if err != nil {
		return GuestAccountRecord{}, fmt.Errorf("insert guest account: %w", err)
	}

	playerID, err := result.LastInsertId()
	if err != nil {
		return GuestAccountRecord{}, fmt.Errorf("read guest account id: %w", err)
	}

	return GuestAccountRecord{
		PlayerID: playerID,
		Token:    token,
		Nickname: nickname,
	}, nil
}

func (s *Store) GetGuestAccountByToken(ctx context.Context, token string) (GuestAccountRecord, error) {
	return s.getGuestAccount(ctx, `SELECT player_id, token, nickname FROM guest_accounts WHERE token = ?`, token)
}

func (s *Store) GetGuestAccountByPlayerID(ctx context.Context, playerID int64) (GuestAccountRecord, error) {
	return s.getGuestAccount(ctx, `SELECT player_id, token, nickname FROM guest_accounts WHERE player_id = ?`, playerID)
}

func (s *Store) LoadModuleState(ctx context.Context, playerID int64, module string, target any) (bool, error) {
	var payload []byte
	err := s.db.QueryRowContext(
		ctx,
		`SELECT state_json FROM module_states WHERE player_id = ? AND module = ?`,
		playerID,
		module,
	).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("load module state: %w", err)
	}
	if err := json.Unmarshal(payload, target); err != nil {
		return false, fmt.Errorf("decode module state: %w", err)
	}
	return true, nil
}

func (s *Store) SaveModuleState(ctx context.Context, playerID int64, module string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode module state: %w", err)
	}

	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO module_states (player_id, module, state_json)
		 VALUES (?, ?, ?)
		 ON DUPLICATE KEY UPDATE state_json = VALUES(state_json), updated_at = CURRENT_TIMESTAMP`,
		playerID,
		module,
		payload,
	)
	if err != nil {
		return fmt.Errorf("save module state: %w", err)
	}
	return nil
}

func (s *Store) getGuestAccount(ctx context.Context, query string, arg any) (GuestAccountRecord, error) {
	var record GuestAccountRecord
	err := s.db.QueryRowContext(ctx, query, arg).Scan(&record.PlayerID, &record.Token, &record.Nickname)
	if err != nil {
		return GuestAccountRecord{}, err
	}
	return record, nil
}

func (s *Store) ensureSchema(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS guest_accounts (
			player_id BIGINT NOT NULL AUTO_INCREMENT,
			token VARCHAR(128) NOT NULL,
			nickname VARCHAR(128) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (player_id),
			UNIQUE KEY uniq_guest_accounts_token (token)
		) AUTO_INCREMENT = 1001`,
		`CREATE TABLE IF NOT EXISTS module_states (
			player_id BIGINT NOT NULL,
			module VARCHAR(64) NOT NULL,
			state_json LONGTEXT NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (player_id, module)
		)`,
	}

	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("ensure mysql schema: %w", err)
		}
	}
	return nil
}
