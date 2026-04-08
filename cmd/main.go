package main

import (
	"fmt"
	"log"

	"manage-product-service/internal/config"
	"manage-product-service/internal/database"
	"manage-product-service/internal/handler"
	"manage-product-service/internal/repository"
	"manage-product-service/internal/router"
	"manage-product-service/internal/service"
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
