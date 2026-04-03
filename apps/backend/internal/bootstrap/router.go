package bootstrap

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/arena"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/home"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/tower"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.InjectTraceID())
	router.Use(middleware.InjectAuthToken())
	router.Use(middleware.InjectIdempotencyKey())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{
			"status": "ok",
			"app":    defaultAppName,
		}, middleware.GetTraceID(c)))
	})
	api := router.Group("/api/v1")

	accountRepo := account.NewMemoryRepository()
	dungeonRepo := dungeon.NewMemoryRepository()
	growthRepo := growth.NewMemoryRepository()
	accountService := account.NewService(accountRepo)
	dungeonService := dungeon.NewService(dungeonRepo, growthRepo)

	authGroup := api.Group("/auth")
	account.NewHandler(accountService).RegisterRoutes(authGroup)

	homeGroup := api.Group("/home")
	home.NewHandler(home.NewService(accountRepo, dungeonService, growthRepo)).RegisterRoutes(homeGroup)

	arenaGroup := api.Group("/arena")
	registerModuleRoot(arenaGroup, "arena")
	arena.NewHandler(arena.NewService(arena.NewMemoryRepository())).RegisterRoutes(arenaGroup)

	dungeonGroup := api.Group("/dungeon")
	registerModuleRoot(dungeonGroup, "dungeon")
	dungeon.NewHandler(dungeonService).RegisterRoutes(dungeonGroup)

	growthGroup := api.Group("/growth")
	registerModuleRoot(growthGroup, "growth")
	growth.NewHandler(growthRepo).RegisterRoutes(growthGroup)

	towerGroup := api.Group("/tower")
	registerModuleRoot(towerGroup, "tower")
	tower.NewHandler(tower.NewService(tower.NewMemoryRepository())).RegisterRoutes(towerGroup)

	return router
}

func registerModuleRoot(group *gin.RouterGroup, module string) {
	group.GET("", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{
			"module": module,
			"status": "ok",
		}, middleware.GetTraceID(c)))
	})
}
