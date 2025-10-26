package models

// User represents a player in the system
type User struct {
	ID        int         `json:"id"`
	Username  string      `json:"username"`
	Currency  int         `json:"currency"`  // Gacha currency
	Inventory []Character `json:"inventory"` // Owned characters
	PityCount int         `json:"pityCount"` // Pity counter
}

// HasCharacter checks if user owns a specific character
func (u *User) HasCharacter(charID int) bool {
	for _, char := range u.Inventory {
		if char.ID == charID {
			return true
		}
	}
	return false
}

// AddCharacter adds a character to user's inventory if not already owned
func (u *User) AddCharacter(char Character) bool {
	if u.HasCharacter(char.ID) {
		return false // Already owned
	}
	u.Inventory = append(u.Inventory, char)
	return true // Newly acquired
}

// DeductCurrency deducts currency from user
func (u *User) DeductCurrency(amount int) bool {
	if u.Currency < amount {
		return false
	}
	u.Currency -= amount
	return true
}

// AddCurrency adds currency to user
func (u *User) AddCurrency(amount int) {
	u.Currency += amount
}

// IncrementPity increments the pity counter
func (u *User) IncrementPity() {
	u.PityCount++
}

// ResetPity resets the pity counter
func (u *User) ResetPity() {
	u.PityCount = 0
}
