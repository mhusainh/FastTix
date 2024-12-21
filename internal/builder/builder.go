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
	notificationRepository := repository.NewNotificationRepository(db)
	//end

	//service
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	productService := service.NewProductService(productRepository)
	userService := service.NewUserService(tokenService, cfg, userRepository)
	submissionService := service.NewSubmissionService(cfg, submissionRepository, transactionRepository, productRepository)
	ticketService := service.NewTicketService(ticketRepository)
	transactionService := service.NewTransactionService(cfg, transactionRepository, productRepository)
	notificationService := service.NewNotificationService(notificationRepository)
	//end

	//handler
	productHandler := handler.NewProductHandler(productService, tokenService)
	userHandler := handler.NewUserHandler(tokenService, userService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, tokenService, productService, transactionService, userService, notificationService)
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
	notificationRepository := repository.NewNotificationRepository(db)
	//end

	//service
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	productService := service.NewProductService(productRepository)
	userService := service.NewUserService(tokenService, cfg, userRepository)
	submissionService := service.NewSubmissionService(cfg, submissionRepository, transactionRepository, productRepository)
	ticketService := service.NewTicketService(ticketRepository)
	transactionService := service.NewTransactionService(cfg, transactionRepository,productRepository)
	notificationService := service.NewNotificationService(notificationRepository)
	//end

	//handler
	productHandler := handler.NewProductHandler(productService, tokenService)
	userHandler := handler.NewUserHandler(tokenService, userService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, tokenService, productService, transactionService, userService, notificationService)
	ticketHandler := handler.NewTicketHandler(ticketService)
	transactionHandler := handler.NewTransactionHandler(transactionService, tokenService, userService, productService, notificationService)
	notificationHandler := handler.NewNotificationHandler(notificationService, tokenService, userService)
	//end

	return router.PrivateRoutes(productHandler, userHandler, submissionHandler, ticketHandler, transactionHandler, notificationHandler)
}
