package managers

import (
	"fmt"
	"resenhalends/internal/models"
	"resenhalends/internal/utils"
	"strings"
)

// PlayerManager handles player operations
type PlayerManager struct {
	jsonManager *JSONManager
	diceRoller  DiceRoller
	levelSystem *utils.LevelSystem
}

// DiceRoller interface for dice rolling
type DiceRoller interface {
	Roll(faces int) int
}

// NewPlayerManager creates a new PlayerManager
func NewPlayerManager(jsonManager *JSONManager, diceRoller DiceRoller) *PlayerManager {
	return &PlayerManager{
		jsonManager: jsonManager,
		diceRoller:  diceRoller,
		levelSystem: utils.NewLevelSystem(),
	}
}

// List returns all players in a campaign
func (p *PlayerManager) List(campaignName string) ([]*models.Player, error) {
	var players []*models.Player
	err := p.jsonManager.LoadInto(campaignName, ".jogadores.json", &players)
	if err != nil {
		return nil, err
	}
	if players == nil {
		return []*models.Player{}, nil
	}
	return players, nil
}

// Add adds a new player to the campaign
func (p *PlayerManager) Add(campaignName string, player *models.Player, rollAttributes bool) (*models.Player, error) {
	players, err := p.List(campaignName)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(player.Name)
	if name == "" {
		return nil, fmt.Errorf("player name is required")
	}

	// Check for duplicate names
	for _, existingPlayer := range players {
		if strings.ToLower(existingPlayer.Name) == strings.ToLower(name) {
			return nil, fmt.Errorf("player '%s' already exists", name)
		}
	}

	// Create new player with defaults
	newPlayer := models.NewPlayer(name)
	newPlayer.Class = player.Class
	newPlayer.Description = player.Description
	newPlayer.Race = player.Race
	if player.Race == "" {
		newPlayer.Race = "Humano"
	}
	newPlayer.Life = player.Life
	if newPlayer.Life == 0 {
		newPlayer.Life = 20
	}
	newPlayer.Inventory = player.Inventory
	if newPlayer.Inventory == nil {
		newPlayer.Inventory = []string{}
	}

	// Set initial XP values
	newPlayer.XP = 0
	newPlayer.Level = 1
	newPlayer.XPToNextLevel = p.levelSystem.XPToNextLevel(1)

	// Roll attributes if requested
	if rollAttributes && p.diceRoller != nil {
		attributes := []string{"for", "hab", "res", "arm", "pdf"}
		atts := make(map[string]int)
		for _, attr := range attributes {
			atts[attr] = p.diceRoller.Roll(20)
		}
		newPlayer.Attributes.Forca = atts["for"]
		newPlayer.Attributes.Habilidade = atts["hab"]
		newPlayer.Attributes.Resistencia = atts["res"]
		newPlayer.Attributes.Armadura = atts["arm"]
		newPlayer.Attributes.PoderFogo = atts["pdf"]
	}

	players = append(players, newPlayer)
	if err := p.jsonManager.Save(campaignName, ".jogadores.json", players); err != nil {
		return nil, err
	}

	return newPlayer, nil
}

