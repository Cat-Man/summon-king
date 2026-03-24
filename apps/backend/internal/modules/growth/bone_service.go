package growth

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) GetWallet(ctx context.Context, playerID int64) Wallet {
	return s.repo.GetWallet(ctx, playerID)
}

func (s *Service) UpgradeBone(ctx context.Context, playerID int64, boneID string) (Bone, error) {
	bone := s.repo.GetBone(ctx, playerID, boneID)
	bone.Level++
	s.repo.SaveBone(ctx, playerID, bone)
	return bone, nil
}
