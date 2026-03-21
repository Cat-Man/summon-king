package bootstrap

import (
	"context"
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	accountmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	arenamod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/arena"
	assetmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	battlemod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/battle"
	dungeonmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
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
	playerService := playermod.NewService(playermod.NewRepository(accountRepo, assetRepo))
	assetService := assetmod.NewService(assetRepo)
	petService := petmod.NewService(petmod.NewMemoryRepository())
	dungeonService := dungeonmod.NewService(dungeonmod.NewMemoryRepository())
	growthService := growthmod.NewService(growthmod.NewMemoryRepository())
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

	return router
}
