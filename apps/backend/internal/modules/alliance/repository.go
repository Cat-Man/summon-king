package alliance

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrLevelTooLow          = errors.New("player level too low")
	ErrAllianceNotFound     = errors.New("alliance not found")
	ErrFireTrainingNotReady = errors.New("fire training not ready")
	ErrFireTrainingClaimed  = errors.New("fire training already claimed")
)

type Repository interface {
	GetPlayerLevel(ctx context.Context, playerID int64) int
	CreateAlliance(ctx context.Context, ownerID int64, name string, now time.Time) (Alliance, error)
	GetAlliance(ctx context.Context, allianceID string) (Alliance, error)
	SaveAlliance(ctx context.Context, alliance Alliance) error
	StartFireTraining(ctx context.Context, allianceID string, playerID int64, roomID string, now time.Time) (FireTrainingRecord, error)
	GetFireTraining(ctx context.Context, allianceID, recordID string) (FireTrainingRecord, error)
	SaveFireTraining(ctx context.Context, record FireTrainingRecord) error
	SaveWarSignUp(ctx context.Context, record AllianceWarRecord) error
	GetWarRecord(ctx context.Context, allianceID string) (AllianceWarRecord, bool)
}

type MemoryRepository struct {
	mu               sync.Mutex
	nextAllianceID   int64
	nextTrainingID   int64
	playerLevels     map[int64]int
	alliances        map[string]Alliance
	fireTraining     map[string]FireTrainingRecord
	allianceWarByAid map[string]AllianceWarRecord
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextAllianceID:   1,
		nextTrainingID:   1,
		playerLevels:     map[int64]int{1001: 10, 3001: 35, 3002: 40},
		alliances:        make(map[string]Alliance),
		fireTraining:     make(map[string]FireTrainingRecord),
		allianceWarByAid: make(map[string]AllianceWarRecord),
	}
}

func (r *MemoryRepository) GetPlayerLevel(_ context.Context, playerID int64) int {
	return r.playerLevels[playerID]
}

func (r *MemoryRepository) CreateAlliance(_ context.Context, ownerID int64, name string, now time.Time) (Alliance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := fmt.Sprintf("alliance-%d", r.nextAllianceID)
	r.nextAllianceID++
	alliance := Alliance{
		AllianceID: id,
		Name:       name,
		OwnerID:    ownerID,
		Notice:     "欢迎加入联盟",
		Funds:      0,
		Buildings:  map[string]int{"hall": 1, "warehouse": 1},
		Members:    []int64{ownerID},
		CreatedAt:  now,
	}
	r.alliances[id] = alliance
	return alliance, nil
}

func (r *MemoryRepository) GetAlliance(_ context.Context, allianceID string) (Alliance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	alliance, ok := r.alliances[allianceID]
	if !ok {
		return Alliance{}, ErrAllianceNotFound
	}
	return alliance, nil
}

func (r *MemoryRepository) SaveAlliance(_ context.Context, alliance Alliance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.alliances[alliance.AllianceID] = alliance
	return nil
}

func (r *MemoryRepository) StartFireTraining(_ context.Context, allianceID string, playerID int64, roomID string, now time.Time) (FireTrainingRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.alliances[allianceID]; !ok {
		return FireTrainingRecord{}, ErrAllianceNotFound
	}
	record := FireTrainingRecord{
		RecordID:      fmt.Sprintf("fire-%d", r.nextTrainingID),
		AllianceID:    allianceID,
		PlayerID:      playerID,
		RoomID:        roomID,
		Status:        "running",
		RewardCrystal: 120,
		StartedAt:     now,
		FinishedAt:    now.Add(2 * time.Hour),
	}
	r.nextTrainingID++
	r.fireTraining[record.RecordID] = record
	return record, nil
}

func (r *MemoryRepository) GetFireTraining(_ context.Context, allianceID, recordID string) (FireTrainingRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.fireTraining[recordID]
	if !ok || record.AllianceID != allianceID {
		return FireTrainingRecord{}, ErrAllianceNotFound
	}
	return record, nil
}

func (r *MemoryRepository) SaveFireTraining(_ context.Context, record FireTrainingRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fireTraining[record.RecordID] = record
	return nil
}

func (r *MemoryRepository) SaveWarSignUp(_ context.Context, record AllianceWarRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.allianceWarByAid[record.AllianceID] = record
	return nil
}

func (r *MemoryRepository) GetWarRecord(_ context.Context, allianceID string) (AllianceWarRecord, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.allianceWarByAid[allianceID]
	return record, ok
}
