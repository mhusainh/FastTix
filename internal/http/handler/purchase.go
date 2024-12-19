package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
)

type PurchaseHandler struct {
	purchaseService service.PurchaseService
	tokenService    service.TokenService
}

func NewPurchaseHandler(purchaseService service.PurchaseService, tokenService service.TokenService) PurchaseHandler {
	return PurchaseHandler{purchaseService, tokenService}
}

// PurchaseTicket handles the ticket purchasing process
func (h PurchaseHandler) PurchaseTicket(ctx echo.Context) error {
	var req dto.PurchaseTicketRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid request payload"))
	}

	// Get user ID from token
	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "Anda Harus Login Terlebih Dahulu"))
	}

	// Initiate the purchase
	purchaseResp, err := h.purchaseService.PurchaseTicket(ctx.Request().Context(), userID, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Purchase initiated", purchaseResp))
}

// CheckPurchaseStatus handles checking the status of a purchase
func (h PurchaseHandler) CheckPurchaseStatus(ctx echo.Context) error {
	orderID := ctx.Param("order_id")
	if orderID == "" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Order ID is required"))
	}

	status, err := h.purchaseService.CheckPurchaseStatus(ctx.Request().Context(), orderID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Purchase status fetched", status))
}
