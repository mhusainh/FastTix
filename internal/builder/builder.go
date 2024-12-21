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
	paymentService := service.NewPaymentService(nil, cfg, userRepository, transactionRepository, productRepository)
	//end

	//handler
	productHandler := handler.NewProductHandler(productService, tokenService)
	userHandler := handler.NewUserHandler(tokenService, userService)

	submissionHandler := handler.NewSubmissionHandler(submissionService, tokenService, paymentService, transactionRepository, userRepository)
	ticketHandler := handler.NewTicketHandler(ticketService)
	webhookHandler := handler.NewWebhookHandler(paymentService, submissionService, transactionRepository)
	//end

	return router.PublicRoutes(userHandler, productHandler, submissionHandler, ticketHandler, webhookHandler)
}

func BuilderPrivateRoutes(cfg *config.Config, db *gorm.DB) []route.Route {
	//repository
	productRepository := repository.NewProductRepository(db)
	userRepository := repository.NewUserRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	submissionRepository := repository.NewSubmissionRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	transactionRepository = repository.NewTransactionRepository(db)
	paymentRepository := repository.NewPaymentRequestRepository(db)
	//end

	//service
	tokenService := service.NewTokenService(cfg.JWTConfig.SecretKey)
	productService := service.NewProductService(productRepository)
	userService := service.NewUserService(tokenService, cfg, userRepository)
	submissionService := service.NewSubmissionService(cfg, submissionRepository, transactionRepository, userRepository, productRepository)
	ticketService := service.NewTicketService(ticketRepository)
	transactionService := service.NewTransactionService(cfg, transactionRepository, userRepository, productRepository)
	paymentService := service.NewPaymentService(paymentRepository, cfg, userRepository, transactionRepository, productRepository)
	//end

	//handler
	productHandler := handler.NewProductHandler(productService, tokenService)
	userHandler := handler.NewUserHandler(tokenService, userService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, tokenService, paymentService, transactionRepository, userRepository)
	ticketHandler := handler.NewTicketHandler(ticketService)
	transactionHandler := handler.NewTransactionHandler(transactionService, tokenService, paymentService, userRepository, transactionRepository, productRepository)
	paymentHandler := handler.NewPaymentHandler(paymentService, tokenService)
	//end

	return router.PrivateRoutes(productHandler, userHandler, submissionHandler, ticketHandler, transactionHandler, paymentHandler)
}
