package superqi

import "time"

type Result struct {
	ResultCode    string `json:"resultCode"`
	ResultStatus  string `json:"resultStatus"`
	ResultMessage string `json:"resultMessage"`
}

type ApplyTokenResponse struct {
	Result                 Result    `json:"result"`
	AccessToken            string    `json:"accessToken"`
	AccessTokenExpiryTime  time.Time `json:"accessTokenExpiryTime"`
	RefreshToken           string    `json:"refreshToken"`
	RefreshTokenExpiryTime time.Time `json:"refreshTokenExpiryTime"`
	CustomerID             string    `json:"customerId"`
}

type InquiryUserCardListResponse struct {
	Result   Result `json:"result"`
	CardList []struct {
		MaskedCardNo  string `json:"maskedCardNo"`
		AccountNumber string `json:"accountNumber"`
	} `json:"cardList"`
}

type InquiryUserInfoResponse struct {
	Result   Result `json:"result"`
	UserInfo struct {
		UserID       string `json:"userId"`
		LoginIDInfos []struct {
			LoginID     string `json:"loginId"`
			HashLoginID string `json:"hashLoginId"`
			MaskLoginID string `json:"maskLoginId"`
			LoginIDType string `json:"loginIdType"`
		} `json:"loginIdInfos"`
		UserName struct {
			FullName   string `json:"fullName"`
			FirstName  string `json:"firstName"`
			SecondName string `json:"secondName"`
			ThirdName  string `json:"thirdName"`
			LastName   string `json:"lastName"`
		} `json:"userName"`
		UserNameInArabic struct {
			FullName   string `json:"fullName"`
			FirstName  string `json:"firstName"`
			SecondName string `json:"secondName"`
			ThirdName  string `json:"thirdName"`
			LastName   string `json:"lastName"`
		} `json:"userNameInArabic"`
		Avatar       string `json:"avatar"`
		Gender       string `json:"gender"`
		BirthDate    string `json:"birthDate"`
		Nationality  string `json:"nationality"`
		ContactInfos []struct {
			ContactType string `json:"contactType"`
			ContactNo   string `json:"contactNo"`
		} `json:"contactInfos"`
	} `json:"userInfo"`
}

type PayResponse struct {
	PaymentId          string `json:"paymentId"`
	Result             Result `json:"result"`
	RedirectActionForm struct {
		Method      string `json:"method"`
		Parameters  string `json:"parameters"`
		RedirectUrl string `json:"redirectUrl"`
	} `json:"redirectActionForm"`
}
