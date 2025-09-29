package superqi

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"superqi-proxy/util"
	"time"
)

var Interface *Client

type Config struct {
	GatewayURL             string
	MerchantPrivateKeyPath string
	ClientID               string
	IsDebug                bool
}

type Client struct {
	config     Config
	privateKey *rsa.PrivateKey

	httpClient *http.Client
}

func loadEnvConfig() (Config, error) {
	gatewayURL := os.Getenv("SUPERQI_GATEWAY_URL")
	if gatewayURL == "" {
		return Config{}, errors.New("SUPERQI_GATEWAY_URL is not set")
	}

	merchantPrivateKeyPath := os.Getenv("SUPERQI_MERCHANT_PRIVATE_KEY_PATH")
	if merchantPrivateKeyPath == "" {
		return Config{}, errors.New("SUPERQI_MERCHANT_PRIVATE_KEY_PATH is not set")
	}

	clientID := os.Getenv("SUPERQI_CLIENT_ID")
	if clientID == "" {
		return Config{}, errors.New("SUPERQI_CLIENT_ID is not set")
	}

	isDebug := util.GetEnvOrDefaultBool("SUPERQI_DEBUG_ENABLED", false)

	return Config{
		GatewayURL:             gatewayURL,
		MerchantPrivateKeyPath: merchantPrivateKeyPath,
		ClientID:               clientID,
		IsDebug:                isDebug,
	}, nil
}

func InitSuperQiClient() error {
	config, err := loadEnvConfig()
	if err != nil {
		return err
	}

	privateKey, err := loadPrivateKey(config.MerchantPrivateKeyPath)
	if err != nil {
		return err
	}

	Interface = &Client{
		config:     config,
		privateKey: privateKey,
		httpClient: &http.Client{
			Timeout: util.GetEnvOrDefaultDuration("SUPERQI_REQUEST_TIMEOUT", time.Second*25),
		},
	}
	return nil
}

func (client *Client) debugPrintln(v ...any) {
	if client.config.IsDebug {
		log.Println(v...)
	}
}

func formatDuration(d time.Duration) string {
	ms := d.Milliseconds()
	if ms < 1000 {
		return fmt.Sprintf("%d ms", ms)
	}
	return fmt.Sprintf("%s ms", fmt.Sprintf("%0.0f", float64(ms)))
}

func (client *Client) buildHeaders(method, path string, params map[string]any) (map[string]string, error) {
	currentTimestamp := time.Now().Format("2006-01-02T15:04:05-07:00")
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	signature, err := client.generateSignature(method, path, currentTimestamp, string(paramsJSON))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		"Client-Id":    client.config.ClientID,
		"Request-Time": currentTimestamp,
		"Signature":    fmt.Sprintf("algorithm=RSA256, keyVersion=1, signature=%s", signature),
	}, nil
}

