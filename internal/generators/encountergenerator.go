package generators

import (
	"fmt"
	"math/rand"
	"resenhalends/internal/managers"
	"strings"
	"time"
)

// EncounterGenerator handles random encounter generation
type EncounterGenerator struct {
	jsonManager *managers.JSONManager
	rand        *rand.Rand
}

// NewEncounterGenerator creates a new EncounterGenerator
func NewEncounterGenerator(jsonManager *managers.JSONManager) *EncounterGenerator {
	return &EncounterGenerator{
		jsonManager: jsonManager,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

var (
	encounterTypes = []string{
		"Combate", "Emboscada", "Armadilha", "Descoberta", "Clima Extremo",
		"Encontro Social", "Enigma", "NPC Perdido", "Tesouro Oculto",
		"Obstáculo Natural", "Sinal Misterioso",
	}

	threats = map[string][]string{
		"fácil": {
			"Bando de Goblins", "Lobos Selvagens", "Bandidos Novatos",
			"Ratos Gigantes", "Slime Ácido", "Cultistas Aprendizes", "Kobolds Ladinos",
		},
		"normal": {
			"Orcs Furiosos", "Esqueletos Armados", "Urso-Coruja",
			"Aranhas Gigantes", "Gárgula de Pedra", "Elemental da Terra",
			"Manticora", "Feiticeiro Renegado",
		},
		"hardcore": {
			"Troll de Pedra", "Vampiro Ancião", "Dragão Jovem",
			"Lich Enfraquecido", "Devorador de Mentes", "Demônio da Fornalha",
			"Beholder (Contemplador)", "Hydra das Cavernas",
		},
	}

	locations = []string{
		// Forests and Nature
		"em uma clareira densa", "atrás de folhagens espessas", "sob as raízes de uma árvore gigante",
		"em um pântano lamacento e fétido", "no coração de um bosque sussurrante",
		// Roads and Passages
		"no meio da estrada", "em uma ponte de pedra antiga", "em uma bifurcação deserta da trilha",
		"no cruzamento de estradas abandonadas", "em um túnel estreito escavado na rocha",
		// Ruins and Structures
		"perto de ruínas desmoronadas", "no interior de uma cabana em ruínas",
		"nos escombros de uma torre de vigia", "no pátio de um castelo esquecido",
		"na entrada de um templo sagrado profanado",
		// Terrain and Caves
		"à beira de um penhasco íngreme", "dentro de uma caverna úmida", "no topo de uma colina enevoada",
		"no fundo de um desfiladeiro rochoso", "na entrada de uma mina abandonada",
		// Water and Coast
		"próximo a um rio turbulento", "à margem de um lago de águas escuras",
		"ao lado de uma cachoeira barulhenta", "em uma praia de areia negra",
		// Dark and Mysterious
		"em um cemitério abandonado", "no centro de um círculo de pedras rúnicas",
		"onde a névoa impede de enxergar um palmo à frente",
	}
)

// GenerateEncounter creates a random encounter for a campaign
func (e *EncounterGenerator) GenerateEncounter(campaignName string) (string, error) {
	// Load campaign rules
	var rules map[string]interface{}
	err := e.jsonManager.LoadInto(campaignName, ".regras.json", &rules)
	if err != nil {
		return "", err
	}

	difficulty := "normal"
	if d, ok := rules["dificuldade"].(string); ok {
		difficulty = strings.ToLower(d)
	}

	encounterType := e.randomChoice(encounterTypes)
	threatList := threats[difficulty]
	if threatList == nil {
		threatList = threats["normal"]
	}
	threat := e.randomChoice(threatList)
	location := e.randomChoice(locations)

	var event string
	switch encounterType {
	case "Descoberta":
		event = fmt.Sprintf("O grupo encontra um baú trancado ou santuário antigo %s.", location)
	case "Armadilha":
		event = fmt.Sprintf("Alguém pisa em uma armadilha escondida %s.", location)
	case "Clima Extremo":
		event = fmt.Sprintf("Uma tempestade violenta ou nevasca repentina atinge o grupo %s.", location)
	default:
		event = fmt.Sprintf("%s: Vocês são surpreendidos por um %s %s.", encounterType, threat, location)
	}

	fmt.Printf("\n\033[31m⚠️ ENCONTRO GERADO: %s\033[0m\n", event)

	// Add to campaign history
	var history []string
	e.jsonManager.LoadInto(campaignName, ".hist.json", &history)
	history = append(history, fmt.Sprintf("Evento Aleatório: %s", event))
	if err := e.jsonManager.Save(campaignName, ".hist.json", history); err != nil {
		return "", err
	}

	return event, nil
}

func (e *EncounterGenerator) randomChoice(choices []string) string {
	return choices[e.rand.Intn(len(choices))]
}
