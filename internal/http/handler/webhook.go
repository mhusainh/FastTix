package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
)

type WebhookHandler struct {
	purchaseService   service.PurchaseService
	submissionService service.SubmissionService // Updated to use SubmissionService
}

func NewWebhookHandler(purchaseService service.PurchaseService, submissionService service.SubmissionService) WebhookHandler {
	return WebhookHandler{purchaseService, submissionService}
}

func (h WebhookHandler) MidtransWebhook(ctx echo.Context) error {
	var notif map[string]interface{}

	// Bind body to payload
	if err := ctx.Bind(&notif); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid payload"))
	}

	// Extract order_id
	orderID, ok := notif["order_id"].(string)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Missing order_id"))
	}

	// Determine the type of transaction based on order_id prefix
	if strings.HasPrefix(orderID, "order_id-") {
		// Purchase Ticket Notification
		err := h.purchaseService.HandleMidtransNotification(ctx.Request().Context(), notif)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}
	} else if strings.HasPrefix(orderID, "daftar_id-") {
		// Submission Ticket Notification
		err := h.submissionService.HandleMidtransNotification(ctx.Request().Context(), notif)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}
	} else {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid order_id prefix"))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Notification handled successfully", nil))
}
