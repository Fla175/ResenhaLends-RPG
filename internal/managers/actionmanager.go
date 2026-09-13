package managers

import (
	"fmt"
	"resenhalends/internal/models"
	"resenhalends/internal/utils"
	"strings"
)

// WeaknessCalculator interface for weakness calculations
type WeaknessCalculator interface {
	GetWeaknessModifier(weaknesses []string, situation string) int
}

// ActionManager handles action resolution and game mechanics
type ActionManager struct {
	diceRoller       DiceRoller
	levelSystem      *utils.LevelSystem
	weaknessCalculator WeaknessCalculator
}

// NewActionManager creates a new ActionManager
func NewActionManager(diceRoller DiceRoller, levelSystem *utils.LevelSystem, weaknessCalculator WeaknessCalculator) *ActionManager {
	return &ActionManager{
		diceRoller:       diceRoller,
		levelSystem:      levelSystem,
		weaknessCalculator: weaknessCalculator,
	}
}

// DetermineActionType determines the type of action based on description
func (a *ActionManager) DetermineActionType(description string) models.ActionType {
	desc := strings.ToLower(description)

	// Quick actions
	for _, quickAction := range []string{"sair", "descansar", "procurar", "esconder", "andar", "mover"} {
		if strings.Contains(desc, quickAction) {
			return models.ActionTypeQuick
		}
	}

	// Combat actions
	combatKeywords := []string{"atacar", "lutar", "bater", "golpear", "defender", "esquivar", "combate", "luta"}
	for _, keyword := range combatKeywords {
		if strings.Contains(desc, keyword) {
			return models.ActionTypeCombat
		}
	}

	// Magic actions
	magicKeywords := []string{"lançar", "conjurar", "magia", "feitiço", "encantar", "curar", "abençoar"}
	for _, keyword := range magicKeywords {
		if strings.Contains(desc, keyword) {
			return models.ActionTypeMagic
		}
	}

	// Social actions
	socialKeywords := []string{"convencer", "negociar", "persuadir", "mentir", "enganar", "diplomacia", "conversar"}
	for _, keyword := range socialKeywords {
		if strings.Contains(desc, keyword) {
			return models.ActionTypeSocial
		}
	}

	// Default to exploration
	return models.ActionTypeExploration
}

// CalculateDifficulty calculates the difficulty class for an action
func (a *ActionManager) CalculateDifficulty(actionType models.ActionType, situation string) int {
	baseDifficulties := map[models.ActionType]int{
		models.ActionTypeCombat:     12,
		models.ActionTypeSocial:     10,
		models.ActionTypeExploration: 8,
		models.ActionTypeMagic:      14,
		models.ActionTypeQuick:      0, // Quick actions don't need dice
	}

	baseDifficulty := baseDifficulties[actionType]

	// Adjust based on situation
	situationModifiers := map[string]int{
		"perigoso":  +3,
		"hostil":    +4,
		"calmo":     -2,
		"favorável": -3,
		"noturno":   +2,
		"extremo":   +5,
	}

	for situation, modifier := range situationModifiers {
		if strings.Contains(strings.ToLower(situation), situation) {
			baseDifficulty += modifier
		}
	}

	// Ensure difficulty is reasonable
	if baseDifficulty < 5 {
		baseDifficulty = 5
	}
	if baseDifficulty > 20 {
		baseDifficulty = 20
	}

	return baseDifficulty
}

// ProcessAction processes a player action through the game system
func (a *ActionManager) ProcessAction(description, playerName, situation string, player *models.Player) (*models.Action, error) {
	actionType := a.DetermineActionType(description)

	var action *models.Action
	if actionType == models.ActionTypeQuick {
		action = models.NewQuickAction(description, playerName)
		action.SuccessLevel = "sucesso" // Quick actions always succeed
		return action, nil
	}

	action = models.NewAction(actionType, description, playerName)

	// Calculate modifiers
	levelBonus := a.levelSystem.CalculateLevelBonus(player.Level)
	var weaknessModifier int
	if a.weaknessCalculator != nil {
		weaknessModifier = a.weaknessCalculator.GetWeaknessModifierForWeaknesses(player.Weaknesses, situation)
	}
	action.Modifier = action.GetModifier(levelBonus, weaknessModifier)

	// Roll dice
	action.DiceRoll = a.diceRoller.Roll(20)

	// Calculate difficulty
	difficulty := a.CalculateDifficulty(actionType, situation)

	// Determine success level
	action.SuccessLevel = action.DetermineSuccessLevel(action.DiceRoll+action.Modifier, difficulty)

	return action, nil
}

// FormatActionResult formats the action result for display
func (a *ActionManager) FormatActionResult(action *models.Action) string {
	if action.Type == models.ActionTypeQuick {
		return fmt.Sprintf("%s realizou: %s (ação rápida)", action.PlayerName, action.Description)
	}

	modifierText := ""
	if action.Modifier != 0 {
		sign := "+"
		if action.Modifier < 0 {
			sign = ""
		}
		modifierText = fmt.Sprintf(" %s%d", sign, action.Modifier)
	}

	successText := map[string]string{
		"critico":       "🎲 CRÍTICO!",
		"sucesso":       "✅ Sucesso",
		"parcial":       "⚠️ Sucesso Parcial",
		"falha":         "❌ Falha",
		"falha_critica": "💀 FALHA CRÍTICA!",
	}

	return fmt.Sprintf("%s tentou: %s (d20%s: %d) %s",
		action.PlayerName, action.Description, modifierText, action.DiceRoll, successText[action.SuccessLevel])
}