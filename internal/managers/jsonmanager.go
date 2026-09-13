package managers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	BaseFolder     = "campanhas"
	ExportFolder   = "exportacoes"
)

// JSONManager handles JSON file operations for campaigns
type JSONManager struct {
	baseFolder string
}

// NewJSONManager creates a new JSONManager
func NewJSONManager() *JSONManager {
	return &JSONManager{
		baseFolder: BaseFolder,
	}
}

// GetPath returns the full path for a campaign file
func (j *JSONManager) GetPath(campaignName string, filename string) string {
	if filename == "" {
		return filepath.Join(j.baseFolder, campaignName)
	}
	return filepath.Join(j.baseFolder, campaignName, filename)
}

// ListCampaigns returns a list of all campaign names
func (j *JSONManager) ListCampaigns() []string {
	if _, err := os.Stat(j.baseFolder); os.IsNotExist(err) {
		return []string{}
	}

	entries, err := os.ReadDir(j.baseFolder)
	if err != nil {
		return []string{}
	}

	var campaigns []string
	for _, entry := range entries {
		if entry.IsDir() {
			campaigns = append(campaigns, entry.Name())
		}
	}

	return campaigns
}

// SetupCampaign creates a new campaign with default files
func (j *JSONManager) SetupCampaign(name, difficulty, description, objective, size string) error {
	campaignPath := j.GetPath(name, "")
	if err := os.MkdirAll(campaignPath, 0755); err != nil {
		return err
	}

	baseFiles := map[string]interface{}{
		".context.json": map[string]interface{}{
			"nome_campanha":      name,
			"descricao":          description,
			"objetivo_principal":  objective,
			"tamanho_campanha":   size,
		},
		".regras.json": map[string]interface{}{
			"dificuldade":      difficulty,
			"magia_permitida":  true,
		},
		".jogadores.json": []interface{}{},
		".hist.json":       []interface{}{},
		".npc.json":        map[string]interface{}{},
		".monstros.json":   map[string]interface{}{},
		".notas_mestre.json": map[string]interface{}{
			"anotacoes": "",
			"segredos":  "",
		},
		"resumo.json": map[string]interface{}{
			"historia":             "",
			"dialogos_importantes": []interface{}{},
		},
		".capitulos.json": []interface{}{},
		".capitulo_atual.json": map[string]interface{}{
			"numero":  1,
			"eventos": []interface{}{},
			"aberto":  true,
		},
	}

	for filename, data := range baseFiles {
		filePath := j.GetPath(name, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := j.Save(name, filename, data); err != nil {
				return fmt.Errorf("failed to save %s: %w", filename, err)
			}
		}
	}

	j.MigrateBestiary(name)
	return nil
}

// Save writes data to a JSON file
func (j *JSONManager) Save(campaignName, filename string, data interface{}) error {
	filePath := j.GetPath(campaignName, filename)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

// Load reads data from a JSON file
func (j *JSONManager) Load(campaignName, filename string) (interface{}, error) {
	filePath := j.GetPath(campaignName, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// LoadInto reads JSON file into a specific struct
func (j *JSONManager) LoadInto(campaignName, filename string, target interface{}) error {
	filePath := j.GetPath(campaignName, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, target)
}

// DeleteCampaign removes a campaign directory
func (j *JSONManager) DeleteCampaign(name string) error {
	campaignPath := j.GetPath(name, "")
	if _, err := os.Stat(campaignPath); os.IsNotExist(err) {
		return fmt.Errorf("campaign '%s' does not exist", name)
	}

	return os.RemoveAll(campaignPath)
}

// MigrateBestiary migrates old bestiary format to new format
func (j *JSONManager) MigrateBestiary(campaignName string) {
	oldData, _ := j.Load(campaignName, ".mostros.json")
	newData, _ := j.Load(campaignName, ".monstros.json")

	if oldData != nil && newData == nil {
		j.Save(campaignName, ".monstros.json", oldData)
	}

	if newData == nil {
		j.Save(campaignName, ".monstros.json", oldData)
	}
}

// ExportCampaign exports a campaign to a zip file
func (j *JSONManager) ExportCampaign(name string, destinationDir string) (string, error) {
	sourcePath := j.GetPath(name, "")
	if info, err := os.Stat(sourcePath); err != nil || !info.IsDir() {
		return "", fmt.Errorf("campaign '%s' does not exist", name)
	}

	if destinationDir == "" {
		destinationDir = ExportFolder
	}

	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102_1504")
	baseName := filepath.Join(destinationDir, fmt.Sprintf("%s_%s", name, timestamp))
	zipPath := baseName + ".zip"

	// Create zip file
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	// Walk through campaign directory and add to zip
	// Note: This is a simplified implementation. For production, use archive/zip
	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// For simplicity, we'll just copy files. In production, use archive/zip
		// This is a placeholder - you'd need to implement proper zip creation
		return nil
	})

	if err != nil {
		return "", err
	}

	return zipPath, nil
}

// ImportCampaign imports a campaign from a zip file
func (j *JSONManager) ImportCampaign(zipPath string) ([]string, error) {
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("zip file not found")
	}

	if err := os.MkdirAll(j.baseFolder, 0755); err != nil {
		return nil, err
	}

	// Extract zip file
	// Note: This is a simplified implementation. For production, use archive/zip
	// This is a placeholder - you'd need to implement proper zip extraction

	campaigns := j.ListCampaigns()
	if len(campaigns) > 0 {
		j.MigrateBestiary(campaigns[len(campaigns)-1])
	}

	return campaigns, nil
}
