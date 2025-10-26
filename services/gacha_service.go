package services

import (
	"gacha/models"
	"math/rand"
)

// GachaService handles gacha logic
type GachaService struct {
	characterPool []models.Character
}

// NewGachaService creates a new gacha service
func NewGachaService() *GachaService {
	return &GachaService{
		characterPool: models.GetCharacterPool(),
	}
}

// PerformSinglePull performs a single gacha pull
func (s *GachaService) PerformSinglePull(user *models.User) models.Character {
	user.IncrementPity()

	// Pity system: guaranteed SSR at 90 pulls
	if user.PityCount >= 90 {
		user.ResetPity()
		return s.getRandomCharacterByRarity(5)
	}

	// Normal gacha logic
	totalRate := 0.0
	for _, char := range s.characterPool {
		totalRate += char.Rate
	}

	roll := rand.Float64() * totalRate
	currentRate := 0.0

	for _, char := range s.characterPool {
		currentRate += char.Rate
		if roll <= currentRate {
			if char.Rarity == 5 {
				user.ResetPity() // Reset pity when SSR obtained
			}
			return char
		}
	}

	return s.characterPool[len(s.characterPool)-1]
}

// PerformTenPull performs a ten-pull with guaranteed SR
func (s *GachaService) PerformTenPull(user *models.User) []models.Character {
	var characters []models.Character
	hasSR := false

	for i := 0; i < 10; i++ {
		char := s.PerformSinglePull(user)

		if i == 9 && !hasSR {
			// Last pull: force SR if no SR+ obtained
			for _, c := range characters {
				if c.Rarity >= 4 {
					hasSR = true
					break
				}
			}
			if !hasSR {
				char = s.getRandomCharacterByRarity(4)
			}
		}

		if char.Rarity >= 4 {
			hasSR = true
		}

		characters = append(characters, char)
	}

	return characters
}

// GetCharacterPool returns the character pool
func (s *GachaService) GetCharacterPool() []models.Character {
	return s.characterPool
}

// getRandomCharacterByRarity gets a random character of specific rarity
func (s *GachaService) getRandomCharacterByRarity(rarity int) models.Character {
	var filtered []models.Character
	for _, char := range s.characterPool {
		if char.Rarity == rarity {
			filtered = append(filtered, char)
		}
	}
	if len(filtered) == 0 {
		return s.characterPool[0]
	}
	return filtered[rand.Intn(len(filtered))]
}
