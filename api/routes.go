package api

import (
	"github.com/gofiber/fiber/v2"
)

// InitAPI initializes all API routes
func InitAPI(api fiber.Router) {
	// Version 1 API routes
	v1 := api.Group("/v1")

	// Authentication endpoints
	v1.Post("/apply-token", ApplyToken)

	// User information endpoints
	v1.Post("/user-info", InquiryUserInfo)
	v1.Post("/user-cards", InquiryUserCardList)

	// Payment endpoints
	v1.Post("/pay", Pay)
	v1.Get("/payment/:paymentId/status", GetPaymentStatus)

	// Health check
	v1.Get("/health", Health)
}
