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

func (s *BoneService) Upgrade(ctx context.Context, playerID int64) (Bone, error) {
	wallet, err := s.repo.UpgradeBoneLevel(ctx, playerID, 1)
	if err != nil {
		return Bone{}, err
	}
	return Bone{
		Name:  "战骨",
		Level: wallet.BoneLevel,
	}, nil
}
