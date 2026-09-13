package models

import "encoding/json"

// Chapter represents a story chapter
type Chapter struct {
	Number    int      `json:"numero"`
	Title     string   `json:"titulo"`
	Summary   string   `json:"resumo"`
	Events    []string `json:"eventos"`
	Type      string   `json:"tipo,omitempty"`
	ClosedAt  string   `json:"fechado_em,omitempty"`
}

// CurrentChapter represents the currently active chapter
type CurrentChapter struct {
	Number  int      `json:"numero"`
	Events  []string `json:"eventos"`
	Open    bool     `json:"aberto"`
}

// NewChapter creates a new chapter
func NewChapter(number int, title string) *Chapter {
	return &Chapter{
		Number:  number,
		Title:   title,
		Summary: "",
		Events:  []string{},
	}
}

// NewCurrentChapter creates a new current chapter
func NewCurrentChapter(number int) *CurrentChapter {
	return &CurrentChapter{
		Number: number,
		Events: []string{},
		Open:   true,
	}
}

// AddEvent adds an event to the chapter
func (c *CurrentChapter) AddEvent(event string) {
	c.Events = append(c.Events, event)
}

// ToJSON converts chapter to JSON
func (c *Chapter) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// FromJSON creates chapter from JSON
func (c *Chapter) FromJSON(data []byte) error {
	return json.Unmarshal(data, c)
}

// ToJSON converts current chapter to JSON
func (c *CurrentChapter) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// FromJSON creates current chapter from JSON
func (c *CurrentChapter) FromJSON(data []byte) error {
	return json.Unmarshal(data, c)
}
