package main

import (
	"fmt"
	"log"

	"prm-product/internal/config"
	"prm-product/internal/database"
	"prm-product/internal/handler"
	"prm-product/internal/repository"
	"prm-product/internal/router"
	"prm-product/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	app := router.NewRouter(productHandler)
	if err := app.Listen(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
