package golibshared

// TokenClaim for token claim data
type TokenClaim struct {
	IsValid bool  `json:"isValid"`
	Claim   Claim `json:"claim"`
}

//Claim data token
type Claim struct {
	Aud    string `json:"aud"`
	Did    string `json:"did"`
	Member struct {
		UserID string `json:"userId"`
		ID     string `json:"id"`
		Email  string `json:"email"`
	} `json:"member"`
}
