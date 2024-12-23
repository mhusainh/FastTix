package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
	webhookHandler handler.WebhookHandler,
) []route.Route {
	return []route.Route{
		{
			Method: http.MethodGet,
			Path:   "/image/*",
			Handler: func(c echo.Context) error {
				filePath := c.Param("*")
				staticDir := "images/"
				fullPath := staticDir + filePath
				return c.File(fullPath)
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/webhook/midtrans",
			Handler: webhookHandler.MidtransWebhook,
		},
		{
			Method:  http.MethodGet,
			Path:    "/checkin/:order_id",
			Handler: webhookHandler.CheckinWebhook,
		},
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
	notificationHandler handler.NotificationHandler,
) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/submissions/:id",
			Handler: submissionHandler.GetSubmission,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/submissions/user",
			Handler: submissionHandler.GetSubmissionByUser,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/submissions",
			Handler: submissionHandler.CreateSubmission, // buat submission baru
			Roles:   userOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/payment/checkout/:tokenid", // payment checkout
			Handler: submissionHandler.CheckoutSubmission,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodPost,
			Path:    "/tickets/:product_id/checkout",
			Handler: transactionHandler.CheckoutTicket, // beli tiket
			Roles:   userOnly,
		},
		{
			Method:  http.MethodPut,
			Path:    "/submissions/:id",
			Handler: submissionHandler.UpdateSubmissionByUser,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodPut,
			Path:    "/submissions/:id/image",
			Handler: submissionHandler.UploadPicture,
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
			Path:    "/submissions/:id/:status",
			Handler: submissionHandler.ApprovalSubmission,
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
			Method:  http.MethodGet,
			Path:    "/transactions/user",
			Handler: transactionHandler.GetTransactionByUserId,
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
		{
			Method:  http.MethodGet,
			Path:    "/users/profile",
			Handler: userHandler.GetProfile,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodPut,
			Path:    "/users/profile",
			Handler: userHandler.UpdateUser,
			Roles:   allRoles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: userHandler.GetUsers,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:id",
			Handler: userHandler.GetUser,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/:id",
			Handler: userHandler.DeleteUser,
			Roles:   adminOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/notifications",
			Handler: notificationHandler.GetNotificationsByUser,
			Roles:   userOnly,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:id/notifications/:id",
			Handler: notificationHandler.GetNotificationsByUserAndID,
			Roles:   userOnly,
		},
	}
}
