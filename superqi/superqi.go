package superqi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	ONLINE_PURCHASE              = "51051000101000000011"
	AGREEMENT_PAYMENT            = "51051000101000100031"
	ONLINE_PURCHASE_AUTH_CAPTURE = "51051000101000000012"
)

func (client *Client) ApplyToken(authCode string) (ApplyTokenResponse, error) {
	const path = "/v1/authorizations/applyToken"
	params := map[string]any{
		"grantType": "AUTHORIZATION_CODE",
		"authCode":  authCode,
	}

	headers, err := client.buildHeaders("POST", path, params)
	if err != nil {
		return ApplyTokenResponse{}, err
	}

	response, err := client.sendRequest(path, "POST", headers, params)
	if err != nil {
		return ApplyTokenResponse{}, err
	}

	var body ApplyTokenResponse
	err = json.Unmarshal(response, &body)
	return body, err
}

func (client *Client) InquiryUserInfo(accessToken string) (InquiryUserInfoResponse, error) {
	const path = "/v1/users/inquiryUserInfo"
	params := map[string]any{
		"accessToken": accessToken,
	}

	headers, err := client.buildHeaders("POST", "/v1/users/inquiryUserInfo", params)
	if err != nil {
		return InquiryUserInfoResponse{}, err
	}

	response, err := client.sendRequest(path, "POST", headers, params)
	if err != nil {
		return InquiryUserInfoResponse{}, err
	}

	var body InquiryUserInfoResponse
	err = json.Unmarshal(response, &body)
	return body, err
}

func (client *Client) InquiryUserCardList(accessToken string) (InquiryUserCardListResponse, error) {
	const path = "/v1/users/inquiryUserCardList"
	params := map[string]any{
		"accessToken": accessToken,
	}

	headers, err := client.buildHeaders("POST", path, params)
	if err != nil {
		return InquiryUserCardListResponse{}, err
	}

	response, err := client.sendRequest(path, "POST", headers, params)
	if err != nil {
		return InquiryUserCardListResponse{}, err
	}

	var body InquiryUserCardListResponse
	err = json.Unmarshal(response, &body)
	return body, err
}

func (client *Client) Pay(amount int, requestId, accessToken, customerId, orderDesc, notifyUrl string) (PayResponse, error) {
	const path = "/v1/payments/pay"
	params := map[string]any{
		"paymentAuthCode": accessToken,
		"paymentAmount": map[string]any{
			"currency": "IQD",
			"value":    strconv.Itoa(amount * 1000),
		},
		"productCode":       ONLINE_PURCHASE,
		"paymentRequestId":  requestId,
		"paymentOrderTitle": orderDesc,
		"order": map[string]any{
			"orderDescription": orderDesc,
			"buyer": map[string]any{
				"referenceBuyerId": customerId,
			},
		},
		"paymentNotifyUrl": notifyUrl,
	}

	headers, err := client.buildHeaders("POST", path, params)
	if err != nil {
		return PayResponse{}, fmt.Errorf("failed to build headers %w", err)
	}

	response, err := client.sendRequest(path, "POST", headers, params)
	if err != nil {
		return PayResponse{}, fmt.Errorf("failed to send request %w", err)
	}

	var body PayResponse
	err = json.Unmarshal(response, &body)

	if err != nil {
		return PayResponse{}, fmt.Errorf("failed to decode response body %w", err)
	}

	if body.Result.ResultStatus == "F" {
		return PayResponse{}, fmt.Errorf("failed to make payment: %v", body.Result.ResultMessage)
	}
	return body, nil
}
