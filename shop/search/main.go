package main

import (
	"elastic-talk-search/config"
	"elastic-talk-search/service"
	"elastic-talk-search/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func searchProducts(c *gin.Context, search *service.SearchService) {
	query := c.Query("q")
	pageStr := c.Query("page")
	page := 0
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}
	total, hits := search.Search(query, page)
	c.JSON(200, gin.H{
		"total": total,
		"hits":  hits,
	})
}

func main() {
	utils.InitializeLogger()
	config := config.Get()
	rootLogger := utils.NewLogger()
	logger := utils.NewApplicationLogger(rootLogger)
	logger.Info().Msg("Starting")
	search := service.SearchService{
		Logger: logger,
	}
	search.LoadProducts()
	logger.Info().Msg("All Products were loaded")
	r := gin.New()
	r.Use(utils.NewAccessLogger(rootLogger))
	r.Use(utils.NewRecoveryHandler(rootLogger))
	r.GET("/api/products/search", func(c *gin.Context) {
		searchProducts(c, &search)
	})
	logger.Info().Msg("Started")
	r.Run(":" + strconv.Itoa(config.Server.Port)) // listen and serve on 0.0.0.0:8080
}
