package bootstrap

import (
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	accountmod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	playermod "github.com/Cat-Man/summon-king/apps/backend/internal/modules/player"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	accountRepo := accountmod.NewMemoryRepository()
	accountService := accountmod.NewService(accountRepo)
	playerService := playermod.NewService(playermod.NewRepository(accountRepo))

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

	return router
}
