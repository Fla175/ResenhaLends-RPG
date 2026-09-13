package managers

import (
	"fmt"
	"resenhalends/internal/models"
	"strings"
)

// CodexManager handles bestiary operations
type CodexManager struct {
	jsonManager *JSONManager
}

// NewCodexManager creates a new CodexManager
func NewCodexManager(jsonManager *JSONManager) *CodexManager {
	return &CodexManager{
		jsonManager: jsonManager,
	}
}

// LoadBestiary loads the bestiary for a campaign
func (c *CodexManager) LoadBestiary(campaignName string) (map[string]*models.Monster, error) {
	c.jsonManager.MigrateBestiary(campaignName)
	
	var bestiary map[string]*models.Monster
	err := c.jsonManager.LoadInto(campaignName, ".monstros.json", &bestiary)
	if err != nil {
		return nil, err
	}
	if bestiary == nil {
		return map[string]*models.Monster{}, nil
	}
	return bestiary, nil
}

// SaveMonster saves a monster to the bestiary
func (c *CodexManager) SaveMonster(campaignName, monsterID string, monster *models.Monster) error {
	bestiary, err := c.LoadBestiary(campaignName)
	if err != nil {
		return err
	}
	
	if bestiary == nil {
		bestiary = make(map[string]*models.Monster)
	}
	
	bestiary[monsterID] = monster
	return c.jsonManager.Save(campaignName, ".monstros.json", bestiary)
}

// SearchMonster searches for monsters by name or ID
func (c *CodexManager) SearchMonster(campaignName, term string) (map[string]*models.Monster, error) {
	bestiary, err := c.LoadBestiary(campaignName)
	if err != nil {
		return nil, err
	}
	
	termLower := strings.ToLower(term)
	results := make(map[string]*models.Monster)
	
	for id, monster := range bestiary {
		nameLower := strings.ToLower(monster.Name)
		idLower := strings.ToLower(id)
		
		if strings.Contains(nameLower, termLower) || strings.Contains(idLower, termLower) {
			results[id] = monster
		}
	}
	
	return results, nil
}

// FormatMonsterSheet formats a monster's data for display
func (c *CodexManager) FormatMonsterSheet(monster *models.Monster) string {
	const separator = "\033[36m----------------------------------------\033[0m"
	
	sheet := fmt.Sprintf(`%s
\033[33m   === %s ===\033[0m
Tipo: %s | ND: %s
%s
PV: %s | CA: %s
FOR: %d | DES: %d | CON: %d
INT: %d | SAB: %d | CAR: %d
%s`,
		separator,
		strings.ToUpper(monster.Name),
		monster.Type,
		monster.Challenge,
		separator,
		monster.HitPoints,
		monster.ArmorClass,
		monster.Attributes.Strength,
		monster.Attributes.Dexterity,
		monster.Attributes.Constitution,
		monster.Attributes.Intelligence,
		monster.Attributes.Wisdom,
		monster.Attributes.Charisma,
		separator,
	)
	
	return sheet
}

// ListMonsters returns a list of all monsters with their IDs
func (c *CodexManager) ListMonsters(campaignName string) ([]string, error) {
	bestiary, err := c.LoadBestiary(campaignName)
	if err != nil {
		return nil, err
	}
	
	var ids []string
	for id := range bestiary {
		ids = append(ids, id)
	}
	
	return ids, nil
}
