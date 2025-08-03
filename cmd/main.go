package main

import (
	"log"

	"github.com/Abdullah-Shaikh01/monk-coupons/config"
	"github.com/Abdullah-Shaikh01/monk-coupons/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()
	defer db.Close()

	router := gin.Default()

	// Pass DB and register all routes
	routes.RegisterRoutes(router, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
