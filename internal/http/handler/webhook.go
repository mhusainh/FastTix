package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
)

type WebhookHandler struct {
	paymentService        service.PaymentService
	submissionService     service.SubmissionService // Updated to use SubmissionService
	transactionRepository repository.TransactionRepository
}

func NewWebhookHandler(paymentService service.PaymentService, submissionService service.SubmissionService, transactionRepository repository.TransactionRepository) WebhookHandler {
	return WebhookHandler{paymentService, submissionService, transactionRepository}
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
		err := h.paymentService.HandleMidtransNotification(ctx.Request().Context(), notif)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}
	} else if strings.HasPrefix(orderID, "daftar_id-") {
		err := h.submissionService.HandleMidtransNotification(ctx.Request().Context(), notif)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}
	} else {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid order_id prefix"))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Notification handled successfully", nil))
}

func (h WebhookHandler) CheckinWebhook(ctx echo.Context) error {
	var req dto.CheckinWebhook
	OrderId := ctx.Param("order_id")
	fmt.Println(OrderId)
	req.OrderID = OrderId
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	GetTransaction, err := h.transactionRepository.GetByOrderID(ctx.Request().Context(), OrderId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	fmt.Println(GetTransaction.CheckIn)
	if GetTransaction.TransactionStatus != "success" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Transaction not success"))
	}
	if GetTransaction.CheckIn != 0 {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Tiket tidak valid atau sudah digunakan"))
	}
	req.TransactionStatus = GetTransaction.TransactionStatus
	err = h.paymentService.HandleCheckinNotification(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Anda telah berhasil melakukan checkin", nil))
}
