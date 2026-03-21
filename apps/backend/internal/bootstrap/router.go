package bootstrap

import (
	"context"
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	accountmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	assetmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/asset"
	battlemod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/battle"
	petmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
	playermod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/player"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	accountRepo := accountmod.NewMemoryRepository()
	assetRepo := assetmod.NewMemoryRepository()
	accountService := accountmod.NewService(accountRepo)
	playerService := playermod.NewService(playermod.NewRepository(accountRepo, assetRepo))
	assetService := assetmod.NewService(assetRepo)
	petService := petmod.NewService(petmod.NewMemoryRepository())
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

	return router
}
