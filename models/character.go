package models

// Character represents a playable character in the gacha system
type Character struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Rarity   int     `json:"rarity"`   // 3, 4, 5
	ImageURL string  `json:"imageUrl"` // Character image URL
	Rate     float64 `json:"rate"`     // Pull weight
}

// GetCharacterPool returns the global character pool
func GetCharacterPool() []Character {
	return []Character{
		// SSR (5-star) - 2% probability
		{ID: 1, Name: "유미", Rarity: 5, ImageURL: "https://example.com/char1.png", Rate: 0.6},
		{ID: 2, Name: "티모", Rarity: 5, ImageURL: "https://example.com/char2.png", Rate: 0.6},
		{ID: 3, Name: "나르", Rarity: 5, ImageURL: "https://example.com/char3.png", Rate: 0.8},

		// SR (4-star) - 10% probability
		{ID: 4, Name: "나미", Rarity: 4, ImageURL: "https://example.com/char4.png", Rate: 2.5},
		{ID: 5, Name: "럼블", Rarity: 4, ImageURL: "https://example.com/char5.png", Rate: 2.5},
		{ID: 6, Name: "레나타", Rarity: 4, ImageURL: "https://example.com/char6.png", Rate: 2.5},
		{ID: 7, Name: "룰루", Rarity: 4, ImageURL: "https://example.com/char7.png", Rate: 2.5},

		// R (3-star) - 88% probability
		{ID: 8, Name: "밀리오", Rarity: 3, ImageURL: "https://example.com/char8.png", Rate: 22},
		{ID: 9, Name: "베이가", Rarity: 3, ImageURL: "https://example.com/char9.png", Rate: 22},
		{ID: 10, Name: "뽀삐", Rarity: 3, ImageURL: "https://example.com/char10.png", Rate: 22},
		{ID: 11, Name: "소라카", Rarity: 3, ImageURL: "https://example.com/char11.png", Rate: 22},
	}
}
