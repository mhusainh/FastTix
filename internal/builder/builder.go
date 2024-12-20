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
	transactionRepository := repository.NewTransactionRepository(db)
	submissionRepository := repository.NewSubmissionRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	//end

	//service
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	productService := service.NewProductService(productRepository)
	userService := service.NewUserService(tokenService, cfg, userRepository)
	submissionService := service.NewSubmissionService(cfg, submissionRepository, transactionRepository, userRepository, productRepository)
	ticketService := service.NewTicketService(ticketRepository)
	//end

	//handler
	productHandler := handler.NewProductHandler(productService, tokenService)
	userHandler := handler.NewUserHandler(tokenService, userService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, tokenService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	//end

	return router.PublicRoutes(userHandler, productHandler, submissionHandler, ticketHandler)
}

func BuilderPrivateRoutes(cfg *config.Config, db *gorm.DB) []route.Route {
	//repository
	productRepository := repository.NewProductRepository(db)
	userRepository := repository.NewUserRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	submissionRepository := repository.NewSubmissionRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	transactionRepository = repository.NewTransactionRepository(db)
	//end

	//service
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	productService := service.NewProductService(productRepository)
	userService := service.NewUserService(tokenService, cfg, userRepository)
	submissionService := service.NewSubmissionService(cfg, submissionRepository, transactionRepository, userRepository, productRepository)
	ticketService := service.NewTicketService(ticketRepository)
	transactionService := service.NewTransactionService(transactionRepository)
	//end

	//handler
	productHandler := handler.NewProductHandler(productService, tokenService)
	userHandler := handler.NewUserHandler(tokenService, userService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, tokenService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	transactionHandler := handler.NewTransactionHandler(transactionService, tokenService)
	//end

	return router.PrivateRoutes(productHandler, userHandler, submissionHandler, ticketHandler, transactionHandler)
}
