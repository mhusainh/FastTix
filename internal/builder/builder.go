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

// BuilderPublicRoutes initializes public routes that do not require authentication
func BuilderPublicRoutes(cfg *config.Config, db *gorm.DB) []route.Route {
	// Repository Initialization
	productRepository := repository.NewProductRepository(db)
	userRepository := repository.NewUserRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	// submissionRepository := repository.NewSubmissionRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	// End Repository Initialization

	// Service Initialization
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	productService := service.NewProductService(productRepository, transactionRepository)
	userService := service.NewUserService(tokenService, cfg, userRepository)
	submissionService := service.NewSubmissionService(cfg, productRepository, userRepository, transactionRepository) // Updated to include cfg and repositories
	ticketService := service.NewTicketService(ticketRepository)
	purchaseService := service.NewPurchaseService(cfg, productRepository, userRepository, transactionRepository, tokenService)
	// End Service Initialization

	// Handler Initialization
	productHandler := handler.NewProductHandler(productService, tokenService)
	userHandler := handler.NewUserHandler(tokenService, userService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	purchaseHandler := handler.NewPurchaseHandler(purchaseService, tokenService)
	webhookHandler := handler.NewWebhookHandler(purchaseService, submissionService) // Pass SubmissionService instead of PengajuanService
	// End Handler Initialization

	return router.PublicRoutes(userHandler, productHandler, ticketHandler, webhookHandler, purchaseHandler)
}

// BuilderPrivateRoutes initializes private routes that require authentication
func BuilderPrivateRoutes(cfg *config.Config, db *gorm.DB) []route.Route {
	// Repository Initialization
	productRepository := repository.NewProductRepository(db)
	userRepository := repository.NewUserRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	// submissionRepository := repository.NewSubmissionRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	// End Repository Initialization

	// Service Initialization
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	productService := service.NewProductService(productRepository, transactionRepository)
	userService := service.NewUserService(tokenService, cfg, userRepository)
	submissionService := service.NewSubmissionService(cfg, productRepository, userRepository, transactionRepository)
	ticketService := service.NewTicketService(ticketRepository)
	purchaseService := service.NewPurchaseService(cfg, productRepository, userRepository, transactionRepository, tokenService)
	// End Service Initialization

	// Handler Initialization
	productHandler := handler.NewProductHandler(productService, tokenService)
	userHandler := handler.NewUserHandler(tokenService, userService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, tokenService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	purchaseHandler := handler.NewPurchaseHandler(purchaseService, tokenService)
	// webhookHandler is only in public routes
	// End Handler Initialization

	return router.PrivateRoutes(productHandler, userHandler, submissionHandler, ticketHandler, purchaseHandler)
}
