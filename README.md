# SuperQI Proxy

A high-performance HTTP proxy server for SuperQi, built with Go and Fiber. This proxy provides a RESTful API interface to interact with SuperQI's payment services, including user authentication, payment processing, and user information management.

## 🚀 Features

- **Payment Processing**: Complete payment flow with SuperQI gateway
- **User Authentication**: OAuth2-like token management (authorization codes → access tokens)
- **User Management**: Retrieve user information and registered payment cards
- **RESTful API**: Clean, well-documented HTTP endpoints
- **Web Interface For Testing**: Built-in frontend for testing and management
- **Cross-Platform**: Builds for 17+ platform combinations
- **High Performance**: Built with Go Fiber for optimal performance
- **Health Monitoring**: Built-in health check endpoints
- **Structured Logging**: Comprehensive request/response logging
- **Security**: Localhost-only access for sensitive endpoints

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [API Documentation](#-api-documentation)
- [Building](#-building)
- [Development](#-development)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

## 🏃 Quick Start

### Prerequisites

- Go 1.25.0 or higher
- SuperQI merchant account and credentials
- Private key file for SuperQI authentication

### 1. Clone and Setup

```bash
git clone https://github.com/MouamleH/superqi-proxy
cd superqi-proxy

# Create environment file
cp .env.example .env
# Edit .env with your SuperQI credentials
```

### 2. Configure Environment

Create a `.env` file with the following variables:

```env
# SuperQI Configuration
SUPERQI_GATEWAY_URL=https://gateway-qiuat.banqinonprod.com
SUPERQI_MERCHANT_PRIVATE_KEY_PATH=/path/to/your/private-key.pem
SUPERQI_CLIENT_ID=your_client_id

# Optional Configuration
SUPERQI_DEBUG_ENABLED=true
SUPERQI_REQUEST_TIMEOUT=25s
PORT=1999
```

### 3. Run the Server

```bash
# Install dependencies
go mod download

# Run the server
make run
# or
go run .
```

The server will start on `http://localhost:1999`

### 4. Test the API

```bash
# Health check
curl http://localhost:1999/api/v1/health

# Web interface (localhost only)
open http://localhost:1999/_
```

## 📦 Installation

### Option 1: Build from Source

```bash
# Clone repository
git clone <repository-url>
cd superqi-proxy

# Build for current platform
make build

# Build for all platforms
make build-all
```

### Option 2: Download Pre-built Binaries

Download the appropriate binary for your platform from the releases page and run:

```bash
chmod +x superqi-proxy-*
./superqi-proxy-linux-amd64  # Example for Linux
```

## ⚙️ Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SUPERQI_GATEWAY_URL` | ✅ | - | SuperQI gateway base URL |
| `SUPERQI_MERCHANT_PRIVATE_KEY_PATH` | ✅ | - | Path to merchant private key file |
| `SUPERQI_CLIENT_ID` | ✅ | - | SuperQI client identifier |
| `SUPERQI_DEBUG_ENABLED` | ❌ | `false` | Enable debug logging |
| `SUPERQI_REQUEST_TIMEOUT` | ❌ | `25s` | Request timeout duration |
| `PORT` | ❌ | `1999` | Server port |

### Private Key Setup

1. Obtain your private key from SuperQI merchant portal
2. Save it as a `.pem` file
3. Set the full path in `SUPERQI_MERCHANT_PRIVATE_KEY_PATH`
4. Ensure the file has proper permissions (`chmod 600`)

## 📖 API Documentation

### Base URL
```
http://localhost:1999/api/v1
```

### Endpoints Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/apply-token` | Exchange auth code for access token |
| `POST` | `/user-info` | Get user information |
| `POST` | `/user-cards` | Get user's registered cards |
| `POST` | `/pay` | Process payment |
| `GET` | `/payment/{id}/status` | Get payment status |
| `GET` | `/health` | Health check |

### Example: Process Payment

```bash
curl -X POST http://localhost:1999/api/v1/pay \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000,
    "requestId": "unique-request-123",
    "accessToken": "your_access_token",
    "customerId": "customer_123",
    "orderDesc": "Payment for Order #12345",
    "notifyUrl": "https://your-domain.com/webhook"
  }'
```

### Response Format

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

For complete API documentation, see [docs/API.md](docs/API.md).

## 🛠️ Development

### Setup Development Environment

```bash
# Clone and setup
git clone <repository-url>
cd superqi-proxy
go mod download
```

### Development Commands

```bash
make run          # Run in development mode
make test         # Run tests
make dev          # Build with race detection
make clean        # Clean build artifacts
make help         # Show all available commands
```

### Project Structure

```
superqi-proxy/
├── api/              # HTTP API handlers and models
│   ├── handlers.go   # Request handlers
│   ├── model.go      # Request/response models
│   └── routes.go     # Route definitions
├── superqi/          # SuperQI client implementation
│   ├── client.go     # HTTP client
│   ├── model.go      # SuperQI API models
│   ├── superqi.go    # Core SuperQI functions
│   └── util.go       # Utilities
├── util/             # General utilities
├── res/frontend/     # Embedded web interface
├── docs/             # Documentation
├── main.go           # Application entry point
└── Makefile          # Build automation
```

### Adding New Endpoints

1. Define request/response models in `api/model.go`
2. Implement handler in `api/handlers.go`
3. Add route in `api/routes.go`
4. Add SuperQI client method in `superqi/superqi.go` if needed
5. Update API documentation

### Testing

> Men test in PROD

## 🚀 Deployment

### Production Deployment

1. **Build for target platform**:
   ```bash
   GOOS=linux GOARCH=amd64 make build
   ```

2. **Prepare environment**:
   ```bash
   # Copy binary and config
   scp superqi-proxy user@server:/opt/superqi-proxy/
   scp .env user@server:/opt/superqi-proxy/
   ```

3. **Setup systemd service** (Linux):
   ```ini
   [Unit]
   Description=SuperQI Proxy
   After=network.target

   [Service]
   Type=simple
   User=superqi
   WorkingDirectory=/opt/superqi-proxy
   ExecStart=/opt/superqi-proxy/superqi-proxy
   Restart=always

   [Install]
   WantedBy=multi-user.target
   ```

### Docker Deployment

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags "-s -w" -o superqi-proxy .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/superqi-proxy .
EXPOSE 1999
CMD ["./superqi-proxy"]
```

### Environment-Specific Configuration

- **Development**: Enable debug logging, use test credentials
- **Staging**: Mirror production setup, use staging SuperQI environment
- **Production**: Disable debug logging, use HTTPS, implement monitoring

### Monitoring

- Health check endpoint: `GET /api/v1/health`
- Structured logs for monitoring systems
- Consider implementing metrics collection (Prometheus, etc.)

## 🔒 Security Considerations

- **HTTPS**: Always use HTTPS in production
- **Access Control**: Web interface restricted to localhost
- **Token Security**: Implement proper token storage and refresh logic
- **Input Validation**: All inputs validated and sanitized
- **Error Handling**: No sensitive data exposed in error messages
- **Private Keys**: Secure private key storage and permissions

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go conventions and formatting (`gofmt`, `golint`)
- Add tests for new functionality
- Update documentation for API changes
- Ensure all builds pass (`make build-all`)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Check [docs/API.md](docs/API.md) for detailed API docs
- **Build Issues**: See [BUILD.md](BUILD.md) for build troubleshooting
- **SuperQI Integration**: Refer to SuperQI official documentation
- **Issues**: Open an issue on GitHub for bugs or feature requests

## 📝 Changelog

### Latest Changes
- ✅ Complete SuperQI payment integration
- ✅ RESTful API with comprehensive validation
- ✅ Multi-platform build system
- ✅ Embedded web interface
- ✅ Health monitoring and logging
- ✅ OAuth2-like authentication flow

---

**Made with ❤️ by AI for seamless SuperQi integration**
