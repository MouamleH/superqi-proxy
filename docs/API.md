# SuperQi Proxy API Documentation

## Overview

The SuperQi Proxy API provides HTTP endpoints to interact with the SuperQi payment gateway. This API serves as a proxy layer that wraps the SuperQi functionality and provides RESTful endpoints for payment processing, user management, and authentication.

**Base URL:** `http://localhost:1999/api/v1`

**Content-Type:** `application/json`

## Authentication Flow

The SuperQi API uses an OAuth2-like flow with authorization codes and access tokens:

1. **Authorization Code** → Exchange for Access Token via `/apply-token`
2. **Access Token** → Use for authenticated requests (user info, payments, etc.)
3. **Refresh Token** → Use to refresh expired access tokens (if needed)

## Endpoints

### 1. Apply Token

Exchange an authorization code for access and refresh tokens.

**Endpoint:** `POST /api/v1/apply-token`

**Request Body:**
```json
{
  "authCode": "string"
}
```

**Request Fields:**
- `authCode` (string, required): The authorization code obtained from SuperQi OAuth flow

**Success Response (200):**
```json
{
  "result": {
    "resultCode": "SUCCESS",
    "resultStatus": "S",
    "resultMessage": "Success"
  },
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "accessTokenExpiryTime": "2024-12-31T23:59:59Z",
  "refreshToken": "refresh_token_here",
  "refreshTokenExpiryTime": "2025-01-30T23:59:59Z",
  "customerId": "customer_123456"
}
```

**Error Response (400/500):**
```json
{
  "error": "authCode is required"
}
```

---

### 2. Get User Information

Retrieve detailed user information using an access token.

**Endpoint:** `POST /api/v1/user-info`

**Request Body:**
```json
{
  "accessToken": "string"
}
```

**Request Fields:**
- `accessToken` (string, required): Valid access token from `/apply-token`

**Success Response (200):**
```json
{
  "result": {
    "resultCode": "SUCCESS",
    "resultStatus": "S",
    "resultMessage": "Success"
  },
  "userInfo": {
    "userId": "user_123456",
    "loginIdInfos": [
      {
        "loginId": "user@example.com",
        "hashLoginId": "hashed_value",
        "maskLoginId": "u***@example.com",
        "loginIdType": "EMAIL"
      }
    ],
    "userName": {
      "fullName": "John Doe Smith",
      "firstName": "John",
      "secondName": "Doe",
      "thirdName": "",
      "lastName": "Smith"
    },
    "userNameInArabic": {
      "fullName": "جون دو سميث",
      "firstName": "جون",
      "secondName": "دو",
      "thirdName": "",
      "lastName": "سميث"
    },
    "avatar": "https://example.com/avatar.jpg",
    "gender": "MALE",
    "birthDate": "1990-01-15",
    "nationality": "IQ",
    "contactInfos": [
      {
        "contactType": "MOBILE",
        "contactNo": "+9647901234567"
      }
    ]
  }
}
```

**Error Response (400/500):**
```json
{
  "error": "accessToken is required"
}
```

---

### 3. Get User Card List

Retrieve the user's registered payment cards.

**Endpoint:** `POST /api/v1/user-cards`

**Request Body:**
```json
{
  "accessToken": "string"
}
```

**Request Fields:**
- `accessToken` (string, required): Valid access token from `/apply-token`

**Success Response (200):**
```json
{
  "result": {
    "resultCode": "SUCCESS",
    "resultStatus": "S",
    "resultMessage": "Success"
  },
  "cardList": [
    {
      "maskedCardNo": "****-****-****-1234",
      "accountNumber": "ACC123456789"
    },
    {
      "maskedCardNo": "****-****-****-5678",
      "accountNumber": "ACC987654321"
    }
  ]
}
```

**Error Response (400/500):**
```json
{
  "error": "accessToken is required"
}
```

---

### 4. Process Payment

Initiate a payment transaction.

**Endpoint:** `POST /api/v1/pay`

**Request Body:**
```json
{
  "amount": 1000,
  "requestId": "unique-request-id-123",
  "accessToken": "string",
  "customerId": "customer_123456",
  "orderDesc": "Payment for Order #12345",
  "notifyUrl": "https://your-domain.com/webhook/payment"
}
```

**Request Fields:**
- `amount` (integer, required): Payment amount in Iraqi Dinar (IQD). Must be > 0
- `requestId` (string, required): Unique identifier for this payment request
- `accessToken` (string, required): Valid access token from `/apply-token`
- `customerId` (string, required): Customer identifier
- `orderDesc` (string, required): Description of the order/payment
- `notifyUrl` (string, required): URL to receive payment notifications (must be valid URL)

**Success Response (200):**
```json
{
  "paymentId": "pay_123456789",
  "result": {
    "resultCode": "SUCCESS",
    "resultStatus": "S",
    "resultMessage": "Payment initiated successfully"
  },
  "redirectActionForm": {
    "method": "POST",
    "parameters": "encoded_form_parameters",
    "redirectUrl": "https://payment-gateway.superqi.com/checkout"
  }
}
```

**Error Response (400/500):**
```json
{
  "error": "amount must be greater than 0"
}
```

**Notes:**
- Amount is automatically multiplied by 1000 internally (SuperQi requirement)
- The `redirectActionForm` contains information needed to redirect the user to complete payment
- `requestId` must be unique for each payment attempt

---

### 5. Get Payment Status

Retrieve the status of a payment transaction.

**Endpoint:** `GET /api/v1/payment/{paymentId}/status`

**Path Parameters:**
- `paymentId` (string, required): The payment ID returned from `/pay` endpoint

