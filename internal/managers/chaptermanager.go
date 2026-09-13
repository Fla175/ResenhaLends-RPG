package managers

import (
	"fmt"
	"resenhalends/internal/models"
	"strings"
	"time"
)

// ChapterManager handles chapter operations
type ChapterManager struct {
	jsonManager *JSONManager
	aiClient    AIClient
}

// AIClient interface for AI operations
type AIClient interface {
	GenerateInitialContext(description, objective, campaignName string) (string, error)
	CloseChapter(number int, events []string, objective string) (*models.Chapter, error)
	ProcessAction(campaignName, action string, diceRoll int) (string, error)
	Narrate(campaignName, outcome string) (string, error)
	CheckConnection() (bool, string)
	MasterResponse(campaignName, action, outcome, mode string) (string, error)
}

// NewChapterManager creates a new ChapterManager
func NewChapterManager(jsonManager *JSONManager, aiClient AIClient) *ChapterManager {
	return &ChapterManager{
		jsonManager: jsonManager,
		aiClient:    aiClient,
	}
}

// ChapterLimits defines event limits per campaign size
var ChapterLimits = map[string]int{
	"pequena": 3,
	"media":   4,
	"medio":   4,
	"extensa": 6,
}

// List returns all chapters in a campaign
func (c *ChapterManager) List(campaignName string) ([]*models.Chapter, error) {
	var chapters []*models.Chapter
	err := c.jsonManager.LoadInto(campaignName, ".capitulos.json", &chapters)
	if err != nil {
		return nil, err
	}
	if chapters == nil {
		return []*models.Chapter{}, nil
	}
	return chapters, nil
}

// Current returns the current active chapter
func (c *ChapterManager) Current(campaignName string) (*models.CurrentChapter, error) {
	var current models.CurrentChapter
	err := c.jsonManager.LoadInto(campaignName, ".capitulo_atual.json", &current)
	if err != nil {
		return nil, err
	}

	// Set defaults if file doesn't exist
	if current.Number == 0 {
		current = *models.NewCurrentChapter(1)
	}

	return &current, nil
}

// saveCurrent saves the current chapter
func (c *ChapterManager) saveCurrent(campaignName string, chapter *models.CurrentChapter) error {
	return c.jsonManager.Save(campaignName, ".capitulo_atual.json", chapter)
}

// EnsureOpening ensures the campaign has an opening prologue and initial context
func (c *ChapterManager) EnsureOpening(campaignName string) (string, error) {
	// Load context
	var context map[string]interface{}
	if err := c.jsonManager.LoadInto(campaignName, ".context.json", &context); err != nil {
		return "", err
	}

	description, _ := context["descricao"].(string)
	objective, _ := context["objetivo_principal"].(string)
	campaignNameFromCtx, _ := context["nome_campanha"].(string)
	if campaignNameFromCtx == "" {
		campaignNameFromCtx = campaignName
	}

	// Generate initial context if not present
	if _, hasContext := context["contexto_inicial"]; !hasContext && c.aiClient != nil {
		initialContext, err := c.aiClient.GenerateInitialContext(description, objective, campaignNameFromCtx)
		if err != nil {
			// Fallback to basic context if AI fails
			initialContext = fmt.Sprintf("Bem-vindos a %s. %s. Objetivo: %s", campaignNameFromCtx, description, objective)
		}
		context["contexto_inicial"] = initialContext
		if err := c.jsonManager.Save(campaignName, ".context.json", context); err != nil {
			return "", err
		}
	}

	initialContext, _ := context["contexto_inicial"].(string)

	// Check for prologue chapter
	chapters, err := c.List(campaignName)
	if err != nil {
		return "", err
	}

	hasPrologue := false
	for _, chapter := range chapters {
		if chapter.Type == "prologo" {
			hasPrologue = true
			break
		}
	}

	// Add prologue if not present
	if !hasPrologue {
		prologue := models.NewChapter(0, "Prólogo — o mestre abre o mundo")
		prologue.Summary = initialContext
		prologue.Events = []string{description, objective}
		prologue.Type = "prologo"
		prologue.ClosedAt = c.now()

		chapters = append([]*models.Chapter{prologue}, chapters...)
		if err := c.jsonManager.Save(campaignName, ".capitulos.json", chapters); err != nil {
			return "", err
		}
	}

	// Ensure current chapter exists
	var current models.CurrentChapter
	if err := c.jsonManager.LoadInto(campaignName, ".capitulo_atual.json", &current); err != nil {
		// Load history for current chapter events
		var history []string
		c.jsonManager.LoadInto(campaignName, ".hist.json", &history)

		current = *models.NewCurrentChapter(1)
		current.Events = history
		if err := c.saveCurrent(campaignName, &current); err != nil {
			return "", err
		}
	}

	return initialContext, nil
}

