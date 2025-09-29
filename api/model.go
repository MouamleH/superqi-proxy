package api

// ApplyTokenRequest represents the request body for apply token endpoint
type ApplyTokenRequest struct {
	AuthCode string `json:"authCode" validate:"required"`
}

// InquiryUserInfoRequest represents the request body for user info endpoint
type InquiryUserInfoRequest struct {
	AccessToken string `json:"accessToken" validate:"required"`
}

// InquiryUserCardListRequest represents the request body for user card list endpoint
type InquiryUserCardListRequest struct {
	AccessToken string `json:"accessToken" validate:"required"`
}

// PayRequest represents the request body for payment endpoint
type PayRequest struct {
	Amount      int    `json:"amount" validate:"required,min=1"`
	RequestId   string `json:"requestId" validate:"required"`
	AccessToken string `json:"accessToken" validate:"required"`
	CustomerId  string `json:"customerId" validate:"required"`
	OrderDesc   string `json:"orderDesc" validate:"required"`
	NotifyUrl   string `json:"notifyUrl" validate:"required,url"`
}

// PaymentStatusResponse represents the response for payment status endpoint (placeholder)
type PaymentStatusResponse struct {
	PaymentId string `json:"paymentId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
