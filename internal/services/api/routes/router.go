package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	httplib "github.com/unanoc/blockchain-indexer/pkg/http"
)

// API contains generic methods that should be used by other API interfaces.
type API interface {
	Setup(router *gin.RouterGroup)
}

func NewRouter(db repository.Storage, p *prometheus.Prometheus) http.Handler {
	var router *gin.Engine

	if config.Default.Gin.Mode == gin.DebugMode {
		router = gin.Default()
	} else {
		router = gin.New()
	}

	// Setup service routers
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup middlewares
	router.Use(httplib.CORSMiddleware())
	router.Use(prometheus.GinMetricsMiddleware(p))

	apiRouteGroup := router.Group("/api/v1")
	// Setup API routes
	NewTransactionsAPI(db).Setup(apiRouteGroup.Group("/transactions"))
	NewCollectionsAPI(db).Setup(apiRouteGroup.Group("/collections"))

	return router
}
