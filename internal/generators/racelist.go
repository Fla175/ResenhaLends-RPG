package generators

// RaceList contains categorized race options
var RaceCategories = map[string][]string{
	"Clássicas": {
		"Humano", "Elfo", "Anão", "Halfling", "Meio-Orc", "Draconato", "Tiefling", "Gnomo",
		"Goliath", "Aasimar", "Firbolg", "Tabaxi", "Kenku", "Lizardfolk", "Changeling",
	},
	"Exóticas": {
		"Warforged", "Centauro", "Minotauro", "Goblinoide", "Orc", "Drow", "Duérgar",
		"Githyanki", "Githzerai", "Satir", "Tritão", "Genasi do Fogo", "Genasi da Água",
		"Genasi da Terra", "Genasi do Ar", "Shifter", "Leonino", "Loxodonte", "Vedalken",
		"Simic Híbrido", "Kalashtar", "Aerie", "Gnomo das Profundezas", "Elfo da Floresta",
		"Anão da Montanha", "Meio-Elfo", "Meio-Anão", "Tiflin", "Dragonborn",
	},
	"Meme": {
		"Jacobiano", "Caneta Azul", "Programador Java", "Gerente de Produto", "Influencer",
		"Youtuber", "Streamer", "Hacker de ChatGPT", "Entusiasta de Crypto", "NFT Bro",
		"Metaverso Habitante", "Prompt Engineer", "Data Scientist Junior", "DevOps Frustrado",
		"QA que testa tudo em produção", "Product Owner sem backlog", "Scrum Master sem scrum",
	},
	"Especiais": {
		"Construto", "Espírito da Natureza", "Semi-Deus", "Celestial", "Infernal",
		"Aberração", "Fada", "Elemental", "Gigante", "Monstro",
	},
}

// GetAllRaces returns all races as a flat list
func GetAllRaces() []string {
	var allRaces []string
	for _, races := range RaceCategories {
		allRaces = append(allRaces, races...)
	}
	return allRaces
}

// GetRaceCategories returns the race categories
func GetRaceCategories() map[string][]string {
	return RaceCategories
}