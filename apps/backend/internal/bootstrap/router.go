package bootstrap

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/arena"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
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
			"app":    LoadConfig().AppName,
		}, middleware.GetTraceID(c)))
	})

	api := router.Group("/api/v1")

	arena.NewHandler(arena.NewService(arena.NewMemoryRepository())).RegisterRoutes(api.Group("/arena"))
	dungeon.NewHandler(dungeon.NewService(dungeon.NewMemoryRepository())).RegisterRoutes(api.Group("/dungeon"))
	growth.NewHandler(growth.NewMemoryRepository()).RegisterRoutes(api.Group("/growth"))
	api.Group("/tower")

	return router
}