**Success Response (501 - Not Implemented):**
```json
{
  "paymentId": "pay_123456789",
  "status": "not_implemented",
  "message": "Payment status inquiry not yet implemented"
}
```

**Error Response (400):**
```json
{
  "error": "paymentId is required"
}
```

**Note:** This endpoint is currently a placeholder and returns HTTP 501 Not Implemented.

---

### 6. Health Check

Check the health status of the API service.

**Endpoint:** `GET /api/v1/health`

**Success Response (200):**
```json
{
  "status": "ok",
  "service": "superqi-proxy"
}
```

## Error Handling

### HTTP Status Codes

- **200 OK** - Request successful
- **400 Bad Request** - Invalid request body or missing required fields
- **500 Internal Server Error** - SuperQi API error or server error
- **501 Not Implemented** - Endpoint not yet implemented

### Error Response Format

All error responses follow this format:
```json
{
  "error": "Descriptive error message"
}
```

### Common Error Messages

- `"Invalid request body"` - JSON parsing failed
- `"authCode is required"` - Missing required authCode field
- `"accessToken is required"` - Missing required accessToken field
- `"amount must be greater than 0"` - Invalid amount value
- `"paymentId is required"` - Missing paymentId parameter

## Usage Examples

### cURL Examples

#### 1. Apply Token
```bash
curl -X POST http://localhost:1999/api/v1/apply-token \
  -H "Content-Type: application/json" \
  -d '{
    "authCode": "AUTH_CODE_FROM_SUPERQI"
  }'
```

#### 2. Get User Information
```bash
curl -X POST http://localhost:1999/api/v1/user-info \
  -H "Content-Type: application/json" \
  -d '{
    "accessToken": "YOUR_ACCESS_TOKEN"
  }'
```

#### 3. Get User Cards
```bash
curl -X POST http://localhost:1999/api/v1/user-cards \
  -H "Content-Type: application/json" \
  -d '{
    "accessToken": "YOUR_ACCESS_TOKEN"
  }'
```

#### 4. Process Payment
```bash
curl -X POST http://localhost:1999/api/v1/pay \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000,
    "requestId": "req_' $(date +%s) '",
    "accessToken": "YOUR_ACCESS_TOKEN",
    "customerId": "customer_123",
    "orderDesc": "Test payment for order #12345",
    "notifyUrl": "https://your-domain.com/webhook/payment"
  }'
```

#### 5. Check Payment Status
```bash
curl -X GET http://localhost:1999/api/v1/payment/pay_123456789/status
```

#### 6. Health Check
```bash
curl -X GET http://localhost:1999/api/v1/health
```

### JavaScript/Node.js Examples

```javascript
const API_BASE = 'http://localhost:1999/api/v1';

// Apply Token
async function applyToken(authCode) {
  const response = await fetch(`${API_BASE}/apply-token`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ authCode })
  });
  return response.json();
}

// Get User Info
async function getUserInfo(accessToken) {
  const response = await fetch(`${API_BASE}/user-info`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ accessToken })
  });
  return response.json();
}

// Process Payment
async function processPayment(paymentData) {
  const response = await fetch(`${API_BASE}/pay`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(paymentData)
  });
  return response.json();
}
```

### Python Examples

```python
import requests
import json

API_BASE = 'http://localhost:1999/api/v1'

def apply_token(auth_code):
    response = requests.post(
        f'{API_BASE}/apply-token',
        json={'authCode': auth_code}
    )
    return response.json()

def get_user_info(access_token):
    response = requests.post(
        f'{API_BASE}/user-info',
        json={'accessToken': access_token}
    )
    return response.json()

def process_payment(payment_data):
    response = requests.post(
        f'{API_BASE}/pay',
        json=payment_data
    )
    return response.json()
```

## Data Models

### Request Models

All request models are defined in `api/model.go`:

- `ApplyTokenRequest`
- `InquiryUserInfoRequest`
- `InquiryUserCardListRequest`
- `PayRequest`

### Response Models

Response models are defined in `superqi/model.go`:

- `ApplyTokenResponse`
- `InquiryUserInfoResponse`
- `InquiryUserCardListResponse`
- `PayResponse`

### Error Models

Error models are defined in `api/model.go`:

- `ErrorResponse`
- `PaymentStatusResponse`
- `HealthResponse`

## Environment Configuration

The API requires the following environment variables:

- `SUPERQI_GATEWAY_URL` - SuperQi gateway URL
- `SUPERQI_MERCHANT_PRIVATE_KEY_PATH` - Path to merchant private key
- `SUPERQI_CLIENT_ID` - SuperQi client identifier
- `SUPERQI_DEBUG_ENABLED` - Enable debug logging (optional, default: false)
- `SUPERQI_REQUEST_TIMEOUT` - Request timeout duration (optional, default: 25s)
- `PORT` - Server port (optional, default: 1999)

## Security Considerations

1. **HTTPS**: Always use HTTPS in production environments
2. **Access Tokens**: Store access tokens securely and implement proper token refresh logic
3. **Request Validation**: All requests are validated for required fields
4. **Error Handling**: Sensitive information is not exposed in error messages
5. **Webhook URLs**: Ensure webhook URLs are properly secured and validated

## Rate Limiting

Currently, no rate limiting is implemented. Consider implementing rate limiting in production environments.

## Monitoring and Logging

The API includes:
- Request/response logging with configurable debug output
- Health check endpoint for monitoring
- Structured error responses for better debugging

## Support

For issues or questions about the SuperQi Proxy API, please refer to:
- SuperQi official documentation
- This API documentation
- Server logs for debugging information
