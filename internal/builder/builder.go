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
	userRepository := repository.NewUserRepository(db)
	//end

	//service
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	productService := service.NewProductService(productRepository)
	userService := service.NewUserService(tokenService, cfg, userRepository)
	//end

	//handler
	productHandler := handler.NewProductHandler(productService)
	userHandler := handler.NewUserHandler(tokenService, userService)
	//end

	return router.PublicRoutes(userHandler, productHandler)
}

func BuilderPrivateRoutes(cfg *config.Config, db *gorm.DB) []route.Route {
	//repository
	productRepository := repository.NewProductRepository(db)
	userRepository := repository.NewUserRepository(db)
	//end

	//service
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	productService := service.NewProductService(productRepository)
	userService := service.NewUserService(tokenService, cfg, userRepository)
	//end

	//handler
	productHandler := handler.NewProductHandler(productService)
	userHandler := handler.NewUserHandler(tokenService, userService)
	//end

	return router.PrivateRoutes(productHandler, userHandler)
}