package generators

import (
	"fmt"
	"math/rand"
	"resenhalends/internal/managers"
	"resenhalends/internal/models"
	"time"
)

// NPCGenerator handles NPC generation
type NPCGenerator struct {
	jsonManager *managers.JSONManager
	rand        *rand.Rand
}

// NewNPCGenerator creates a new NPCGenerator
func NewNPCGenerator(jsonManager *managers.JSONManager) *NPCGenerator {
	return &NPCGenerator{
		jsonManager: jsonManager,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

var (
	names = []string{
		"Thalor", "Elara", "Bram", "Kael", "Lyra", "Gorn", "Silas", "Maeve", "Doran", "Vex",
		"Eldrin", "Mira", "Tavish", "Kira", "Borin", "Zalriel", "Faelar", "Yvaine", "Garrick", "Sariel",
		"Fenric", "Liora", "Dax", "Valen", "Rowan", "Talia", "Kaelen", "Morwen", "Orin", "Reyna",
		"Zephyr", "Aveline", "Caelum", "Elowen", "Gideon", "Hadrian", "Isolde", "Jarek", "Malakor", "Theron",
		"Ignis", "Storm", "Thurn", "Zorvath", "Kallista", "Brottor", "Fira", "Keth", "Milo", "Zanna",
		"Xalthor", "Ruk", "Garnet", "Marrow", "Echo", "Rusty", "Torinn", "Nisha", "Aurelius", "Bryn",
	}

	sexes = []string{"Masculino", "Feminino"}

	// Use centralized race list
	races = GetAllRaces()

	occupations = []string{
		"Ferreiro", "Mercador Suspeito", "Guarda da Cidade", "Mago Andarilho", "Ladino",
		"Alquimista Excêntrico", "Coveiro Ranzinza", "Taberneiro Falador", "Caçador de Recompensas",
		"Sacerdote Corrupto", "Bardo Fracassado", "Ervanária da Vila", "Mensageiro Real",
		"Navegador Aposentado", "Estudante de Arqueologia", "Guarda-costas Silencioso",
		"Cozinheiro de Taverna", "Contrabandista de Carga", "Alfaiate da Alta Sociedade", "Eremita da Floresta",
	}

	motivations = []string{
		"Busca vingança", "Quer ficar rico rápido", "Foge do passado", "Procura um artefato perdido",
		"Protege um segredo antigo", "Busca o perdão por um crime", "Quer provar seu valor a uma guilda",
		"Deseja curar uma doença incurável", "Pesquisa uma lenda esquecida", "Busca fama e glória no combate",
		"Quer reencontrar um amor perdido", "Tenta saldar uma dívida de jogo", "Protege sua família ameaçada",
		"Quer dominar uma magia proibida", "Foge de um casamento arranjado", "Deseja recuperar a honra da família",
		"Busca entender uma visão profética", "Quer derrubar um tirano local", "Explora o mundo por puro tédio",
		"Coleta espécimes raros para estudo",
	}

	secrets = []string{
		"É um espião disfarçado", "Tem uma maldição oculta", "Sabe onde está o chefe final",
		"Deve dinheiro à guilda das sombras", "Na verdade, é um lobisomem", "Roubou a identidade de um nobre falecido",
		"Serve secretamente a uma divindade maligna", "Possui um mapa para o covil do dragão",
		"Foi o responsável por incendiar a vila vizinha", "É o herdeiro legítimo do trono que foi deposto",
		"Sua riqueza vem de pactos com demônios", "Está envenenando lentamente o líder local",
		"Vendeu informações sigilosas para os inimigos", "É uma criatura mágica disfarçada de humano",
		"Está fingindo ter amnésia para esconder um crime", "Possui a última cópia de um grimório proibido",
		"Foi enviado para sabotar a missão do grupo", "É o verdadeiro assassino do antigo rei",
		"Sua linhagem sanguínea ativa portais mágicos", "Guarda a chave que liberta o grande vilão",
	}
)

// GenerateNPC creates a new random NPC
func (n *NPCGenerator) GenerateNPC(campaignName string) (*models.NPC, error) {
	npc := &models.NPC{
		Name:       n.randomChoice(names),
		Sex:        n.randomChoice(sexes),
		Race:       n.randomChoice(races),
		Occupation: n.randomChoice(occupations),
		Motivation: n.randomChoice(motivations),
		Secret:     n.randomChoice(secrets),
	}

	// Save to campaign's NPC database
	var npcs map[string]*models.NPC
	err := n.jsonManager.LoadInto(campaignName, ".npc.json", &npcs)
	if err != nil {
		return nil, err
	}

	if npcs == nil {
		npcs = make(map[string]*models.NPC)
	}

	npcs[npc.Name] = npc
	if err := n.jsonManager.Save(campaignName, ".npc.json", npcs); err != nil {
		return nil, err
	}

	n.printNPC(npc)
	return npc, nil
}

func (n *NPCGenerator) randomChoice(choices []string) string {
	return choices[n.rand.Intn(len(choices))]
}

func (n *NPCGenerator) printNPC(npc *models.NPC) {
	fmt.Println("\n\033[35m--- NOVO NPC CRIADO ---\033[0m")
	fmt.Printf("\033[36mNome:\033[0m %s\n", npc.Name)
	fmt.Printf("\033[36mSexo:\033[0m %s\n", npc.Sex)
	fmt.Printf("\033[36mRaça:\033[0m %s\n", npc.Race)
	fmt.Printf("\033[36mOcupação:\033[0m %s\n", npc.Occupation)
	fmt.Printf("\033[36mMotivação:\033[0m %s\n", npc.Motivation)
	fmt.Printf("\033[36mSegredo:\033[0m %s\n", npc.Secret)
	fmt.Println("\033[35m-----------------------\033[0m")
}
