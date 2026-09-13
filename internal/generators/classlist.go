package generators

// ClassList contains categorized class options
var ClassCategories = map[string][]string{
	"Combate": {
		"Guerreiro", "Paladino", "Bárbaro", "Cavaleiro", "Ladino", "Ranger", "Mercenário",
		"Gladiador", "Samurai", "Ninja", "Espadachim", "Arqueiro", "Lanceiro", "Mestre de Armas",
		"Duelista", "Berserker", "Guarda Real", "Cavaleiro Negro", "Assassino", "Moto-boy",
	},
	"Magia": {
		"Mago", "Bruxo", "Clérigo", "Druida", "Feiticeiro", "Bardo", "Místico", "Necromante",
		"Invocador", "Alquimista", "Encantador", "Elementalista", "Ilusionista", "Adivinho",
		"Xamã", "Conjurador", "Teurgista", "Arcanista", "Caneta Azul",
	},
	"Suporte": {
		"Curandeiro", "Bardo", "Clérigo", "Druida", "Paladino", "Monge", "Sacerdote",
		"Inspetor", "Conselheiro", "Diplomata", "Escriba", "Historiador", "Sábio",
		"Mentor", "Guardião", "Protetor", "Boticário", "Herbalista",
	},
	"Meme": {
		"Moto-boy", "Lambizome", "Caneta Azul", "Influencer", "Streamer", "Youtuber",
		"Prompt Engineer", "Hacker de ChatGPT", "Gerente de Produto sem roadmap", "DevOps que quebra produção",
		"QA que não testa", "Product Owner sem backlog", "Scrum Master sem scrum",
		"Gerente de Tráfego", "SEO Specialist", "Growth Hacker", "Full Stack Developer",
		"Tech Lead sem código", "CTO de startup sem produto",
	},
	"Especializadas": {
		"Artífice", "Inventor", "Engenheiro", "Alquimista", "Ferramenteiro", "Jewel",
		"Explorador", "Cartógrafo", "Navegador", "Batedor", "Investigador", "Detetive",
		"Coveiro", "Taberneiro", "Mercador", "Ladrão de Guilda", "Espião", "Batedor de Cartas",
		"Astrologo", "Runemaster", "Cavalheiro da Ordem", "Caçador de Demônios",
	},
}

// GetAllClasses returns all classes as a flat list
func GetAllClasses() []string {
	var allClasses []string
	for _, classes := range ClassCategories {
		allClasses = append(allClasses, classes...)
	}
	return allClasses
}

// GetClassCategories returns the class categories
func GetClassCategories() map[string][]string {
	return ClassCategories
}