package bootstrap

import (
	"fmt"

	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/alliance"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/arena"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/commerce"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/home"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/ranking"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/tower"
	"github.com/Cat-Man/summon-king/apps/backend/internal/storage/mysqlstore"
)

type dependencies struct {
	accountService  *account.Service
	homeService     *home.Service
	rankingService  *ranking.Service
	arenaService    *arena.Service
	dungeonService  *dungeon.Service
	allianceService *alliance.Service
	commerceService *commerce.Service
	growthRepo      growth.Repository
	petService      *pet.Service
	towerService    *tower.Service
}

type mysqlStore interface {
	mysqlstore.GuestAccountStore
	mysqlstore.ModuleStateStore
}

var buildMySQLDependenciesFn = buildMySQLDependencies
var openMySQLStore = func(cfg Config) (mysqlStore, error) {
	return mysqlstore.Open(cfg.MySQLDSN)
}

func buildDependencies(cfg Config) (dependencies, error) {
	switch cfg.StorageDriver {
	case "", defaultStorageDriver:
		return buildMemoryDependencies(), nil
	case "mysql":
		return buildMySQLDependenciesFn(cfg)
	default:
		return dependencies{}, fmt.Errorf("unsupported storage driver %q", cfg.StorageDriver)
	}
}

func buildMySQLDependencies(cfg Config) (dependencies, error) {
	store, err := openMySQLStore(cfg)
	if err != nil {
		return dependencies{}, err
	}

	accountRepo := account.NewMySQLRepository(store)
	dungeonRepo := dungeon.NewMySQLRepository(store)
	growthRepo := growth.NewMySQLRepository(store)
	petRepo := pet.NewMySQLRepository(store)
	arenaRepo := arena.NewMySQLRepository(store)
	towerRepo := tower.NewMySQLRepository(store)
	commerceRepo := commerce.NewMySQLRepository(store)

	assetService := asset.NewService(growthRepo)
	accountService := account.NewService(accountRepo)
	allianceService := alliance.NewService(alliance.NewMemoryRepository())
	commerceService := commerce.NewService(commerceRepo, assetService)
	petService := pet.NewService(petRepo, pet.WithGrowthReader(growthRepo))
	dungeonService := dungeon.NewService(
		dungeonRepo,
		assetService,
		dungeon.WithBattleTeamReader(petService),
		dungeon.WithPetProgressor(petService),
	)
	arenaService := arena.NewService(arenaRepo, assetService, arena.WithBattleTeamReader(petService))
	towerService := tower.NewService(towerRepo, assetService, tower.WithBattleTeamReader(petService))
	rankingService := ranking.NewService(accountRepo, growthRepo, dungeonService, arenaService)

	return dependencies{
		accountService:  accountService,
		homeService:     home.NewService(accountRepo, dungeonService, growthRepo, petService, towerService, arenaService, rankingService),
		rankingService:  rankingService,
		arenaService:    arenaService,
		dungeonService:  dungeonService,
		allianceService: allianceService,
		commerceService: commerceService,
		growthRepo:      growthRepo,
		petService:      petService,
		towerService:    towerService,
	}, nil
}

func buildMemoryDependencies() dependencies {
	accountRepo := account.NewMemoryRepository()
	dungeonRepo := dungeon.NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	assetService := asset.NewService(growthRepo)
	accountService := account.NewService(accountRepo)
	allianceService := alliance.NewService(alliance.NewMemoryRepository())
	commerceService := commerce.NewService(commerce.NewMemoryRepository(), assetService)
	petService := pet.NewService(pet.NewMemoryRepository(), pet.WithGrowthReader(growthRepo))
	dungeonService := dungeon.NewService(
		dungeonRepo,
		assetService,
		dungeon.WithBattleTeamReader(petService),
		dungeon.WithPetProgressor(petService),
	)
	arenaService := arena.NewService(arena.NewMemoryRepository(), assetService, arena.WithBattleTeamReader(petService))
	towerService := tower.NewService(tower.NewMemoryRepository(), assetService, tower.WithBattleTeamReader(petService))
	rankingService := ranking.NewService(accountRepo, growthRepo, dungeonService, arenaService)

	return dependencies{
		accountService:  accountService,
		homeService:     home.NewService(accountRepo, dungeonService, growthRepo, petService, towerService, arenaService, rankingService),
		rankingService:  rankingService,
		arenaService:    arenaService,
		dungeonService:  dungeonService,
		allianceService: allianceService,
		commerceService: commerceService,
		growthRepo:      growthRepo,
		petService:      petService,
		towerService:    towerService,
	}
}
