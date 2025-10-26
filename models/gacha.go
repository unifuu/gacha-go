package models

// GachaResult represents the result of a gacha pull
type GachaResult struct {
	Characters []Character `json:"characters"`
	IsNew      []bool      `json:"isNew"` // Whether each character is new
	Timestamp  int64       `json:"timestamp"`
}

// PoolInfo represents gacha pool information
type PoolInfo struct {
	Characters []Character       `json:"characters"`
	Rates      map[string]string `json:"rates"`
	PitySystem string            `json:"pitySystem"`
}

// UserInfoResponse represents user information for API response
type UserInfoResponse struct {
	Username  string `json:"username"`
	Currency  int    `json:"currency"`
	PityCount int    `json:"pityCount"`
}

// InventoryResponse represents user inventory for API response
type InventoryResponse struct {
	Inventory []Character `json:"inventory"`
	Count     int         `json:"count"`
}

// AddCurrencyRequest represents request to add currency
type AddCurrencyRequest struct {
	Amount int `json:"amount"`
}

// CurrencyResponse represents currency update response
type CurrencyResponse struct {
	Currency int `json:"currency"`
}
