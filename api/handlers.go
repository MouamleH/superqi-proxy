package api

import (
	"superqi-proxy/superqi"

	"github.com/gofiber/fiber/v2"
)

// ApplyToken handles the apply token endpoint
func ApplyToken(c *fiber.Ctx) error {
	var req ApplyTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if req.AuthCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "authCode is required",
		})
	}

	response, err := superqi.Interface.ApplyToken(req.AuthCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(response)
}

// InquiryUserInfo handles the user info inquiry endpoint
func InquiryUserInfo(c *fiber.Ctx) error {
	var req InquiryUserInfoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if req.AccessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "accessToken is required",
		})
	}

	response, err := superqi.Interface.InquiryUserInfo(req.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(response)
}

// InquiryUserCardList handles the user card list inquiry endpoint
func InquiryUserCardList(c *fiber.Ctx) error {
	var req InquiryUserCardListRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if req.AccessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "accessToken is required",
		})
	}

	response, err := superqi.Interface.InquiryUserCardList(req.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(response)
}

// Pay handles the payment endpoint
func Pay(c *fiber.Ctx) error {
	var req PayRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	// Validate required fields
	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "amount must be greater than 0",
		})
	}
	if req.RequestId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "requestId is required",
		})
	}
	if req.AccessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "accessToken is required",
		})
	}
	if req.CustomerId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "customerId is required",
		})
	}
	if req.OrderDesc == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "orderDesc is required",
		})
	}
	if req.NotifyUrl == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "notifyUrl is required",
		})
	}

	response, err := superqi.Interface.Pay(
		req.Amount,
		req.RequestId,
		req.AccessToken,
		req.CustomerId,
		req.OrderDesc,
		req.NotifyUrl,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(response)
}

// GetPaymentStatus handles payment status inquiry by payment ID
func GetPaymentStatus(c *fiber.Ctx) error {
	paymentId := c.Params("paymentId")
	if paymentId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "paymentId is required",
		})
	}

	// This would require implementing a payment status inquiry method in superqi
	// For now, return a placeholder response
	return c.Status(fiber.StatusNotImplemented).JSON(PaymentStatusResponse{
		PaymentId: paymentId,
		Status:    "not_implemented",
		Message:   "Payment status inquiry not yet implemented",
	})
}

// Health check endpoint
func Health(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{
		Status:  "ok",
		Service: "superqi-proxy",
	})
}
