package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
)

type NotificationHandler struct {
	notificationService service.NotificationService
	tokenService        service.TokenService
	userService         service.UserService
}

func NewNotificationHandler(
	notificationService service.NotificationService,
	tokenService service.TokenService,
	userService service.UserService,
) NotificationHandler {
	return NotificationHandler{notificationService, tokenService, userService}
}

func (h *NotificationHandler) GetNotificationsByUser(ctx echo.Context) error {
	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	notifications, err := h.notificationService.GetByUserID(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing all notifications", notifications))
}

func (h *NotificationHandler) GetNotificationsByUserAndID(ctx echo.Context) error {
	var req dto.GetNotificationByIDRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	user, err := h.userService.GetById(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if userID != user.ID {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "Anda tidak memiliki hak untuk melihat notifikasi ini"))
	}

	notification, err := h.notificationService.GetByID(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	err = h.notificationService.MarkAsRead(ctx.Request().Context(), notification)

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing a notifications", notification))
}
