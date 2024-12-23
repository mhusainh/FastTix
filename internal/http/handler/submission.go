package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

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

func (h *SubmissionHandler) ApprovalSubmission(ctx echo.Context) error {
	var req dto.GetProductByIDRequest
	var n dto.CreateNotificationRequest
	var m dto.UpdateProductStatusRequest
	status := ctx.Param("status")
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

	user, err := h.userService.GetById(ctx.Request().Context(), Submission.UserID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	m.Status = status
	submission, err := h.submissionService.Approval(ctx.Request().Context(), m, Submission, user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	n.Message = submission.ProductStatus

	err = h.notificationService.SendNotificationSubmission(ctx.Request().Context(), n, submission)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully approval a submission", submission))
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

func saveUploadedFile(file *multipart.FileHeader, path string) error {
	// Open the uploaded file.
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create a destination file for the uploaded content.
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the uploaded content to the destination file.
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func (h *SubmissionHandler) UploadPicture(ctx echo.Context) error {
	var req dto.GetProductByIDRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Define the file path to save the uploaded image.
	pathImage := "images/" + file.Filename

	// Save the uploaded file to the specified path.
	if err := saveUploadedFile(file, pathImage); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	// Construct the URL for the saved picture.
	baseURL := "http://localhost:8080/api/v1"
	pictureURL := baseURL + "/image/" + file.Filename

	// Update the user's profile with the picture URL using the user service.
	if err := h.submissionService.UpdatePictureURL(ctx.Request().Context(), req, pictureURL); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("Successfully upload picture", pictureURL))
}
