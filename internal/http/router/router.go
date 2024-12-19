package router

import (
	"net/http"

	"github.com/mhusainh/FastTix/internal/http/handler"
	"github.com/mhusainh/FastTix/pkg/route"
)

var (
	adminOnly = []string{"Administrator"}

	allRoles = []string{"Administrator", "User"}
)

func PublicRoutes(
	userHandler handler.UserHandler,
	productHandler handler.ProductHandler,
	// submissionHandler handler.SubmissionHandler,
	ticketHandler handler.TicketHandler,
	webhookHandler handler.WebhookHandler,
	purchaseHandler handler.PurchaseHandler,
) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/webhook/midtrans",
			Handler: webhookHandler.MidtransWebhook,
		},
		// {
		// 	Method:  http.MethodGet,
		// 	Path:    "/submissions",
		// 	Handler: submissionHandler.GetSubmissionDetail,
		// },
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
	purchaseHandler handler.PurchaseHandler,
) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/purchase",
			Handler: purchaseHandler.PurchaseTicket,
		},
		{
			Method:  http.MethodGet,
			Path:    "/purchase/status/:order_id",
			Handler: purchaseHandler.CheckPurchaseStatus,
		},
		{
			Method:  http.MethodGet,
			Path:    "/submissions/:id",
			Handler: submissionHandler.GetSubmissionDetail,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tickets/:id",
			Handler: ticketHandler.GetTicket,
		},
		{
			Method:  http.MethodGet,
			Path:    "/products/:id",
			Handler: productHandler.GetProduct,
		},
		{
			Method:  http.MethodGet,
			Path:    "/products",
			Handler: productHandler.GetProducts,
		},
		{
			Method:  http.MethodPost,
			Path:    "/create/products",
			Handler: productHandler.CreateProduct,
		},
		{
			Method:  http.MethodPut,
			Path:    "/products/:id",
			Handler: productHandler.UpdateProduct,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/products/:id",
			Handler: productHandler.DeleteProduct,
		},
		{
			Method:  http.MethodGet,
			Path:    "/submissions/detail",
			Handler: submissionHandler.GetSubmissionDetail,
		},
		{
			Method:  http.MethodGet,
			Path:    "/submissions",
			Handler: productHandler.GetStatusPending,
		},
		{
			Method:  http.MethodPost,
			Path:    "/create/submissions",
			Handler: submissionHandler.CreateSubmission,
		},
		{
			Method:  http.MethodPost,
			Path:    "/submissions/approve/:id",
			Handler: submissionHandler.ApproveSubmission,
		},
		{
			Method:  http.MethodPost,
			Path:    "/submissions/reject/:id",
			Handler: submissionHandler.RejectSubmission,
		},
	}
}
