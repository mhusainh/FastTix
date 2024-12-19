package router

import (
	"net/http"

	"github.com/mhusainh/FastTix/internal/http/handler"
	"github.com/mhusainh/FastTix/pkg/route"
)

func PublicRoutes(product *handler.ProductHandler) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/submissions",
			Handler: product.GetSubmissions,
		},
		{
			Method:  http.MethodGet,
			Path:    "/submission",
			Handler: product.GetSubmission,
		},
		{
			Method:  http.MethodGet,
			Path:    "/verify-submission/:token",
			Handler: product.VerifySubmission,
		},
		{
			Method:  http.MethodGet,
			Path:    "/products",
			Handler: product.GetProducts,
		},
		{
			Method:  http.MethodGet,
			Path:    "/products/:id",
			Handler: product.GetProduct,
		},
		{
			Method:  http.MethodPost,
			Path:    "/products",
			Handler: product.CreateProduct,
		},
		{
			Method:  http.MethodPut,
			Path:    "/products/:id",
			Handler: product.UpdateProduct,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/products/:id",
			Handler: product.DeleteProduct,
		},

		{
			Method:  http.MethodGet,
			Path:    "/search",
			Handler: product.SearchProduct,
		},

		{
			Method:  http.MethodGet,
			Path:    "/filter/address",
			Handler: product.FilterProductsByAddress,
		},

		{
			Method:  http.MethodGet,
			Path:    "/filter/category",
			Handler: product.FilterProductsByCategory,
		},

		{
			Method:  http.MethodGet,
			Path:    "/filter/price",
			Handler: product.FilterProductsByPrice,
		},

		{
			Method:  http.MethodGet,
			Path:    "/filter/time",
			Handler: product.FilterProductsByTime,
		},

		{
			Method:  http.MethodGet,
			Path:    "/filter/date",
			Handler: product.FilterProductsByDate,
		},

		{
			Method:  http.MethodGet,
			Path:    "/filter/status",
			Handler: product.FilterProductsByStatus,
		},

		{
			Method:  http.MethodGet,
			Path:    "/sort/newest",
			Handler: product.SortProductByNewest,
		},

		{
			Method:  http.MethodGet,
			Path:    "/sort/expensive",
			Handler: product.SortProductByExpensive,
		},

		{
			Method:  http.MethodGet,
			Path:    "/sort/cheapest",
			Handler: product.SortProductByCheapest,
		},

		{
			Method:  http.MethodGet,
			Path:    "/sort/mostbought",
			Handler: product.SortProductByMostBought,
		},

		{
			Method:  http.MethodGet,
			Path:    "/sort/available",
			Handler: product.SortProductByAvailable,
		},
	}
}

func PrivateRoutes(product handler.ProductHandler) []route.Route {
	return []route.Route{}
}
