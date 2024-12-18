package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
)

type PengajuanHandler struct {
	pengajuanService service.PengajuanService
	tokenService     service.TokenService
}

func NewPengajuanHandler(pengajuanService service.PengajuanService, tokenService service.TokenService) PengajuanHandler {
	return PengajuanHandler{pengajuanService, tokenService}
}

// Detail pengajuan
func (h PengajuanHandler) GetPengajuanDetail(ctx echo.Context) error {
	detail := map[string]interface{}{
		"fields": []string{"product_name", "product_address", "product_time", "product_date", "product_price", "product_description", "product_category", "product_quantity"},
		"note":   "Isi seluruh field di atas untuk membuat pengajuan tiket.",
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Pengajuan detail fetched", detail))
}

// Buat pengajuan
func (h PengajuanHandler) CreatePengajuan(ctx echo.Context) error {
	var req dto.CreateProductRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}
	req.UserID = userID

	redirectURL, err := h.pengajuanService.CreatePengajuan(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	data := map[string]interface{}{
		"message":      "Pengajuan dibuat. Jika berbayar, cek email untuk konfirmasi pembayaran.",
		"payment_link": redirectURL,
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Pengajuan berhasil", data))
}

// Midtrans notification
// internal/http/handler/pengajuan.go

func (h PengajuanHandler) MidtransNotification(ctx echo.Context) error {
	var notif map[string]interface{}
	if err := ctx.Bind(&notif); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err := h.pengajuanService.HandleMidtransNotification(ctx.Request().Context(), notif)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("notification handled", nil))
}

// Admin approve
func (h PengajuanHandler) ApprovePengajuan(ctx echo.Context) error {
	idStr := ctx.Param("id")
	pid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
	}

	// Cek role admin jika diperlukan

	err = h.pengajuanService.ApprovePengajuan(ctx.Request().Context(), pid)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Pengajuan disetujui", nil))
}

// Admin reject
func (h PengajuanHandler) RejectPengajuan(ctx echo.Context) error {
	idStr := ctx.Param("id")
	pid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid ID"))
	}

	err = h.pengajuanService.RejectPengajuan(ctx.Request().Context(), pid)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Pengajuan ditolak", nil))
}
