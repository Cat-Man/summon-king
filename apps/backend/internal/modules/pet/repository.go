package pet

import "context"

type Repository interface {
	GetBattleTeam(ctx context.Context, playerID int64) (TeamSnapshot, error)
}

type MemoryRepository struct{}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) GetBattleTeam(_ context.Context, playerID int64) (TeamSnapshot, error) {
	starter := BattlePet{
		PetID:    playerID*10 + 1,
		Slot:     1,
		Name:     "初始灵狐",
		Level:    1,
		Power:    120,
		IsActive: true,
	}

	return TeamSnapshot{
		PlayerID: playerID,
		Pets:     []BattlePet{starter},
	}, nil
}
