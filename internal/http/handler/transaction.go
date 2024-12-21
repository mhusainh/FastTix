package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
)

type TransactionHandler struct {
	transactionService  service.TransactionService
	tokenService        service.TokenService
	userService         service.UserService
	productService      service.ProductService
	notificationService service.NotificationService
}

func NewTransactionHandler(
	transactionService service.TransactionService,
	tokenService service.TokenService,
	userService service.UserService,
	productService service.ProductService,
	notificationService service.NotificationService,
) TransactionHandler {
	return TransactionHandler{transactionService, tokenService, userService, productService, notificationService}
}

func (h *TransactionHandler) GetTransactions(ctx echo.Context) error {
	transactions, err := h.transactionService.GetAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing all transactions", transactions))
}

func (h *TransactionHandler) GetTransaction(ctx echo.Context) error {
	var req dto.GetTransactionByIDRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	transaction, err := h.transactionService.GetById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing a transaction", transaction))
}

func (h *TransactionHandler) GetTransactionByUserId(ctx echo.Context) error {
	var req dto.GetTransactionByUserIDRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	req.UserID = userID

	transactions, err := h.transactionService.GetByUserId(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing all user transactions", transactions))
}

func (h *TransactionHandler) CheckoutTicket(ctx echo.Context) error {
	var req dto.CreateTransactionRequest
	var n dto.CreateNotificationRequest

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

	product, err := h.productService.GetById(ctx.Request().Context(), req.ProductID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if product.ProductQuantity < req.TransactionQuantity {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Ticket yang tersedia tidak cukup"))
	}

	transaction, err := h.transactionService.Create(ctx.Request().Context(), req, user, product)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	n.Message = "Checkout Ticket"

	err = h.notificationService.SendNotificationTransaction(ctx.Request().Context(), n, product, transaction)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully create a transaction", transaction))
}

func (h *TransactionHandler) PaymentTicket(ctx echo.Context) error {
	var req dto.UpdateTransactionRequest
	var n dto.CreateNotificationRequest

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

	transaction, err := h.transactionService.GetById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	product, err := h.productService.GetById(ctx.Request().Context(), transaction.ProductID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if product.ProductQuantity < transaction.TransactionQuantity {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Ticket yang tersedia tidak cukup"))
	}

	product.ProductQuantity -= transaction.TransactionQuantity

	err = h.productService.Update(ctx.Request().Context(), product)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	transaction, err = h.transactionService.PaymentTicket(ctx.Request().Context(), req, user, product, transaction)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	n.Message = "Payment"

	err = h.notificationService.SendNotificationTransaction(ctx.Request().Context(), n, product, transaction)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully payment a ticket", transaction))
}

func (h *TransactionHandler) PaymentSubmission(ctx echo.Context) error {
	var req dto.UpdateTransactionRequest
	var n dto.CreateNotificationRequest

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

	transaction, err := h.transactionService.GetById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	product, err := h.productService.GetById(ctx.Request().Context(), transaction.ProductID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	product.ProductStatus = "pending"

	err = h.productService.Update(ctx.Request().Context(), product)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	transaction, err = h.transactionService.PaymentSubmission(ctx.Request().Context(), req, user, product, transaction)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	n.Message = "Payment"

	err = h.notificationService.SendNotificationTransaction(ctx.Request().Context(), n, product, transaction)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully payment a submission", transaction))
}
