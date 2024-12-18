package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/pkg/response"
)

type WebhookHandler struct{}

func NewWebhookHandler() WebhookHandler {
	return WebhookHandler{}
}

func (h WebhookHandler) MidtransWebhook(ctx echo.Context) error {
	var payload map[string]interface{}

	// Bind body ke payload
	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid payload"))
	}

	// Cetak notifikasi untuk debug sementara
	println("Webhook Received:", payload)

	// Logika penanganan notifikasi Midtrans bisa ditambahkan disini
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Webhook diterima", nil))
}
