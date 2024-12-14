package router

import (
	"net/http"

	"github.com/mhusainh/FastTix/internal/http/handler"
	"github.com/mhusainh/FastTix/pkg/route"
)

func PublicRoutes(product handler.ProductHandler) []route.Route {
	return []route.Route{
		{
			Method: http.MethodGet,
			Path: "/Submissions",
			Handler: product.GetSubmissions,
		},
		{
			Method: http.MethodGet,
			Path: "/Submission",
			Handler: product.GetSubmission,
		},
		{
			Method: http.MethodGet,
			Path: "/products/:id",
			Handler: product.GetProduct,
		},
		{
			Method: http.MethodPost,
			Path: "/products",
			Handler: product.CreateProduct,
		},
		{
			Method: http.MethodPut,
			Path: "/products/:id",
			Handler: product.UpdateProduct,
		},
		{
			Method: http.MethodDelete,
			Path: "/products/:id",
			Handler: product.DeleteProduct,
		},
	}
}

func PrivateRoutes(product handler.ProductHandler) []route.Route{
	return []route.Route{
	}
}