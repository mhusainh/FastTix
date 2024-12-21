package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/repository"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
	"github.com/mhusainh/FastTix/utils"
)

type SubmissionHandler struct {
	submissionService     service.SubmissionService
	tokenService          service.TokenService
	paymentService        service.PaymentService
	transactionRepository repository.TransactionRepository
	userRepository        repository.UserRepository
}

func NewSubmissionHandler(submissionService service.SubmissionService,
	tokenService service.TokenService,
	paymentService service.PaymentService,
	transactionRepository repository.TransactionRepository,
	userRepository repository.UserRepository,
) SubmissionHandler {
	return SubmissionHandler{submissionService, tokenService, paymentService, transactionRepository, userRepository}
}

func (h SubmissionHandler) GetSubmissions(ctx echo.Context) error {
	var req dto.GetAllProductsRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	submission, err := h.submissionService.GetAll(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing all submission", submission))
}

func (h SubmissionHandler) GetSubmission(ctx echo.Context) error {
	var req dto.GetProductByIDRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	submission, err := h.submissionService.GetById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing a submission", submission))
}

// func (h SubmissionHandler) checkoutSubmission(ctx echo.Context) error {
// 	var req dto.CreateProductRequest
// 	var payment dto.CreatePaymentRequest

// 	if err := ctx.Bind(&req); err != nil {
// 		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
// 	}
// 	userID, err := h.tokenService.GetUserIDFromToken(ctx)
// 	if err != nil {
// 		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
// 	}

// 	req.UserID = userID
// 	// Buat transaksi hanya sekali
// 	// Siapkan pembayaran
// 	payment.OrderID = req.ProductName // Sesuaikan format OrderID sesuai kebutuhan
// 	payment.Amount = 1000             // Pastikan TransactionAmount diisi dengan benar
// 	payment.Email = "iiwandila01@gmail.com"
// 	payment.NameProduct = "FastTix" // Gantilah dengan email pengguna yang sesuai
// 	payment.UserID = userID

// 	purchase, err := h.paymentService.CreatePayment(ctx.Request().Context(), payment)
// 	if err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
// 	}
// 	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully created a transaction", purchase))
// }

func (h SubmissionHandler) CreateSubmission(ctx echo.Context) error {
	var t dto.CreateTransactionRequest
	var payment dto.CreatePaymentRequest
	token := ctx.Param("tokenid")
	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}
	if err := ctx.Bind(&t); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	Orderid, err := h.transactionRepository.GetTransactionByToken(ctx.Request().Context(), token)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	GetUserAll, err := h.userRepository.GetById(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	t.VerificationToken = Orderid.VerificationToken

	if token != t.VerificationToken {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid token"))
	}
	payment.UserID = userID
	payment.OrderID = Orderid.OrderID
	payment.Amount = Orderid.TransactionAmount
	payment.Email = GetUserAll.Email
	payment.NameProduct = t.NameProduct
	purchase, err := h.paymentService.CreatePayment(ctx.Request().Context(), payment)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully created a transaction", purchase))
}

func (h *SubmissionHandler) CheckoutSubmission(ctx echo.Context) error {
	var req dto.CreateProductRequest
	var t dto.CreateTransactionRequest
	var payment dto.CreatePaymentRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	t.VerificationToken = utils.RandomString(6)
	req.UserID = userID

	GetAll, err := h.userRepository.GetById(ctx.Request().Context(), userID)
	if GetAll == nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "User not found"))
	}

	err = h.submissionService.Create(ctx.Request().Context(), req, t)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	payment.UserID = userID
	payment.Email = GetAll.Email
	payment.VerificationToken = t.VerificationToken
	payment.NameProduct = req.ProductName // Gantilah dengan email pengguna yang sesuai
	// payment.UserID = userID

	purchase, err := h.paymentService.CreateTokenPayment(ctx.Request().Context(), payment)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully created a transaction", purchase))
}

func (h SubmissionHandler) UpdateSubmissionByUser(ctx echo.Context) error {
	var req dto.UpdateProductRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	req.UserID = userID

	err = h.submissionService.UpdateByUser(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully update a submission", nil))
}

func (h SubmissionHandler) ApproveSubmission(ctx echo.Context) error {
	var req dto.GetProductByIDRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err := h.submissionService.Approve(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully approve a submission", nil))
}

func (h SubmissionHandler) RejectSubmission(ctx echo.Context) error {
	var req dto.GetProductByIDRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err := h.submissionService.Reject(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully reject a submission", nil))
}

func (h SubmissionHandler) CancelSubmission(ctx echo.Context) error {
	var req dto.GetProductByIDRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	req.UserID = userID

	submission, err := h.submissionService.GetById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	err = h.submissionService.Cancel(ctx.Request().Context(), submission, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully cancel a submission", nil))
}
