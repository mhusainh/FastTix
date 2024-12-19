package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
)

type TicketHandler struct {
	ticketService service.TicketService
}

func NewTicketHandler(ticketService service.TicketService) TicketHandler {
	return TicketHandler{ticketService}
}

func (h TicketHandler) GetTickets(ctx echo.Context) error {
	ticket, err := h.ticketService.GetAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing all tickets", ticket))
}

func (h TicketHandler) GetTicket(ctx echo.Context) error {
	var req dto.GetProductByIDRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	ticket, err := h.ticketService.GetById(ctx.Request().Context(), int(req.ID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing a ticket", ticket))
}
