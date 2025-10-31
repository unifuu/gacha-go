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
		{ID: 1, Name: "Sona", Rarity: 5, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Sona_0.jpg", Rate: 0.6},
		{ID: 2, Name: "Soraka", Rarity: 5, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Soraka_0.jpg", Rate: 0.6},
		{ID: 3, Name: "Syndra", Rarity: 5, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Syndra_0.jpg", Rate: 0.8},

		// SR (4-star) - 10% probability
		{ID: 4, Name: "Anivia", Rarity: 4, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Anivia_0.jpg", Rate: 2.5},
		{ID: 5, Name: "Annie", Rarity: 4, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Annie_0.jpg", Rate: 2.5},
		{ID: 6, Name: "Ashe", Rarity: 4, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Ashe_0.jpg", Rate: 2.5},
		{ID: 7, Name: "Azir", Rarity: 4, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Azir_0.jpg", Rate: 2.5},

		// R (3-star) - 88% probability
		{ID: 8, Name: "Bard", Rarity: 3, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Bard_0.jpg", Rate: 22},
		{ID: 9, Name: "Braum", Rarity: 3, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Braum_0.jpg", Rate: 22},
		{ID: 10, Name: "Briar", Rarity: 3, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Briar_0.jpg", Rate: 22},
		{ID: 11, Name: "Blitzcrank", Rarity: 3, ImageURL: "https://ddragon.leagueoflegends.com/cdn/img/champion/splash/Blitzcrank_0.jpg", Rate: 22},
	}
}
