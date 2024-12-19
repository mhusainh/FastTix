package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
)

type SubmissionHandler struct {
	submissionService service.SubmissionService
	tokenService      service.TokenService
}

func NewSubmissionHandler(submissionService service.SubmissionService, tokenService service.TokenService) SubmissionHandler {
	return SubmissionHandler{submissionService, tokenService}
}

// GetSubmissionDetail fetches submission details
func (h SubmissionHandler) GetSubmissionDetail(ctx echo.Context) error {
	detail := map[string]interface{}{
		"fields": []string{"product_name", "product_address", "product_time", "product_date", "product_price", "product_description", "product_category", "product_quantity"},
		"note":   "Isi seluruh field di atas untuk membuat pengajuan tiket.",
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Submission detail fetched", detail))
}

// CreateSubmission handles the creation of a new submission
func (h SubmissionHandler) CreateSubmission(ctx echo.Context) error {
	var req dto.CreateProductRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}
	req.UserID = userID

	redirectURL, err := h.submissionService.CreateSubmission(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	data := map[string]interface{}{
		"message":      "Submission created. If payment is required, check your email for payment confirmation.",
		"payment_link": redirectURL,
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Submission successful", data))
}

// ApproveSubmission allows admin to approve a submission
func (h SubmissionHandler) ApproveSubmission(ctx echo.Context) error {
	idStr := ctx.Param("id")
	pid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
	}

	// Optionally, verify admin role here

	err = h.submissionService.ApproveSubmission(ctx.Request().Context(), pid)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Submission approved", nil))
}

// RejectSubmission allows admin to reject a submission
func (h SubmissionHandler) RejectSubmission(ctx echo.Context) error {
	idStr := ctx.Param("id")
	pid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
	}

	err = h.submissionService.RejectSubmission(ctx.Request().Context(), pid)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Submission rejected", nil))
}

// MidtransWebhook handles Midtrans notifications for submissions
func (h SubmissionHandler) MidtransWebhook(ctx echo.Context) error {
	var notif map[string]interface{}
	if err := ctx.Bind(&notif); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid payload"))
	}

	err := h.submissionService.HandleMidtransNotification(ctx.Request().Context(), notif)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Notification handled", nil))
}
