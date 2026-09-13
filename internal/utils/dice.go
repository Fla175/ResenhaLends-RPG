package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// Dice handles dice rolling operations
type Dice struct {
	showOutput bool
	rand       *rand.Rand
	shouldRoll bool // Flag to control if dice should roll
}

// NewDice creates a new Dice instance
func NewDice(showOutput bool) *Dice {
	return &Dice{
		showOutput: showOutput,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
		shouldRoll: true, // Default to rolling dice
	}
}

// SetShouldRoll sets whether dice should roll
func (d *Dice) SetShouldRoll(shouldRoll bool) {
	d.shouldRoll = shouldRoll
}

// Roll rolls a die with the specified number of faces
func (d *Dice) Roll(faces int) int {
	if !d.shouldRoll {
		if d.showOutput {
			fmt.Println("🎲 Dado anulado (ação rápida)")
		}
		return 0 // Return 0 for no-roll actions
	}

	if faces < 2 {
		panic("number of faces must be at least 2")
	}

	result := d.rand.Intn(faces) + 1

	if d.showOutput {
		var color string
		if result == faces {
			color = "\033[32m" // Green
		} else if result == 1 {
			color = "\033[31m" // Red
		} else {
			color = "\033[33m" // Yellow
		}
		reset := "\033[0m"
		fmt.Printf("%s🎲 D%d: %d%s\n", color, faces, result, reset)
	}

	return result
}

// RollMultiple rolls multiple dice and returns the sum
func (d *Dice) RollMultiple(faces, count int) int {
	if !d.shouldRoll {
		return 0
	}

	total := 0
	for i := 0; i < count; i++ {
		total += d.Roll(faces)
	}
	return total
}

// RollWithModifier rolls a die and adds a modifier
func (d *Dice) RollWithModifier(faces, modifier int) int {
	if !d.shouldRoll {
		return modifier
	}

	result := d.Roll(faces)
	return result + modifier
}

// ShouldRoll returns whether dice should roll
func (d *Dice) ShouldRoll() bool {
	return d.shouldRoll
}
