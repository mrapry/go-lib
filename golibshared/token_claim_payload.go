package golibshared

// TokenClaim for token claim data
type TokenClaim struct {
	DeviceID string `json:"did"`
	User     struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
	Alg string `json:"-"`
}
