package models

import "encoding/json"

// NPC represents a non-player character
type NPC struct {
	Name       string `json:"nome"`
	Sex        string `json:"sexo"`
	Race       string `json:"raca"`
	Occupation string `json:"ocupacao"`
	Motivation string `json:"motivacao"`
	Secret     string `json:"segredo"`
}

// NewNPC creates a new NPC
func NewNPC(name string) *NPC {
	return &NPC{
		Name: name,
	}
}

// ToJSON converts NPC to JSON
func (n *NPC) ToJSON() ([]byte, error) {
	return json.Marshal(n)
}

// FromJSON creates NPC from JSON
func (n *NPC) FromJSON(data []byte) error {
	return json.Unmarshal(data, n)
}