// Annotate adds an event to the current chapter
func (c *ChapterManager) Annotate(campaignName, outcome, narration string) (*models.Chapter, error) {
	current, err := c.Current(campaignName)
	if err != nil {
		return nil, err
	}

	event := outcome
	if narration != "" {
		event = fmt.Sprintf("%s // %s", outcome, narration)
	}

	current.AddEvent(event)
	if err := c.saveCurrent(campaignName, current); err != nil {
		return nil, err
	}

	// Check if chapter should be closed
	var context map[string]interface{}
	c.jsonManager.LoadInto(campaignName, ".context.json", &context)
	size, _ := context["tamanho_campanha"].(string)
	if size == "" {
		size = "media"
	}

	limit := ChapterLimits[strings.ToLower(size)]
	if limit == 0 {
		limit = 4
	}

	if len(current.Events) >= limit {
		return c.Close(campaignName)
	}

	return nil, nil
}

// Close closes the current chapter and creates a new one
func (c *ChapterManager) Close(campaignName string) (*models.Chapter, error) {
	current, err := c.Current(campaignName)
	if err != nil {
		return nil, err
	}

	if len(current.Events) == 0 {
		return nil, nil
	}

	// Load context for objective
	var context map[string]interface{}
	c.jsonManager.LoadInto(campaignName, ".context.json", &context)
	objective, _ := context["objetivo_principal"].(string)
	if objective == "" {
		objective = "a jornada"
	}

	var closedChapter *models.Chapter

	if c.aiClient != nil {
		// Use AI to generate chapter summary
		closedChapter, err = c.aiClient.CloseChapter(current.Number, current.Events, objective)
		if err != nil {
			return nil, fmt.Errorf("failed to close chapter with AI: %w", err)
		}
	} else {
		// Fallback: create basic chapter without AI
		closedChapter = models.NewChapter(current.Number, fmt.Sprintf("Capítulo %d", current.Number))
		closedChapter.Summary = fmt.Sprintf("Resumo do capítulo %d", current.Number)
		closedChapter.Events = current.Events
	}

	closedChapter.ClosedAt = c.now()

	// Save to chapters list
	chapters, err := c.List(campaignName)
	if err != nil {
		return nil, err
	}
	chapters = append(chapters, closedChapter)
	if err := c.jsonManager.Save(campaignName, ".capitulos.json", chapters); err != nil {
		return nil, err
	}

	// Create new current chapter
	newCurrent := models.NewCurrentChapter(closedChapter.Number + 1)
	if err := c.saveCurrent(campaignName, newCurrent); err != nil {
		return nil, err
	}

	return closedChapter, nil
}

// now returns current timestamp in the expected format
func (c *ChapterManager) now() string {
	return time.Now().Format("2006-01-02 15:04")
}

// GetMasterResponse gets AI master response to player actions
func (c *ChapterManager) GetMasterResponse(campaignName, action, outcome, mode string) (string, error) {
	if c.aiClient == nil {
		return "", fmt.Errorf("AI client not available")
	}
	return c.aiClient.MasterResponse(campaignName, action, outcome, mode)
}
