package growth

import "context"

type BoneService struct {
	repo Repository
}

func NewBoneService(repo Repository) *BoneService {
	return &BoneService{repo: repo}
}

func (s *BoneService) GetBoneState(ctx context.Context, playerID int64) Bone {
	wallet, _ := s.repo.GetWallet(ctx, playerID)
	return Bone{
		Name:  "战骨",
		Level: wallet.BoneLevel,
	}
}
