package main

import (
	"fmt"
	"log"

	"companyInternalManagement/config"
	"companyInternalManagement/models"
	"companyInternalManagement/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.InitConfig()
	
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	
	err = models.InitDefaultData(db)
	if err != nil {
		log.Fatalf("Failed to initialize default data: %v", err)
	}
	
	router := gin.Default()
	
	routes.SetupRoutes(router)
	
	log.Printf("Server starting on port %d...", cfg.Port)
	if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
