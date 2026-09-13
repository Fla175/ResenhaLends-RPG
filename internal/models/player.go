package models

import "encoding/json"

// Macro represents a player macro with name, description, and command
type Macro struct {
	Name        string `json:"nome"`
	Description string `json:"descricao"`
	Command     string `json:"comando"`
	Key         string `json:"atalho"` // Numeric key (1-9)
}

// Player represents a RPG player character
type Player struct {
	Name         string   `json:"nome"`
	Class        string   `json:"classe"`
	Description  string   `json:"descricao"`
	Race         string   `json:"raca"`
	Life         int      `json:"vida"`
	Inventory    []string `json:"inventario"`
	Macros       []Macro  `json:"macros"`
	Attributes   Attributes `json:"atributos"`
	QuickActions []string `json:"acoes_rapidas"`
	Level        int      `json:"nivel"`
	XP           int      `json:"xp"`
	XPToNextLevel int     `json:"xp_proximo_nivel"`
	Weaknesses   []string `json:"fraquezas"`
}

// Attributes represents player stats
type Attributes struct {
	Forca       int `json:"for"`
	Habilidade  int `json:"hab"`
	Resistencia int `json:"res"`
	Armadura    int `json:"arm"`
	PoderFogo   int `json:"pdf"`
}

// NewPlayer creates a new player with default values
func NewPlayer(name string) *Player {
	return &Player{
		Name:         name,
		Class:        "",
		Description:  "",
		Race:         "Humano",
		Life:         20,
		Inventory:    []string{},
		Macros:       []Macro{},
		QuickActions: []string{"Sair", "Descansar", "Procurar", "Esconder"},
		Level:        1,
		XP:           0,
		XPToNextLevel: 100,
		Weaknesses:   []string{},
		Attributes: Attributes{
			Forca:       10,
			Habilidade:  10,
			Resistencia: 10,
			Armadura:    10,
			PoderFogo:   10,
		},
	}
}

// Merge updates the player with data from a patch
func (p *Player) Merge(patch map[string]interface{}) {
	if name, ok := patch["nome"].(string); ok && name != "" {
		p.Name = name
	}
	if class, ok := patch["classe"].(string); ok {
		p.Class = class
	}
	if desc, ok := patch["descricao"].(string); ok {
		p.Description = desc
	}
	if race, ok := patch["raca"].(string); ok {
		p.Race = race
	}
	if life, ok := patch["vida"].(float64); ok {
		p.Life = int(life)
		if p.Life < 0 {
			p.Life = 0
		}
		if p.Life > 99 {
			p.Life = 99
		}
	}
	if inv, ok := patch["inventario"]; ok {
		switch v := inv.(type) {
		case string:
			p.Inventory = parseInventory(v)
		case []interface{}:
			p.Inventory = parseInventoryInterface(v)
		case []string:
			p.Inventory = v
		}
	}
	if macros, ok := patch["macros"]; ok {
		switch v := macros.(type) {
		case []interface{}:
			macros := make([]Macro, len(v))
			for i, item := range v {
				if macroMap, ok := item.(map[string]interface{}); ok {
					macros[i] = Macro{
						Name:        getString(macroMap, "nome"),
						Description: getString(macroMap, "descricao"),
						Command:     getString(macroMap, "comando"),
						Key:         getString(macroMap, "atalho"),
					}
				}
			}
			p.Macros = macros
		case []Macro:
			p.Macros = v
		}
	}
	if atts, ok := patch["atributos"].(map[string]interface{}); ok {
		if p.Attributes.Forca == 0 {
			p.Attributes = Attributes{Forca: 10, Habilidade: 10, Resistencia: 10, Armadura: 10, PoderFe: 10}
		}
		if val, ok := atts["for"].(float64); ok {
			p.Attributes.Forca = int(val)
		}
		if val, ok := atts["hab"].(float64); ok {
			p.Attributes.Habilidade = int(val)
		}
		if val, ok := atts["res"].(float64); ok {
			p.Attributes.Resistencia = int(val)
		}
		if val, ok := atts["arm"].(float64); ok {
			p.Attributes.Armadura = int(val)
		}
		if val, ok := atts["pdf"].(float64); ok {
			p.Attributes.PoderFe = int(val)
		}
	}
	if level, ok := patch["nivel"].(float64); ok {
		p.Level = int(level)
	}
	if xp, ok := patch["xp"].(float64); ok {
		p.XP = int(xp)
	}
	if xpNext, ok := patch["xp_proximo_nivel"].(float64); ok {
		p.XPToNextLevel = int(xpNext)
	}
	if weaknesses, ok := patch["fraquezas"]; ok {
		switch v := weaknesses.(type) {
		case []interface{}:
			weaknessList := make([]string, len(v))
			for i, item := range v {
				if str, ok := item.(string); ok {
					weaknessList[i] = str
				}
			}
			p.Weaknesses = weaknessList
		case []string:
			p.Weaknesses = v
		}
	}
	if quickActions, ok := patch["acoes_rapidas"]; ok {
		switch v := quickActions.(type) {
		case []interface{}:
			actions := make([]string, len(v))
			for i, item := range v {
				if str, ok := item.(string); ok {
					actions[i] = str
				}
			}
			p.QuickActions = actions
		case []string:
			p.QuickActions = v
		}
	}
}

func parseInventory(s string) []string {
	var result []string
	for _, item := range splitString(s, ",") {
		if trimmed := trimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func parseInventoryInterface(items []interface{}) []string {
	result := make([]string, len(items))
	for i, item := range items {
		if str, ok := item.(string); ok {
			result[i] = str
		}
	}
	return result
}

// Helper functions for string manipulation
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// ToJSON converts player to JSON
func (p *Player) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON creates player from JSON
func (p *Player) FromJSON(data []byte) error {
	return json.Unmarshal(data, p)
}
