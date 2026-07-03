package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *Handler, mode string) *gin.Engine {
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	r.GET("/health", handler.HealthCheck)

	r.Static("/web", "./web")
	r.StaticFile("/", "./web/public/index.html")

	api := r.Group("/api/v1")
	{
		products := api.Group("/products")
		{
			products.GET("", handler.GetProducts)
			products.GET("/hot", handler.GetHotProducts)
			products.GET("/search", handler.SearchProducts)
			products.GET("/:id", handler.GetProduct)
			products.GET("/:id/trends", handler.GetPriceTrends)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", handler.GetCategories)
			categories.GET("/stats", handler.GetCategoryStats)
		}

		platforms := api.Group("/platforms")
		{
			platforms.GET("", handler.GetPlatforms)
			platforms.GET("/stats", handler.GetPlatformStats)
		}

		insights := api.Group("/insights")
		{
			insights.GET("/market", handler.GetMarketInsights)
			insights.GET("/keywords", handler.GetTrendingKeywords)
		}
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
