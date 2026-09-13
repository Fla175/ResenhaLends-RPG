package models

import "encoding/json"

// Monster represents a creature in the bestiary
type Monster struct {
	Name       string             `json:"nome"`
	Type       string             `json:"tipo"`
	Challenge  string             `json:"nd"`
	HitPoints  string             `json:"pv"`
	ArmorClass string             `json:"ca"`
	Attributes MonsterAttributes  `json:"atributos"`
}

// MonsterAttributes represents monster stats
type MonsterAttributes struct {
	Strength     int `json:"for"`
	Dexterity    int `json:"des"`
	Constitution int `json:"con"`
	Intelligence int `json:"int"`
	Wisdom       int `json:"sab"`
	Charisma     int `json:"car"`
}

// NewMonster creates a new monster
func NewMonster(name string) *Monster {
	return &Monster{
		Name: name,
		Attributes: MonsterAttributes{
			Strength:     10,
			Dexterity:    10,
			Constitution: 10,
			Intelligence: 10,
			Wisdom:       10,
			Charisma:     10,
		},
	}
}

// ToJSON converts monster to JSON
func (m *Monster) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON creates monster from JSON
func (m *Monster) FromJSON(data []byte) error {
	return json.Unmarshal(data, m)
}
