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
	submission handler.SubmissionHandler,
	ticketHandler handler.TicketHandler,
	) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/products",
			Handler: productHandler.FilterProducts,
		},
		{
			Method:  http.MethodGet,
			Path:    "/products/sort",
			Handler: productHandler.SortProducts,
		},
		{
			Method: http.MethodGet,
			Path: "/submissions",
			Handler: submission.GetSubmissions,
		},
		{
			Method: http.MethodGet,
			Path: "/tickets",
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
	submission handler.SubmissionHandler,
	ticketHandler handler.TicketHandler,
	) []route.Route{
	return []route.Route{
		{
			Method: http.MethodGet,
			Path: "/submissions/:id",
			Handler: submission.GetSubmission,
		},
		{
			Method: http.MethodGet,
			Path: "/tickets/:id",
			Handler: ticketHandler.GetTicket,
		},
		{
			Method: http.MethodGet,
			Path: "/products/:id",
			Handler: productHandler.GetProduct,
		},
		{
			Method: http.MethodPost,
			Path: "/products",
			Handler: productHandler.CreateProduct,
		},
		{
			Method: http.MethodPut,
			Path: "/products/:id",
			Handler: productHandler.UpdateProduct,
		},
		{
			Method: http.MethodDelete,
			Path: "/products/:id",
			Handler: productHandler.DeleteProduct,
		},
	}
}