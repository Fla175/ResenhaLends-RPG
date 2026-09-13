package models

// ActionType represents the type of action
type ActionType string

const (
	ActionTypeCombat     ActionType = "combate"
	ActionTypeSocial     ActionType = "social"
	ActionTypeExploration ActionType = "exploracao"
	ActionTypeMagic      ActionType = "magia"
	ActionTypeQuick      ActionType = "rapida"
)

// Action represents a player action in the game
type Action struct {
	Type           ActionType `json:"tipo"`
	Description    string     `json:"descricao"`
	PlayerName     string     `json:"nome_jogador"`
	DiceRoll       int        `json:"resultado_dado"`
	RequiresDice   bool       `json:"requer_dado"`
	Modifier       int        `json:"modificador"`
	SuccessLevel   string     `json:"nivel_sucesso"` // "critico", "sucesso", "parcial", "falha", "falha_critica"
	MasterResponse string     `json:"resposta_mestre"`
}

// NewAction creates a new action
func NewAction(actionType ActionType, description, playerName string) *Action {
	return &Action{
		Type:         actionType,
		Description:  description,
		PlayerName:   playerName,
		RequiresDice: true,
		Modifier:     0,
	}
}

// NewQuickAction creates a quick action that doesn't require dice
func NewQuickAction(description, playerName string) *Action {
	return &Action{
		Type:         ActionTypeQuick,
		Description:  description,
		PlayerName:   playerName,
		RequiresDice: false,
		Modifier:     0,
	}
}

// DetermineSuccessLevel determines the success level based on dice roll and difficulty
func (a *Action) DetermineSuccessLevel(diceRoll, difficulty int) string {
	finalRoll := diceRoll + a.Modifier

	if finalRoll >= 20 {
		return "critico"
	} else if finalRoll >= difficulty {
		return "sucesso"
	} else if finalRoll >= difficulty - 5 {
		return "parcial"
	} else if finalRoll > 1 {
		return "falha"
	} else {
		return "falha_critica"
	}
}

// GetModifier returns the total modifier for the action
func (a *Action) GetModifier(levelBonus, weaknessModifier int) int {
	return a.Modifier + levelBonus + weaknessModifier
}