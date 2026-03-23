package bootstrap

import (
	"context"
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	accountmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	alliancemod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/alliance"
	arenamod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/arena"
	assetmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	battlemod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/battle"
	wxminimod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/bridge/wxmini"
	commercemod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/commerce"
	configmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/config"
	dungeonmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	gmmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/gm"
	growthmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	petmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
	playermod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/player"
	rankingmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/ranking"
	towermod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/tower"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	accountRepo := accountmod.NewMemoryRepository()
	assetRepo := assetmod.NewMemoryRepository()
	accountService := accountmod.NewService(accountRepo)
	assetService := assetmod.NewService(assetRepo)
	petService := petmod.NewService(petmod.NewMemoryRepository())
	dungeonService := dungeonmod.NewService(dungeonmod.NewMemoryRepository())
	growthService := growthmod.NewService(growthmod.NewMemoryRepository())
	allianceService := alliancemod.NewService(alliancemod.NewMemoryRepository())
	commerceService := commercemod.NewService(commercemod.NewMemoryRepository())
	playerService := playermod.NewService(
		playermod.NewRepository(accountRepo, assetRepo),
		playermod.WithPetReader(petService),
		playermod.WithDungeonReader(dungeonService),
		playermod.WithCommerceReader(commerceService),
	)
	configService := configmod.NewService()
	gmService := gmmod.NewService()
	wxminiService := wxminimod.NewService()
	towerService := towermod.NewService()
	arenaService := arenamod.NewService()
	rankingService := rankingmod.NewService()
	battleStore := battlemod.NewMemoryStore()
	_ = battleStore.Save(context.Background(), battlemod.RunBattle(battlemod.BattleInput{
		BattleNo: "battle-demo-1",
		Left: []battlemod.Unit{
			{Name: "烈焰狼", HP: 160, Attack: 72, Defense: 18, Speed: 22},
		},
		Right: []battlemod.Unit{
			{Name: "林地史莱姆", HP: 130, Attack: 50, Defense: 12, Speed: 14},
		},
	}))
	rankingService.RecordScore("power", 1001, 1280)
	rankingService.RecordScore("power", 1002, 1170)
	rankingService.RecordScore("arena-daily", 1001, 9)
	rankingService.RecordScore("arena-daily", 1002, 6)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.TraceID())
	router.Use(middleware.AuthToken())
	router.Use(middleware.IdempotencyKeyMiddleware())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{"status": "ok"}, middleware.GetTraceID(c)))
	})
	router.GET("/metrics", func(c *gin.Context) {
		c.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		c.String(stdhttp.StatusOK, "# HELP summon_king_build_info Build information.\n# TYPE summon_king_build_info gauge\nsummon_king_build_info{service=\"api\"} 1\n# HELP summon_king_health_status Service health status.\n# TYPE summon_king_health_status gauge\nsummon_king_health_status 1\n")
	})

	playerGroup := router.Group("/api/v1/player")
	accountmod.NewHandler(accountService).RegisterRoutes(playerGroup.Group("/auth"))
	playermod.NewHandler(playerService).RegisterRoutes(playerGroup.Group("/home"))
	assetmod.NewHandler(assetService).RegisterRoutes(playerGroup.Group("/assets"))
	petmod.NewHandler(petService).RegisterRoutes(playerGroup.Group("/pets"))
	battlemod.NewHandler(battleStore).RegisterRoutes(playerGroup.Group("/battles"))
	dungeonmod.NewHandler(dungeonService).RegisterRoutes(playerGroup)
	growthmod.NewHandler(growthService).RegisterRoutes(playerGroup.Group("/growth"))
	towermod.NewHandler(towerService).RegisterRoutes(playerGroup.Group("/tower"))
	arenamod.NewHandler(arenaService).RegisterRoutes(playerGroup.Group("/arena"))
	rankingmod.NewHandler(rankingService).RegisterRoutes(playerGroup.Group("/rankings"))
	alliancemod.NewHandler(allianceService).RegisterRoutes(playerGroup.Group("/alliance"))
	commercemod.NewHandler(commerceService).RegisterRoutes(playerGroup.Group("/commerce"))

	bridgeGroup := router.Group("/api/v1/bridge")
	wxminimod.NewHandler(wxminiService).RegisterRoutes(bridgeGroup.Group("/wxmini"))

	adminGroup := router.Group("/api/v1/admin")
	configmod.NewHandler(configService).RegisterRoutes(adminGroup.Group("/config"))
	gmmod.NewHandler(gmService).RegisterRoutes(adminGroup.Group("/gm"))

	return router
}
