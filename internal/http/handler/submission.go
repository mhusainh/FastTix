package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mhusainh/FastTix/internal/http/dto"
	"github.com/mhusainh/FastTix/internal/service"
	"github.com/mhusainh/FastTix/pkg/response"
	"github.com/mhusainh/FastTix/utils"
)

type SubmissionHandler struct {
	submissionService   service.SubmissionService
	tokenService        service.TokenService
	productService      service.ProductService
	transactionService  service.TransactionService
	userService         service.UserService
	notificationService service.NotificationService
	paymentService      service.PaymentService
}

func NewSubmissionHandler(
	submissionService service.SubmissionService,
	tokenService service.TokenService,
	productService service.ProductService,
	transactionService service.TransactionService,
	userService service.UserService,
	notificationService service.NotificationService,
	paymentService service.PaymentService,
) SubmissionHandler {
	return SubmissionHandler{submissionService, tokenService, productService, transactionService, userService, notificationService, paymentService}
}

func (h *SubmissionHandler) GetSubmissions(ctx echo.Context) error {
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

func (h *SubmissionHandler) GetSubmission(ctx echo.Context) error {
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

func (h *SubmissionHandler) GetSubmissionByUser(ctx echo.Context) error {
	var req dto.GetProductByUserIDRequest

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

	submission, err := h.submissionService.GetByUserId(ctx.Request().Context(), req, user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully showing all user's submission", submission))
}

func (h *SubmissionHandler) CreateSubmission(ctx echo.Context) error {
	var req dto.CreateProductRequest
	var t dto.CreateTransactionRequest
	var n dto.CreateNotificationRequest

	if err := ctx.Bind(&req); err != nil {
		fmt.Println("Kontol", err.Error())
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	fmt.Println(req.ProductName)
	userID, err := h.tokenService.GetUserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	user, err := h.userService.GetById(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	t.VerificationToken = utils.RandomString(6)

	submission, err := h.submissionService.Create(ctx.Request().Context(), req, t, user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	fmt.Println(submission.OrderID)
	if submission.ProductPrice != 0 {
		transaction, err := h.transactionService.GetByOrderID(ctx.Request().Context(), submission.OrderID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}

		err = h.paymentService.CreateTokenPayment(ctx.Request().Context(), user, submission, transaction)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
		}
	}

	n.Message = "create"

	err = h.notificationService.SendNotificationSubmission(ctx.Request().Context(), n, submission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully create a product", nil))
}

func (h *SubmissionHandler) CheckoutSubmission(ctx echo.Context) error {

	verificationToken := ctx.Param("tokenid")
	var t dto.GetTransactionByVerificationTokenRequest
	t.VerificationToken = verificationToken

	if err := ctx.Bind(&t); err != nil {
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

	transaction, err := h.transactionService.GetTransactionByToken(ctx.Request().Context(), t.VerificationToken)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	product, err := h.productService.GetById(ctx.Request().Context(), transaction.ProductID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if transaction.UserID != user.ID {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "Anda tidak memiliki hak untuk melihat transaksi ini"))
	}

	purchase, err := h.paymentService.CreatePayment(ctx.Request().Context(), product, user, transaction)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully generate payment url", purchase))
}

func (h *SubmissionHandler) UpdateSubmissionByUser(ctx echo.Context) error {
	var req dto.UpdateProductRequest
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

	submission, err := h.submissionService.GetById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if submission.UserID != user.ID {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "Anda tidak memiliki hak untuk mengupdate pengajuan ini"))
	}

	if submission.ProductStatus != "pending" && submission.ProductStatus != "unpaid" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Pengajuan ini sudah tidak dapat diupdate karena sudah diterima atau ditolak oleh admin"))
	}

	submission, err = h.submissionService.UpdateByUser(ctx.Request().Context(), req, user, submission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	n.Message = "update"

	err = h.notificationService.SendNotificationSubmission(ctx.Request().Context(), n, submission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully update a submission", submission))
}

func (h *SubmissionHandler) ApproveSubmission(ctx echo.Context) error {
	var req dto.GetProductByIDRequest
	var n dto.CreateNotificationRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	Submission, err := h.submissionService.GetById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if Submission.ProductStatus != "pending" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Pengajuan ini sudah tidak dapat diapprove karena sudah diterima atau ditolak oleh admin"))
	}

	submission, err := h.submissionService.Approve(ctx.Request().Context(), Submission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	n.Message = "accept"

	err = h.notificationService.SendNotificationSubmission(ctx.Request().Context(), n, submission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully approve a submission", submission))
}

func (h *SubmissionHandler) RejectSubmission(ctx echo.Context) error {
	var req dto.GetProductByIDRequest
	var n dto.CreateNotificationRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	submission, err := h.submissionService.GetById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if submission.ProductStatus != "pending" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Pengajuan ini sudah tidak dapat direject karena sudah diterima atau ditolak oleh admin"))
	}

	submission, err = h.submissionService.Reject(ctx.Request().Context(), submission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	n.Message = "reject"

	err = h.notificationService.SendNotificationSubmission(ctx.Request().Context(), n, submission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully reject a submission", submission))
}

func (h SubmissionHandler) CancelSubmission(ctx echo.Context) error {
	var req dto.GetProductByIDRequest
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

	submission, err := h.submissionService.GetById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	if user.ID != submission.UserID {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "Anda tidak memiliki hak untuk membatalkan pengajuan ini"))
	}

	if submission.ProductStatus != "pending" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Pengajuan ini sudah tidak dapat dicancel karena sudah diterima atau ditolak oleh admin"))
	}

	submission, err = h.submissionService.GetById(ctx.Request().Context(), req.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	err = h.submissionService.Cancel(ctx.Request().Context(), submission, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	n.Message = "delete"

	err = h.notificationService.SendNotificationSubmission(ctx.Request().Context(), n, submission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully cancel a submission", nil))
}
