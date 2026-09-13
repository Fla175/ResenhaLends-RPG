package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Questionnaire handles interactive selection for terminal UI
type Questionnaire struct {
	scanner *bufio.Scanner
}

// NewQuestionnaire creates a new Questionnaire instance
func NewQuestionnaire() *Questionnaire {
	return &Questionnaire{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// SelectFromList presents a numbered list and returns the selected item
func (q *Questionnaire) SelectFromList(title string, options []string) (string, error) {
	fmt.Printf("\n\033[33m%s\033[0m\n", title)
	for i, option := range options {
		fmt.Printf("\033[32m[%d]\033[0m %s\n", i+1, option)
	}
	fmt.Print("\033[36mEscolha (número): \033[0m")

	q.scanner.Scan()
	input := strings.TrimSpace(q.scanner.Text())

	choice, err := strconv.Atoi(input)
	if err != nil {
		return "", fmt.Errorf("entrada inválida: %s", input)
	}

	if choice < 1 || choice > len(options) {
		return "", fmt.Errorf("escolha fora do intervalo: %d", choice)
	}

	return options[choice-1], nil
}

// SelectFromCategories presents categorized options
func (q *Questionnaire) SelectFromCategories(title string, categories map[string][]string) (string, error) {
	fmt.Printf("\n\033[33m%s\033[0m\n", title)

	categoryNames := make([]string, 0, len(categories))
	for name := range categories {
		categoryNames = append(categoryNames, name)
	}

	// Select category
	category, err := q.SelectFromList("Selecione uma categoria:", categoryNames)
	if err != nil {
		return "", err
	}

	// Select option from category
	options := categories[category]
	if len(options) == 0 {
		return "", fmt.Errorf("categoria vazia: %s", category)
	}

	return q.SelectFromList(fmt.Sprintf("Selecione (%s):", category), options)
}

// Confirm asks for yes/no confirmation
func (q *Questionnaire) Confirm(message string) bool {
	fmt.Printf("\033[33m%s (s/n): \033[0m", message)
	q.scanner.Scan()
	input := strings.ToLower(strings.TrimSpace(q.scanner.Text()))
	return input == "s" || input == "sim" || input == "y" || input == "yes"
}

// GetInput gets text input from user
func (q *Questionnaire) GetInput(prompt string) string {
	fmt.Printf("\033[36m%s: \033[0m", prompt)
	q.scanner.Scan()
	return strings.TrimSpace(q.scanner.Text())
}

// GetNumber gets numeric input from user
func (q *Questionnaire) GetNumber(prompt string, defaultValue int) int {
	fmt.Printf("\033[36m%s (padrão: %d): \033[0m", prompt, defaultValue)
	q.scanner.Scan()
	input := strings.TrimSpace(q.scanner.Text())

	if input == "" {
		return defaultValue
	}

	num, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	}
	return num
}