// Update updates an existing player
func (p *PlayerManager) Update(campaignName, playerName string, patch map[string]interface{}) (*models.Player, error) {
	players, err := p.List(campaignName)
	if err != nil {
		return nil, err
	}

	var found bool
	var updatedPlayer *models.Player

	for i, player := range players {
		if player.Name == playerName {
			player.Merge(patch)
			players[i] = player
			updatedPlayer = player
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("player '%s' not found", playerName)
	}

	if err := p.jsonManager.Save(campaignName, ".jogadores.json", players); err != nil {
		return nil, err
	}

	return updatedPlayer, nil
}

// Remove removes a player from the campaign
func (p *PlayerManager) Remove(campaignName, playerName string) ([]*models.Player, error) {
	players, err := p.List(campaignName)
	if err != nil {
		return nil, err
	}

	var filteredPlayers []*models.Player
	for _, player := range players {
		if player.Name != playerName {
			filteredPlayers = append(filteredPlayers, player)
		}
	}

	if err := p.jsonManager.Save(campaignName, ".jogadores.json", filteredPlayers); err != nil {
		return nil, err
	}

	return filteredPlayers, nil
}

// AdjustLife adjusts a player's life points
func (p *PlayerManager) AdjustLife(campaignName, playerName string, delta int) (int, error) {
	players, err := p.List(campaignName)
	if err != nil {
		return 0, err
	}

	var found bool
	var newLife int

	for i, player := range players {
		if player.Name == playerName {
			player.Life += delta
			if player.Life < 0 {
				player.Life = 0
			}
			if player.Life > 99 {
				player.Life = 99
			}
			players[i] = player
			newLife = player.Life
			found = true
			break
		}
	}

	if !found {
		return 0, fmt.Errorf("player '%s' not found", playerName)
	}

	if err := p.jsonManager.Save(campaignName, ".jogadores.json", players); err != nil {
		return 0, err
	}

	return newLife, nil
}

// AddMacro adds a new macro to a player
func (p *PlayerManager) AddMacro(campaignName, playerName string, macro models.Macro) error {
	players, err := p.List(campaignName)
	if err != nil {
		return err
	}

	var found bool
	for i, player := range players {
		if player.Name == playerName {
			// Check if key already exists
			for _, existingMacro := range player.Macros {
				if existingMacro.Key == macro.Key {
					return fmt.Errorf("macro key '%s' already exists", macro.Key)
				}
			}
			player.Macros = append(player.Macros, macro)
			players[i] = player
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("player '%s' not found", playerName)
	}

	return p.jsonManager.Save(campaignName, ".jogadores.json", players)
}

// RemoveMacro removes a macro from a player by key
func (p *PlayerManager) RemoveMacro(campaignName, playerName, key string) error {
	players, err := p.List(campaignName)
	if err != nil {
		return err
	}

	var found bool
	for i, player := range players {
		if player.Name == playerName {
			var newMacros []models.Macro
			for _, macro := range player.Macros {
				if macro.Key != key {
					newMacros = append(newMacros, macro)
				}
			}
			player.Macros = newMacros
			players[i] = player
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("player '%s' not found", playerName)
	}

	return p.jsonManager.Save(campaignName, ".jogadores.json", players)
}

// GetMacroBy returns a macro by its key for a player
func (p *PlayerManager) GetMacroBy(campaignName, playerName, key string) (*models.Macro, error) {
	players, err := p.List(campaignName)
	if err != nil {
		return nil, err
	}

	for _, player := range players {
		if player.Name == playerName {
			for _, macro := range player.Macros {
				if macro.Key == key {
					return &macro, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("macro with key '%s' not found for player '%s'", key, playerName)
}

// AddXP adds experience points to a player and handles level ups
func (p *PlayerManager) AddXP(campaignName, playerName string, xpGained int) (int, int, bool, error) {
	players, err := p.List(campaignName)
	if err != nil {
		return 0, 0, false, err
	}

	var found bool
	var newXP, newLevel int
	var leveledUp bool

	for i, player := range players {
		if player.Name == playerName {
			newXP, newLevel, leveledUp = p.levelSystem.AddXP(player.XP, player.Level, xpGained)
			player.XP = newXP
			player.Level = newLevel
			player.XPToNextLevel = p.levelSystem.XPToNextLevel(newLevel)

			// Increase attributes on level up
			if leveledUp {
				player.Attributes.Forca += 1
				player.Attributes.Habilidade += 1
				player.Attributes.Resistencia += 1
				player.Attributes.PoderFogo += 1
				player.Life += 5 // Increase max HP
			}

			players[i] = player
			found = true
			break
		}
	}

	if !found {
		return 0, 0, false, fmt.Errorf("player '%s' not found", playerName)
	}

	if err := p.jsonManager.Save(campaignName, ".jogadores.json", players); err != nil {
		return 0, 0, false, err
	}

	return newXP, newLevel, leveledUp, nil
}

// GetLevelProgress returns the formatted level progress for a player
func (p *PlayerManager) GetLevelProgress(player *models.Player) string {
	return p.levelSystem.FormatXPProgress(player.XP, player.Level)
}

// GetLevelBonus returns the level bonus modifier for dice rolls
func (p *PlayerManager) GetLevelBonus(player *models.Player) int {
	return p.levelSystem.CalculateLevelBonus(player.Level)
}

// FormatWeaknesses returns a formatted string of character weaknesses
func (p *PlayerManager) FormatWeaknesses(player *models.Player) string {
	if len(player.Weaknesses) == 0 {
		return "Nenhuma fraqueza"
	}

	result := "Fraquezas: "
	for i, weakness := range player.Weaknesses {
		if i > 0 {
			result += ", "
		}
		result += weakness
	}
	return result
}

// FormatAttributes returns a formatted string of player attributes with emojis
func (p *PlayerManager) FormatAttributes(player *models.Player) string {
	atts := player.Attributes
	return fmt.Sprintf(
		"💪 Força: %d | ⚡ Habilidade: %d | 🛡️ Resistência: %d | 🎽 Armadura: %d | 🔥 Poder de Fogo: %d",
		atts.Forca, atts.Habilidade, atts.Resistencia, atts.Armadura, atts.PoderFogo,
	)
}
