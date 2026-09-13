package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"resenhalends/internal/models"
	"time"
)

// PythonClient communicates with the Python Why-Man AI backend
type PythonClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewPythonClient creates a new PythonClient
func NewPythonClient(baseURL string) *PythonClient {
	return &PythonClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateInitialContext requests initial context from Python AI
func (p *PythonClient) GenerateInitialContext(description, objective, campaignName string) (string, error) {
	request := map[string]interface{}{
		"action":     "generate_initial_context",
		"descricao":  description,
		"objetivo":   objective,
		"campanha":   campaignName,
	}

	response, err := p.sendRequest(request)
	if err != nil {
		return "", err
	}

	if context, ok := response["contexto_inicial"].(string); ok {
		return context, nil
	}

	return "", fmt.Errorf("invalid response from AI backend")
}

// CloseChapter requests chapter summary from Python AI
func (p *PythonClient) CloseChapter(number int, events []string, objective string) (*models.Chapter, error) {
	request := map[string]interface{}{
		"action":    "close_chapter",
		"numero":   number,
		"eventos":  events,
		"objetivo": objective,
	}

	response, err := p.sendRequest(request)
	if err != nil {
		return nil, err
	}

	chapter := models.NewChapter(number, fmt.Sprintf("Capítulo %d", number))
	if titulo, ok := response["titulo"].(string); ok {
		chapter.Title = titulo
	}
	if resumo, ok := response["resumo"].(string); ok {
		chapter.Summary = resumo
	}
	chapter.Events = events

	return chapter, nil
}

// ProcessAction processes an action through the AI
func (p *PythonClient) ProcessAction(campaignName, action string, diceRoll int) (string, error) {
	request := map[string]interface{}{
		"action":       "process_action",
		"campanha":     campaignName,
		"acao":         action,
		"resultado_dado": diceRoll,
	}

	response, err := p.sendRequest(request)
	if err != nil {
		return "", err
	}

	if desfecho, ok := response["desfecho"].(string); ok {
		return desfecho, nil
	}

	return "", fmt.Errorf("invalid response from AI backend")
}

// Narrate requests narration from the AI
func (p *PythonClient) Narrate(campaignName, outcome string) (string, error) {
	request := map[string]interface{}{
		"action":    "narrate",
		"campanha":  campaignName,
		"desfecho":  outcome,
	}

	response, err := p.sendRequest(request)
	if err != nil {
		return "", err
	}

	if narracao, ok := response["narracao"].(string); ok {
		return narracao, nil
	}

	return "", fmt.Errorf("invalid response from AI backend")
}

// MasterResponse generates master response to player actions
func (p *PythonClient) MasterResponse(campaignName, action, outcome string, mode string) (string, error) {
	request := map[string]interface{}{
		"action":     "master_response",
		"campanha":   campaignName,
		"acao":       action,
		"desfecho":   outcome,
		"modo":       mode,
	}

	response, err := p.sendRequest(request)
	if err != nil {
		return "", err
	}

	if resposta, ok := response["resposta"].(string); ok {
		return resposta, nil
	}

	return "", fmt.Errorf("invalid response from AI backend")
}

// CheckConnection verifies if the Python AI backend is available
func (p *PythonClient) CheckConnection() (bool, string) {
	request := map[string]interface{}{
		"action": "ping",
	}

	response, err := p.sendRequest(request)
	if err != nil {
		return false, fmt.Sprintf("AI backend unavailable: %v", err)
	}

	if status, ok := response["status"].(string); ok {
		return true, status
	}

	return true, "AI backend connected"
}

// sendRequest sends a request to the Python backend
func (p *PythonClient) sendRequest(data map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := p.httpClient.Post(p.baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if errorMsg, ok := result["error"].(string); ok {
			return nil, fmt.Errorf("AI backend error: %s", errorMsg)
		}
		return nil, fmt.Errorf("AI backend returned status %d", resp.StatusCode)
	}

	return result, nil
}

// MockClient is a fallback client when Python backend is unavailable
type MockClient struct{}

// NewMockClient creates a new MockClient
func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) GenerateInitialContext(description, objective, campaignName string) (string, error) {
	return fmt.Sprintf("Contexto inicial para %s: %s. Objetivo: %s", campaignName, description, objective), nil
}

func (m *MockClient) CloseChapter(number int, events []string, objective string) (*models.Chapter, error) {
	chapter := models.NewChapter(number, fmt.Sprintf("Capítulo %d", number))
	chapter.Summary = fmt.Sprintf("Resumo do capítulo %d com %d eventos", number, len(events))
	chapter.Events = events
	return chapter, nil
}

func (m *MockClient) ProcessAction(campaignName, action string, diceRoll int) (string, error) {
	return fmt.Sprintf("Ação '%s' processada com resultado %d", action, diceRoll), nil
}

func (m *MockClient) Narrate(campaignName, outcome string) (string, error) {
	return fmt.Sprintf("Narração: %s", outcome), nil
}

func (m *MockClient) CheckConnection() (bool, string) {
	return false, "Using mock AI client (Python backend unavailable)"
}

func (m *MockClient) MasterResponse(campaignName, action, outcome, mode string) (string, error) {
	return fmt.Sprintf("Mestre responde à ação '%s': %s", action, outcome), nil
}
