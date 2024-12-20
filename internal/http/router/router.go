package router

import (
	"net/http"

	"github.com/mhusainh/FastTix/internal/http/handler"
	"github.com/mhusainh/FastTix/pkg/route"
)

var (
	adminOnly = []string{"Administrator"}
	userOnly  = []string{"User"}
	allRoles  = []string{"Administrator", "User"}
)

func PublicRoutes(
	userHandler handler.UserHandler,
	productHandler handler.ProductHandler,
	submission handler.SubmissionHandler,
	ticketHandler handler.TicketHandler,
) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/submissions",
			Handler: submission.GetSubmissions,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tickets",
			Handler: ticketHandler.GetTickets,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: userHandler.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/request-reset-password",
			Handler: userHandler.ResetPasswordRequest,
		},
		{
			Method:  http.MethodPost,
			Path:    "/reset-password/:token",
			Handler: userHandler.ResetPassword,
		},
		{
			Method:  http.MethodGet,
			Path:    "/verify-email/:token",
			Handler: userHandler.VerifyEmail,
		},
	}
}

func PrivateRoutes(
	productHandler handler.ProductHandler,
	userHandler handler.UserHandler,
	submissionHandler handler.SubmissionHandler,
	ticketHandler handler.TicketHandler,
	transactionHandler handler.TransactionHandler,
) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/submissions/:id",
			Handler: submissionHandler.GetSubmission,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPost,
			Path:    "/submissions",
			Handler: submissionHandler.CreateSubmission,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodPut,
			Path:    "/submissions/:id",
			Handler: submissionHandler.UpdateSubmissionByUser,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/submissions/:id",
			Handler: submissionHandler.CancelSubmission,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodPut,
			Path:    "/submissions/:id/approve",
			Handler: submissionHandler.ApproveSubmission,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodPut,
			Path:    "/submissions/:id/reject",
			Handler: submissionHandler.RejectSubmission,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tickets/:id",
			Handler: ticketHandler.GetTicket,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/products/user",
			Handler: productHandler.GetProductByUserId,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/products",
			Handler: productHandler.GetProducts,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/products/:id",
			Handler: productHandler.GetProduct,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/products'////////////////:id",
			Handler: productHandler.DeleteProduct,
		},
		{
			Method:  http.MethodGet,
			Path:    "/transactions/user",
			Handler: transactionHandler.GetTransactionByUserId,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodPut,
			Path:    "/submissions/:id/payment",
			Handler: transactionHandler.PaymentSubmission,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/tickets/:product_id/checkout",
			Handler: transactionHandler.CheckoutTicket,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodPut,
			Path:    "/tickets/:id/payment",
			Handler: transactionHandler.PaymentTicket,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/transactions",
			Handler: transactionHandler.GetTransactions,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/transactions/:id",
			Handler: transactionHandler.GetTransaction,
			Roles:   adminOnly,
		},
	}
}
