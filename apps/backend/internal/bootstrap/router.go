package bootstrap

import (
	stdhttp "net/http"

	"github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.TraceID())
	router.Use(middleware.AuthToken())
	router.Use(middleware.IdempotencyKeyMiddleware())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{"status": "ok"}, middleware.GetTraceID(c)))
	})

	return router
}
