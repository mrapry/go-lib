package auth_service

// ValidateResp model
type ValidateResponse struct {
	IsExpired      bool  `json:"isExpired"`
	IsDeviceActive bool  `json:"isDeviceActive"`
	Claim          Claim `json:"claim"`
}

// ValidateTokenRequest model
type ValidateTokenRequest struct {
	Token string `json:"token"`
}

// Claim model
type Claim struct {
	Aud    string `json:"aud"`
	Did    string `json:"did"`
	Member struct {
		ID     string `json:"id"`
		UserID string `json:"user_id"`
		Email  string `json:"email"`
	} `json:"member"`
}

// GenerateTokenRequest model
type GenerateTokenRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Device struct {
		ID       string `json:"id,omitempty"`
		Platform string `json:"platform"`
	} `json:"device"`
	IPAddress string `json:"ip_address,omitempty"`
}

// TokenRequest model
type TokenRequest struct {
	UserID   string `url:"id"`
	DeviceID string `url:"device"`
}

// ResponseGenerate model
type ResponseGenerate struct {
	Token  string                `json:"token"`
	Member *GenerateTokenRequest `json:"member"`
}
