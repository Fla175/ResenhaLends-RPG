package generators

import (
	"fmt"
	"math/rand"
	"time"
)

// WeaknessGenerator handles random weakness generation for characters
type WeaknessGenerator struct {
	rand *rand.Rand
}

// NewWeaknessGenerator creates a new WeaknessGenerator
func NewWeaknessGenerator() *WeaknessGenerator {
	return &WeaknessGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

var (
	weaknesses = []string{
		// Physical weaknesses
		"Medo de fogo",
		"Alérgico a magia arcana",
		"Fragilidade em combate corpo a corpo",
		"Sensibilidade à luz solar",
		"Vulnerabilidade a venenos",
		"Fraqueza contra ataques de sombra",
		"Resistência reduzida a gelo",
		"Suscetível a ataques mentais",
		"Mobilidade reduzida em terreno difícil",
		"Dificuldade em segurar respiração",

		// Mental weaknesses
		"Fobia de altura",
		"Trauma com água profunda",
		"Incapacidade de mentir",
		"Compulsão de honestidade",
		"Memória fraca para nomes",
		"Distradil fácil",
		"Susceptível a ilusões",
		"Medo irracional de aranhas",
		"Pânico em espaços fechados",
		"Confusão em combate",

		// Social weaknesses
		"Dificuldade com autoridades",
		"Incapacidade de negociar",
		"Repulsão por mortos-vivos",
		"Desconfiança de estranhos",
		"Incapacidade de recusar desafios",
		"Fascinação por ouro",
		"Gula excessiva",
		"Ira descontrolada",
		"Susceptível a bajulação",
		"Lealdade cega",

		// Magical weaknesses
		"Incapacidade de usar magia divina",
		"Resistência nula a magia de mente",
		"Vulnerabilidade a magia de tempo",
		"Fascinação por magia antiga",
		"Consumo excessivo de mana",
		"Rejeição por espíritos",
		"Attrai criaturas mágicas",
		"Imunidade impossível de curar",
		"Instabilidade mágica",
		"Consume vida ao usar magia",

		// Combat weaknesses
		"Aim ruim com arcos",
		"Incapacidade de usar escudos",
		"Lentidão com armas pesadas",
		"Precipitação em combate",
		"Incapacidade de golpear críticos",
		"Falta de reflexos",
		"Dificuldade em defender",
		"Exaustão rápida",
		"Suscetível a surpresas",
		"Panico em números inferiores",

		// Environmental weaknesses
		"Sensibilidade a altitude",
		"Desorientação em escuridão",
		"Medo de tempestades",
		"Fragilidade em desertos",
		"Incapacidade de nadar",
		"Vulnerabilidade a plantas",
		"Medo de animais grandes",
		"Suscetível a doenças",
		"Fragilidade em extremos",
		"Consumo excessivo de recursos",
	}
)

// GenerateWeaknesses generates 1-2 random weaknesses for a character
func (w *WeaknessGenerator) GenerateWeaknesses() []string {
	count := w.rand.Intn(2) + 1 // 1-2 weaknesses
	weaknessList := make([]string, 0, count)

	// Create a copy of weaknesses to avoid duplicates
	availableWeaknesses := make([]string, len(weaknesses))
	copy(availableWeaknesses, weaknesses)

	for i := 0; i < count && len(availableWeaknesses) > 0; i++ {
		index := w.rand.Intn(len(availableWeaknesses))
		weakness := availableWeaknesses[index]
		weaknessList = append(weaknessList, weakness)

		// Remove selected weakness to avoid duplicates
		availableWeaknesses = append(availableWeaknesses[:index], availableWeaknesses[index+1:]...)
	}

	return weaknessList
}

// GenerateThematicWeakness generates a weakness based on character theme
func (w *WeaknessGenerator) GenerateThematicWeakness(theme string) string {
	themeWeaknesses := map[string][]string{
		"guerreiro": {"Fragilidade em combate corpo a corpo", "Lentidão com armas pesadas", "Exaustão rápida"},
		"mago":     {"Alérgico a magia arcana", "Consumo excessivo de mana", "Instabilidade mágica"},
		"ladino":   {"Aim ruim com arcos", "Panic em combate", "Desconfiança de estranhos"},
		"clérigo":  {"Incapacidade de usar magia divina", "Compulsão de honestidade", "Repulsão por mortos-vivos"},
		"ranger":   {"Incapacidade de nadar", "Fascinação por magia antiga", "Vulnerabilidade a plantas"},
		"bárbaro":  {"Ira descontrolada", "Precipitação em combate", "Incapacidade de negociar"},
	}

	if themed, ok := themeWeaknesses[theme]; ok {
		return themed[w.rand.Intn(len(themed))]
	}

	return w.GenerateWeaknesses()[0]
}

// FormatWeaknesses returns a formatted string of weaknesses
func (w *WeaknessGenerator) FormatWeaknesses(weaknesses []string) string {
	if len(weaknesses) == 0 {
		return "Nenhuma fraqueza"
	}

	result := "Fraquezas: "
	for i, weakness := range weaknesses {
		if i > 0 {
			result += ", "
		}
		result += weakness
	}
	return result
}

// HasWeakness checks if a character has a specific weakness
func (w *WeaknessGenerator) HasWeakness(weaknesses []string, weakness string) bool {
	for _, w := range weaknesses {
		if w == weakness {
			return true
		}
	}
	return false
}

// GetWeaknessModifier returns a modifier for dice rolls based on weakness
func (w *WeaknessGenerator) GetWeaknessModifier(weaknesses []string, situation string) int {
	// Simple logic: return -2 if weakness matches situation
	// In a full implementation, this would be more sophisticated
	situationWeaknesses := map[string]string{
		"fogo":     "Medo de fogo",
		"magia":    "Alérgico a magia arcana",
		"sombra":   "Fraqueza contra ataques de sombra",
		"altitude": "Sensibilidade a altitude",
		"água":     "Trauma com água profunda",
	}

	if expectedWeakness, ok := situationWeaknesses[situation]; ok {
		if w.HasWeakness(weaknesses, expectedWeakness) {
			return -2
		}
	}

	return 0
}

// GetWeaknessModifierForWeaknesses is the interface implementation
func (w *WeaknessGenerator) GetWeaknessModifierForWeaknesses(weaknesses []string, situation string) int {
	return w.GetWeaknessModifier(weaknesses, situation)
}