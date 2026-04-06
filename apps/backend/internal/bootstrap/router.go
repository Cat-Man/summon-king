package bootstrap

import (
	"fmt"
	stdhttp "net/http"

	httpx "github.com/Cat-Man/summon-king/apps/backend/internal/infra/http"
	"github.com/Cat-Man/summon-king/apps/backend/internal/middleware"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/account"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/alliance"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/arena"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/bridge/wxmini"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/commerce"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/dungeon"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/growth"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/home"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/pet"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/ranking"
	"github.com/Cat-Man/summon-king/apps/backend/internal/modules/tower"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	cfg, err := LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	router, err := NewRouterWithConfig(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to build router: %v", err))
	}
	return router
}

func NewRouterWithConfig(cfg Config) (*gin.Engine, error) {
	deps, err := buildDependencies(cfg)
	if err != nil {
		return nil, err
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.InjectTraceID())
	router.Use(middleware.InjectAuthToken())
	router.Use(middleware.InjectIdempotencyKey())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{
			"status": "ok",
			"app":    cfg.AppName,
		}, middleware.GetTraceID(c)))
	})
	api := router.Group("/api/v1")

	authGroup := api.Group("/auth")
	account.NewHandler(deps.accountService).RegisterRoutes(authGroup)

	wxminiBridgeGroup := api.Group("/bridge/wxmini")
	registerModuleRoot(wxminiBridgeGroup, "bridge/wxmini")
	wxmini.NewHandler(deps.wxminiService).RegisterRoutes(wxminiBridgeGroup)

	homeGroup := api.Group("/home")
	home.NewHandler(deps.homeService).RegisterRoutes(homeGroup)

	allianceGroup := api.Group("/alliance")
	registerModuleRoot(allianceGroup, "alliance")
	alliance.NewHandler(deps.allianceService).RegisterRoutes(allianceGroup)

	allianceWarGroup := api.Group("/alliance-war")
	registerModuleRoot(allianceWarGroup, "alliance-war")
	alliance.NewWarHandler(deps.allianceWar).RegisterRoutes(allianceWarGroup)

	commerceHandler := commerce.NewHandler(deps.commerceService)
	signinGroup := api.Group("/signin")
	registerModuleRoot(signinGroup, "signin")
	commerceHandler.RegisterSigninRoutes(signinGroup)

	vipGroup := api.Group("/vip")
	registerModuleRoot(vipGroup, "vip")
	commerceHandler.RegisterVIPRoutes(vipGroup)

	rankingGroup := api.Group("/ranking")
	ranking.NewHandler(deps.rankingService).RegisterRoutes(rankingGroup)

	arenaGroup := api.Group("/arena")
	registerModuleRoot(arenaGroup, "arena")
	arena.NewHandler(deps.arenaService).RegisterRoutes(arenaGroup)

	dungeonGroup := api.Group("/dungeon")
	registerModuleRoot(dungeonGroup, "dungeon")
	dungeon.NewHandler(deps.dungeonService).RegisterRoutes(dungeonGroup)

	growthGroup := api.Group("/growth")
	registerModuleRoot(growthGroup, "growth")
	growth.NewHandler(deps.growthRepo).RegisterRoutes(growthGroup)

	petGroup := api.Group("/pet")
	registerModuleRoot(petGroup, "pet")
	pet.NewHandler(deps.petService).RegisterRoutes(petGroup)

	towerGroup := api.Group("/tower")
	registerModuleRoot(towerGroup, "tower")
	tower.NewHandler(deps.towerService).RegisterRoutes(towerGroup)

	return router, nil
}

func registerModuleRoot(group *gin.RouterGroup, module string) {
	group.GET("", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, httpx.Success(gin.H{
			"module": module,
			"status": "ok",
		}, middleware.GetTraceID(c)))
	})
}
