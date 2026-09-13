package utils

import (
	"fmt"
	"math"
)

// LevelSystem handles character progression and XP calculations
type LevelSystem struct{}

// NewLevelSystem creates a new LevelSystem
func NewLevelSystem() *LevelSystem {
	return &LevelSystem{}
}

// XPForLevel calculates the XP required for a given level
func (l *LevelSystem) XPForLevel(level int) int {
	if level <= 1 {
		return 0
	}
	// Exponential XP curve: 100 * (level ^ 1.5)
	return int(100 * math.Pow(float64(level), 1.5))
}

// XPToNextLevel calculates XP needed to reach the next level
func (l *LevelSystem) XPToNextLevel(currentLevel int) int {
	return l.XPForLevel(currentLevel + 1) - l.XPForLevel(currentLevel)
}

// AddXP adds XP to a character and checks for level up
func (l *LevelSystem) AddXP(currentXP, currentLevel int, xpGained int) (newXP, newLevel int, leveledUp bool) {
	newXP = currentXP + xpGained
	newLevel = currentLevel
	leveledUp = false

	// Check for level ups
	for l.XPForLevel(newLevel+1) <= newXP {
		newLevel++
		leveledUp = true
	}

	return newXP, newLevel, leveledUp
}

// CalculateLevelBonus returns the bonus modifier based on level
func (l *LevelSystem) CalculateLevelBonus(level int) int {
	return (level - 1) / 2 // +1 bonus every 2 levels
}

// GetLevelProgress returns the progress percentage towards next level
func (l *LevelSystem) GetLevelProgress(currentXP, currentLevel int) float64 {
	if currentLevel >= 20 { // Max level
		return 100.0
	}

	currentLevelXP := l.XPForLevel(currentLevel)
	nextLevelXP := l.XPForLevel(currentLevel + 1)
	xpInRange := nextLevelXP - currentLevelXP
	xpProgress := currentXP - currentLevelXP

	if xpInRange == 0 {
		return 0.0
	}

	progress := (float64(xpProgress) / float64(xpInRange)) * 100.0
	if progress > 100.0 {
		progress = 100.0
	}
	return progress
}

// GetXPNeeded returns the XP needed to reach the next level
func (l *LevelSystem) GetXPNeeded(currentXP, currentLevel int) int {
	if currentLevel >= 20 {
		return 0
	}

	nextLevelXP := l.XPForLevel(currentLevel + 1)
	return nextLevelXP - currentXP
}

// FormatXPProgress returns a formatted string showing XP progress
func (l *LevelSystem) FormatXPProgress(currentXP, currentLevel int) string {
	if currentLevel >= 20 {
		return "Nível Máximo"
	}

	progress := l.GetLevelProgress(currentXP, currentLevel)
	nextLevelXP := l.XPForLevel(currentLevel + 1)
	xpNeeded := l.GetXPNeeded(currentXP, currentLevel)

	return fmt.Sprintf("Nível %d | XP: %d/%d (%.1f%%) - Faltam: %d XP", currentLevel, currentXP, nextLevelXP, progress, xpNeeded)
}