func (client *Client) generateSignature(httpMethod, path, requestTime, content string) (string, error) {
	signContent := fmt.Sprintf("%s %s\n%s.%s.%s", httpMethod, path, client.config.ClientID, requestTime, content)
	hash := sha256.Sum256([]byte(signContent))

	signature, err := rsa.SignPKCS1v15(nil, client.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func (client *Client) sendRequest(path, method string, headers map[string]string, params map[string]any) ([]byte, error) {
	startTime := time.Now()
	debugBuilder := strings.Builder{}

	defer func(builder *strings.Builder) {
		client.debugPrintln(builder.String())
	}(&debugBuilder)

	// Format request info section
	debugBuilder.WriteString("\n================================================\n")
	debugBuilder.WriteString("================= Request Info =================\n")
	debugBuilder.WriteString("================================================\n")
	debugBuilder.WriteString(fmt.Sprintf("[%s] %s\n\n", method, path))
	debugBuilder.WriteString(fmt.Sprintf("Full URL        : %s%s\n", client.config.GatewayURL, path))
	debugBuilder.WriteString(fmt.Sprintf("Request Date    : %s\n", startTime.Format("2006-01-02 15:04:05.000")))

	// Format request headers section
	debugBuilder.WriteString("================================\n")
	debugBuilder.WriteString("======= Request Headers ========\n")
	debugBuilder.WriteString("================================\n")

	// Find the longest header name for alignment
	maxHeaderLen := 0
	for key := range headers {
		if len(key) > maxHeaderLen {
			maxHeaderLen = len(key)
		}
	}

	for key, value := range headers {
		// Handle sensitive headers
		displayValue := value
		if strings.Contains(strings.ToLower(key), "signature") ||
			strings.Contains(strings.ToLower(key), "authorization") ||
			strings.Contains(strings.ToLower(key), "client-id") {
			displayValue = "REDACTED"
		}
		debugBuilder.WriteString(fmt.Sprintf("%-*s : %s\n", maxHeaderLen, key, displayValue))
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// Format request body section
	debugBuilder.WriteString("================================\n")
	debugBuilder.WriteString("========= Request Body =========\n")
	debugBuilder.WriteString("================================\n")

	// Pretty print JSON
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, paramsJSON, "", "  ")
	if err != nil {
		debugBuilder.WriteString(string(paramsJSON))
	} else {
		debugBuilder.WriteString(prettyJSON.String())
	}
	debugBuilder.WriteString("\n")

	req, err := http.NewRequest(method, client.config.GatewayURL+path, strings.NewReader(string(paramsJSON)))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.httpClient.Do(req)
	responseTime := time.Now()
	executionTime := responseTime.Sub(startTime)

	if err != nil {
		debugBuilder.WriteString("================================================\n")
		debugBuilder.WriteString("================= Request Error ================\n")
		debugBuilder.WriteString("================================================\n")
		debugBuilder.WriteString(fmt.Sprintf("Error           : %v\n", err))
		debugBuilder.WriteString(fmt.Sprintf("Response Time   : %s\n", formatDuration(executionTime)))
		debugBuilder.WriteString("================================================\n")
		return nil, err
	}
	defer resp.Body.Close()

	// Format response info section
	debugBuilder.WriteString("================================================\n")
	debugBuilder.WriteString("================ Response Info =================\n")
	debugBuilder.WriteString("================================================\n")
	debugBuilder.WriteString(fmt.Sprintf("HTTP Status      : %s\n", resp.Status))
	debugBuilder.WriteString(fmt.Sprintf("Response Date    : %s\n", responseTime.Format("2006-01-02 15:04:05.000")))
	debugBuilder.WriteString(fmt.Sprintf("Response Time    : %s\n", formatDuration(executionTime)))

	// Format response headers section
	debugBuilder.WriteString("================================\n")
	debugBuilder.WriteString("======= Response Headers =======\n")
	debugBuilder.WriteString("================================\n")

	// Find the longest response header name for alignment
	maxRespHeaderLen := 0
	for key := range resp.Header {
		if len(key) > maxRespHeaderLen {
			maxRespHeaderLen = len(key)
		}
	}

	for key, values := range resp.Header {
		for _, value := range values {
			debugBuilder.WriteString(fmt.Sprintf("%-*s : %s\n", maxRespHeaderLen, key, value))
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		debugBuilder.WriteString("================================\n")
		debugBuilder.WriteString("======= Response Read Error ====\n")
		debugBuilder.WriteString("================================\n")
		debugBuilder.WriteString(fmt.Sprintf("Error           : %v\n", err))
		debugBuilder.WriteString("================================================\n")
		return nil, err
	}

	// Format response body section
	debugBuilder.WriteString("================================\n")
	debugBuilder.WriteString("======== Response Body =========\n")
	debugBuilder.WriteString("================================\n")

	// Pretty print JSON response
	var prettyResponseJSON bytes.Buffer
	err = json.Indent(&prettyResponseJSON, body, "", "  ")
	if err != nil {
		debugBuilder.WriteString(string(body))
	} else {
		debugBuilder.WriteString(prettyResponseJSON.String())
	}
	debugBuilder.WriteString("\n")
	debugBuilder.WriteString("================================================\n")

	return body, nil
}
