package builder

import (
	"github.com/mhusainh/FastTix/config"
	"github.com/mhusainh/FastTix/internal/http/handler"
	"github.com/mhusainh/FastTix/internal/http/router"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/route"
	"gorm.io/gorm"
)

func BuilderPublicRoutes(cfg *config.Config, db *gorm.DB) []route.Route {
	//repository
	productRepository := repository.NewProductRepository(db)
	//end

	//service
	productService := service.NewProductService(productRepository)
	//end

	//handler
	productHandler := handler.NewProductHandler(productService)
	//end

	return router.PublicRoutes(productHandler)
}

func BuilderPrivateRoutes(cfg *config.Config, db *gorm.DB) []route.Route {
	//repository
	_ = repository.NewProductRepository(db)
	//end

	//service
	// _ = service.NewProductService()
	//end

	//handler
	// _ = handler.NewProductHandler()
	//end

	return nil
}