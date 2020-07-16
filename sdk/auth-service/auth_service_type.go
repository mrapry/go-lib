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
	Aud string `json:"aud"`
}

// GenerateTokenRequest model
type GenerateTokenRequest struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
	Email    string `json:"email"`
	UnitCode string `json:"unit_code"`
}

// TokenRequest model
type TokenRequest struct {
	UserID   string `url:"id"`
	DeviceID string `url:"device"`
}
