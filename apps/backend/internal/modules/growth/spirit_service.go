package growth

import "context"

func (s *Service) WashSpirit(ctx context.Context, playerID int64, spiritID int64, _ WashOption) error {
	wallet := s.repo.GetWallet(ctx, playerID)
	if wallet.FreeSpiritWashCount > 0 {
		wallet.FreeSpiritWashCount--
	} else {
		wallet.SpiritPower -= 10
	}
	s.repo.SaveWallet(ctx, wallet)

	spirit := s.repo.GetSpirit(ctx, playerID, spiritID)
	spirit.WashCount++
	spirit.AttrLine = "攻击+12"
	s.repo.SaveSpirit(ctx, playerID, spirit)
	return nil
}
