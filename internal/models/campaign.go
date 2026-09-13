package models

import "encoding/json"

const (
	// Master response modes
	MasterResponseAuto     = "automatico"
	MasterResponseSemi     = "semi_automatico"
	MasterResponseManual   = "manual"
	DefaultMasterResponse  = MasterResponseSemi
)

// Campaign represents an RPG campaign
type Campaign struct {
	Name              string `json:"nome_campanha"`
	Description       string `json:"descricao"`
	Objective         string `json:"objetivo_principal"`
	Size              string `json:"tamanho_campanha"`
	Difficulty        string `json:"dificuldade"`
	MagicAllowed      bool   `json:"magia_permitida"`
	InitialContext    string `json:"contexto_inicial"`
	MasterResponseMode string `json:"modo_resposta_mestre"`
}

// CampaignContext represents the campaign context file
type CampaignContext struct {
	Name              string `json:"nome_campanha"`
	Description       string `json:"descricao"`
	Objective         string `json:"objetivo_principal"`
	Size              string `json:"tamanho_campanha"`
	InitialContext    string `json:"contexto_inicial"`
}

// CampaignRules represents the campaign rules
type CampaignRules struct {
	Difficulty   string `json:"dificuldade"`
	MagicAllowed bool   `json:"magia_permitida"`
}

// CampaignSummary represents the campaign summary
type CampaignSummary struct {
	History            string   `json:"historia"`
	ImportantDialogues []string `json:"dialogos_importantes"`
}

// MasterNotes represents master's notes
type MasterNotes struct {
	Annotations string `json:"anotacoes"`
	Secrets     string `json:"segredos"`
}

// NewCampaign creates a new campaign with default values
func NewCampaign(name string) *Campaign {
	return &Campaign{
		Name:               name,
		Description:        "",
		Objective:          "",
		Size:               "medio",
		Difficulty:         "Normal",
		MagicAllowed:       true,
		MasterResponseMode: DefaultMasterResponse,
	}
}

// ToJSON converts campaign to JSON
func (c *Campaign) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// FromJSON creates campaign from JSON
func (c *Campaign) FromJSON(data []byte) error {
	return json.Unmarshal(data, c)
}
