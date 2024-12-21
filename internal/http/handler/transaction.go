package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
	"github.com/mhusainh/FastTix/utils"
)

type TransactionHandler struct {
	transactionService    service.TransactionService
	tokenService          service.TokenService
	paymentService        service.PaymentService
	userRepository        repository.UserRepository
	transactionRepository repository.TransactionRepository
	productRepository     repository.ProductRepository
}

func NewTransactionHandler(
	transactionService service.TransactionService,
	tokenService service.TokenService,
	paymentService service.PaymentService,
	userRepository repository.UserRepository,
	transactionRepository repository.TransactionRepository,
	productRepository repository.ProductRepository,
) TransactionHandler {
	return TransactionHandler{transactionService, tokenService, paymentService, userRepository, transactionRepository, productRepository}
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
	var payment dto.CreatePaymentRequest
	idProduct := ctx.Param("product_id")
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	req.VerificationToken = utils.RandomString(6)
	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}
	GetUserAll, err := h.userRepository.GetById(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	idProductInt64, err := strconv.ParseInt(idProduct, 10, 64)
	Product, err := h.productRepository.GetById(ctx.Request().Context(), idProductInt64)
	if Product == nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Product not found"))
	}
	req.TransactionAmount = Product.ProductPrice * float64(req.TransactionQuantity)
	req.UserID = userID
	req.OrderID = fmt.Sprintf("order_id-%d", req.ProductID, time.Now().Unix())
	err = h.transactionService.Create(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	fmt.Println(req.TransactionAmount)
	payment.VerificationToken = req.VerificationToken // Sesuaikan format OrderID sesuai kebutuhan// Pastikan TransactionAmount diisi dengan benar
	payment.Email = GetUserAll.Email
	payment.NameProduct = req.NameProduct // Gantilah dengan email pengguna yang sesuai
	payment.UserID = userID

	// Buat pembayaran
	purchase, err := h.paymentService.CreateTokenPayment(ctx.Request().Context(), payment)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully created a transaction", purchase))
}

func (h *TransactionHandler) PaymentTicket(ctx echo.Context) error {
	var req dto.UpdateTransactionRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	req.UserID = userID

	err = h.transactionService.PaymentTicket(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully payment a ticket", nil))
}

func (h *TransactionHandler) PaymentSubmission(ctx echo.Context) error {
	var req dto.UpdateTransactionRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	req.UserID = userID

	err = h.transactionService.PaymentSubmission(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully payment a submission", nil))
}
