package growth

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrCropNotReady = errors.New("crop not ready")
)

type Repository interface {
	GetWallet(ctx context.Context, playerID int64) Wallet
	SaveWallet(ctx context.Context, wallet Wallet)
	GetSpirit(ctx context.Context, playerID int64, spiritID int64) Spirit
	SaveSpirit(ctx context.Context, playerID int64, spirit Spirit)
	GetBone(ctx context.Context, playerID int64, boneID string) Bone
	SaveBone(ctx context.Context, playerID int64, bone Bone)
	GetSoul(ctx context.Context, playerID int64, soulID string) Soul
	SaveSoul(ctx context.Context, playerID int64, soul Soul)
	GetCauldron(ctx context.Context, playerID int64) CauldronRecord
	SaveCauldron(ctx context.Context, record CauldronRecord)
	GetAscension(ctx context.Context, playerID int64) AscensionRecord
	SaveAscension(ctx context.Context, record AscensionRecord)
	GetCrop(ctx context.Context, playerID int64) (Crop, bool)
	SaveCrop(ctx context.Context, crop Crop)
}

type MemoryRepository struct {
	mu         sync.Mutex
	wallets    map[int64]Wallet
	spirits    map[int64]map[int64]Spirit
	bones      map[int64]map[string]Bone
	souls      map[int64]map[string]Soul
	cauldrons  map[int64]CauldronRecord
	ascensions map[int64]AscensionRecord
	crops      map[int64]Crop
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		wallets:    make(map[int64]Wallet),
		spirits:    make(map[int64]map[int64]Spirit),
		bones:      make(map[int64]map[string]Bone),
		souls:      make(map[int64]map[string]Soul),
		cauldrons:  make(map[int64]CauldronRecord),
		ascensions: make(map[int64]AscensionRecord),
		crops:      make(map[int64]Crop),
	}
}

func (r *MemoryRepository) GetWallet(_ context.Context, playerID int64) Wallet {
	r.mu.Lock()
	defer r.mu.Unlock()
	wallet, ok := r.wallets[playerID]
	if !ok {
		wallet = Wallet{PlayerID: playerID, SpiritPower: 300, FreeSpiritWashCount: 1, BonePower: 100, SoulPower: 120}
		r.wallets[playerID] = wallet
	}
	return wallet
}

func (r *MemoryRepository) SaveWallet(_ context.Context, wallet Wallet) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.wallets[wallet.PlayerID] = wallet
}

func (r *MemoryRepository) GetSpirit(_ context.Context, playerID int64, spiritID int64) Spirit {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.spirits[playerID]; !ok {
		r.spirits[playerID] = make(map[int64]Spirit)
	}
	spirit, ok := r.spirits[playerID][spiritID]
	if !ok {
		spirit = Spirit{SpiritID: spiritID, AttrLine: "生命+10", SellPrice: 30}
		r.spirits[playerID][spiritID] = spirit
	}
	return spirit
}

func (r *MemoryRepository) SaveSpirit(_ context.Context, playerID int64, spirit Spirit) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.spirits[playerID]; !ok {
		r.spirits[playerID] = make(map[int64]Spirit)
	}
	r.spirits[playerID][spirit.SpiritID] = spirit
}

func (r *MemoryRepository) GetBone(_ context.Context, playerID int64, boneID string) Bone {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.bones[playerID]; !ok {
		r.bones[playerID] = make(map[string]Bone)
	}
	bone, ok := r.bones[playerID][boneID]
	if !ok {
		bone = Bone{BoneID: boneID, Level: 1}
		r.bones[playerID][boneID] = bone
	}
	return bone
}

func (r *MemoryRepository) SaveBone(_ context.Context, playerID int64, bone Bone) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.bones[playerID]; !ok {
		r.bones[playerID] = make(map[string]Bone)
	}
	r.bones[playerID][bone.BoneID] = bone
}

func (r *MemoryRepository) GetSoul(_ context.Context, playerID int64, soulID string) Soul {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.souls[playerID]; !ok {
		r.souls[playerID] = make(map[string]Soul)
	}
	soul, ok := r.souls[playerID][soulID]
	if !ok {
		soul = Soul{SoulID: soulID, Level: 1}
		r.souls[playerID][soulID] = soul
	}
	return soul
}

func (r *MemoryRepository) SaveSoul(_ context.Context, playerID int64, soul Soul) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.souls[playerID]; !ok {
		r.souls[playerID] = make(map[string]Soul)
	}
	r.souls[playerID][soul.SoulID] = soul
}

func (r *MemoryRepository) GetCauldron(_ context.Context, playerID int64) CauldronRecord {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.cauldrons[playerID]
	if !ok {
		rec = CauldronRecord{PlayerID: playerID, CurrentQuality: 1}
		r.cauldrons[playerID] = rec
	}
	return rec
}

func (r *MemoryRepository) SaveCauldron(_ context.Context, record CauldronRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cauldrons[record.PlayerID] = record
}

func (r *MemoryRepository) GetAscension(_ context.Context, playerID int64) AscensionRecord {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.ascensions[playerID]
	if !ok {
		rec = AscensionRecord{PlayerID: playerID, State: "idle"}
		r.ascensions[playerID] = rec
	}
	return rec
}

func (r *MemoryRepository) SaveAscension(_ context.Context, record AscensionRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ascensions[record.PlayerID] = record
}

func (r *MemoryRepository) GetCrop(_ context.Context, playerID int64) (Crop, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	crop, ok := r.crops[playerID]
	return crop, ok
}

func (r *MemoryRepository) SaveCrop(_ context.Context, crop Crop) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.crops[crop.PlayerID] = crop
}

func cropReadyAt(now time.Time) time.Time { return now.Add(8 * time.Hour) }
