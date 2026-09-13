package main

import (
	"bufio"
	"fmt"
	"os"
	"resenhalends/internal/ai"
	"resenhalends/internal/generators"
	"resenhalends/internal/managers"
	"resenhalends/internal/models"
	"resenhalends/internal/utils"
	"strings"
)

var (
	jsonManager       *managers.JSONManager
	playerManager     *managers.PlayerManager
	chapterManager    *managers.ChapterManager
	codexManager      *managers.CodexManager
	npcGenerator      *generators.NPCGenerator
	encounterGen      *generators.EncounterGenerator
	weaknessGenerator *generators.WeaknessGenerator
	dice              *utils.Dice
	aiClient          managers.AIClient
	scanner           *bufio.Scanner
	questionnaire     *utils.Questionnaire
)

func init() {
	jsonManager = managers.NewJSONManager()
	dice = utils.NewDice(true)
	
	// Try to connect to Python AI backend, fall back to mock
	pythonClient := ai.NewPythonClient("http://localhost:8080/api")
	if connected, _ := pythonClient.CheckConnection(); connected {
		aiClient = pythonClient
		fmt.Println("\033[32m✓ AI backend connected\033[0m")
	} else {
		aiClient = ai.NewMockClient()
		fmt.Println("\033[33m⚠ Using mock AI client (Python backend unavailable)\033[0m")
	}
	
	playerManager = managers.NewPlayerManager(jsonManager, dice)
	chapterManager = managers.NewChapterManager(jsonManager, aiClient)
	codexManager = managers.NewCodexManager(jsonManager)
	npcGenerator = generators.NewNPCGenerator(jsonManager)
	encounterGen = generators.NewEncounterGenerator(jsonManager)
	weaknessGenerator = generators.NewWeaknessGenerator()

	scanner = bufio.NewScanner(os.Stdin)
	questionnaire = utils.NewQuestionnaire()
}

func main() {
	displayLogo()
	
	for {
		choice := mainMenu()
		switch choice {
		case 1:
			createCampaign()
		case 2:
			loadCampaign()
		case 3:
			manageRules()
		case 4:
			exportCampaign()
		case 5:
			deleteCampaign()
		case 6:
			fmt.Println("\033[32mAté a próxima!\033[0m")
			return
		}
	}
}

func displayLogo() {
	fmt.Println("\033[36m    =================================================")
	fmt.Println(" ____                      _             _                  _")
	fmt.Println("|  _ \\ ___  ___  ___ _ __ | |__   ____  | |    ___ _ __   __| |____")
	fmt.Println("| |_) / _ \\ / __|/ _ \\ '_ \\| '_ \\ | '_ \\ | |   / _ \\ '_ \\ / _` / ___|")
	fmt.Println("|  _ <  __/\\__ \\  __/ | | | | | || |_| || |__|  __/ | | | (_| \\__  \\")
	fmt.Println("|_| \\_\\___||___/\\___|_| |_|_| |_||_| |_||____\\___|__| |_|\\__,_|____/")
	fmt.Println("                                                     ")
	fmt.Println("        RPG SELF-HOSTED - MESTRE DE MESA PARA RPG")
	fmt.Println("    =================================================\033[0m")
}

func mainMenu() int {
	fmt.Println("\n\033[33mEscolha o seu destino:\033[0m")
	fmt.Println("\033[32m1.\033[0m Nova Campanha")
	fmt.Println("\033[32m2.\033[0m Carregar Campanha")
	fmt.Println("\033[32m3.\033[0m Gerenciar Regras")
	fmt.Println("\033[32m4.\033[0m Exportar Campanha")
	fmt.Println("\033[32m5.\033[0m Excluir Campanha")
	fmt.Println("\033[32m6.\033[0m Sair")
	
	fmt.Print("\033[36m$\033[0m ")
	scanner.Scan()
	input := scanner.Text()
	
	var choice int
	fmt.Sscanf(input, "%d", &choice)
	if choice < 1 || choice > 6 {
		return 6
	}
	return choice
}

func chooseCampaign() string {
	campaigns := jsonManager.ListCampaigns()
	if len(campaigns) == 0 {
		fmt.Println("\033[31mNenhuma campanha encontrada.\033[0m")
		return ""
	}
	
	fmt.Println("\nCampanhas disponíveis:")
	for i, name := range campaigns {
		fmt.Printf("[%d] %s\n", i+1, name)
	}
	
	fmt.Print("\nNúmero ou Nome da campanha: ")
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())
	
	// Try as number
	var num int
	if _, err := fmt.Sscanf(input, "%d", &num); err == nil {
		if num >= 1 && num <= len(campaigns) {
			return campaigns[num-1]
		}
	}
	
	// Try as name
	for _, name := range campaigns {
		if strings.ToLower(name) == strings.ToLower(input) {
			return name
		}
	}
	
	fmt.Println("\033[31mInválido.\033[0m")
	return ""
}

func createCampaign() {
	fmt.Print("Nome da campanha: ")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ToLower(name)
	
	if name == "" {
		return
	}
	
	fmt.Print("Descrição da campanha: ")
	scanner.Scan()
	description := scanner.Text()
	
	fmt.Print("Objetivo principal: ")
	scanner.Scan()
	objective := scanner.Text()

	// Use questionnaire for difficulty selection
	difficulty, err := questionnaire.SelectFromList("Dificuldade:", []string{"Fácil", "Normal", "Hardcore"})
	if err != nil {
		difficulty = "Normal"
	}

	// Use questionnaire for size selection
	sizeSelection, err := questionnaire.SelectFromList("Tamanho:", []string{"Pequena", "Média", "Extensa"})
	if err != nil {
		sizeSelection = "Média"
	}

	size := "media"
	switch sizeSelection {
	case "Pequena":
		size = "pequena"
	case "Extensa":
		size = "extensa"
	}
	
	if err := jsonManager.SetupCampaign(name, difficulty, description, objective, size); err != nil {
		fmt.Printf("\033[31mErro ao criar campanha: %v\033[0m\n", err)
		return
	}

	// Configure master response mode
	masterMode, err := questionnaire.SelectFromList("Modo de resposta do mestre:", []string{
		"Automático (IA responde sempre)",
		"Semi-automático (IA sugere, você aprova)",
		"Manual (vê sugestão e edita)",
	})
	if err != nil {
		masterMode = "Semi-automático (IA sugere, você aprova)"
	}

	// Save master response mode to campaign
	var context map[string]interface{}
	jsonManager.LoadInto(name, ".context.json", &context)
	if context == nil {
		context = make(map[string]interface{})
	}

	modeMapping := map[string]string{
		"Automático (IA responde sempre)":            models.MasterResponseAuto,
		"Semi-automático (IA sugere, você aprova)":  models.MasterResponseSemi,
		"Manual (vê sugestão e edita)":              models.MasterResponseManual,
	}

	context["modo_resposta_mestre"] = modeMapping[masterMode]
	jsonManager.Save(name, ".context.json", context)
	
	// Generate initial context
	initialContext, err := chapterManager.EnsureOpening(name)
	if err != nil {
		fmt.Printf("\033[31mErro ao gerar contexto inicial: %v\033[0m\n", err)
	} else {
		fmt.Println("\033[35mMestre:\033[0m")
		fmt.Println(initialContext)
	}
	
	createPlayers(name)
	gameLoop(name)
}

func createPlayers(campaignName string) {
	fmt.Println("\n\033[33mCriar jogadores:\033[0m")
	
	for {
		fmt.Print("\nNome do jogador (ou Enter para finalizar): ")
		scanner.Scan()
		name := strings.TrimSpace(scanner.Text())
		
		if name == "" {
			break
		}
		
		// Use questionnaire for class selection
		classCategories := generators.GetClassCategories()
		class, err := questionnaire.SelectFromCategories(fmt.Sprintf("Classe de %s:", name), classCategories)
		if err != nil {
			class = "Aventureiro"
		}

		fmt.Printf("Descrição de %s: ", name)
		scanner.Scan()
		description := scanner.Text()

		// Use questionnaire for race selection
		raceCategories := generators.GetRaceCategories()
		race, err := questionnaire.SelectFromCategories(fmt.Sprintf("Raça de %s:", name), raceCategories)
		if err != nil {
			race = "Humano"
		}
		
		player := models.NewPlayer(name)
		player.Class = class
		player.Description = description
		player.Race = race

		// Generate random weaknesses
		player.Weaknesses = weaknessGenerator.GenerateWeaknesses()

		if _, err := playerManager.Add(campaignName, player, true); err != nil {
			fmt.Printf("\033[31mErro: %v\033[0m\n", err)
		} else {
			fmt.Printf("\033[32m%s adicionado.\033[0m\n", name)
			// Get the created player to show formatted attributes and weaknesses
			players, _ := playerManager.List(campaignName)
			for _, p := range players {
				if p.Name == name {
					fmt.Printf("\033[36m%s\033[0m\n", playerManager.FormatAttributes(p))
					fmt.Printf("\033[35m%s\033[0m\n", playerManager.FormatWeaknesses(p))
					break
				}
			}
		}
	}
}

func loadCampaign() {
	name := chooseCampaign()
	if name == "" {
		return
	}
	
	// Ensure opening
	if _, err := chapterManager.EnsureOpening(name); err != nil {
		fmt.Printf("\033[31mErro: %v\033[0m\n", err)
		return
	}
	
	gameLoop(name)
}

func gameLoop(campaignName string) {
	players, err := playerManager.List(campaignName)
	if err != nil {
		fmt.Printf("\033[31mErro ao carregar jogadores: %v\033[0m\n", err)
		return
	}
	
	if len(players) == 0 {
		fmt.Println("Nenhum jogador encontrado. Criando jogadores...")
		createPlayers(campaignName)
		players, _ = playerManager.List(campaignName)
	}
	
	turnIndex := 0
	
	for {
		displayHUD(campaignName, players)
		
		if turnIndex >= len(players) {
			turnIndex = 0
		}
		
		currentPlayer := players[turnIndex]
		
		fmt.Printf("\n\033[33m[Turno de %s] ação / capitulos / encerrar / exportar / sair: \033[0m", currentPlayer.Name)
		scanner.Scan()
		action := strings.TrimSpace(scanner.Text())
		
		actionLower := strings.ToLower(action)
		
		if actionLower == "sair" || actionLower == "menu" {
			break
		}
		
		if actionLower == "capitulos" || actionLower == "capítulos" {
			listChapters(campaignName)
			continue
		}
		
		if actionLower == "encerrar" {
			if chapter, err := chapterManager.Close(campaignName); err != nil {
				fmt.Printf("\033[31mErro: %v\033[0m\n", err)
			} else if chapter != nil {
				fmt.Printf("\033[33m%s\033[0m\n%s\n", chapter.Title, chapter.Summary)
			} else {
				fmt.Println("Nada para fechar ainda.")
			}
			continue
		}
		
		if actionLower == "exportar" {
			if path, err := jsonManager.ExportCampaign(campaignName, ""); err != nil {
				fmt.Printf("\033[31mErro: %v\033[0m\n", err)
			} else {
				fmt.Printf("\033[32mExportado: %s\033[0m\n", path)
			}
			continue
		}
		
		if actionLower == "!encontro" {
			if _, err := encounterGen.GenerateEncounter(campaignName); err != nil {
				fmt.Printf("\033[31mErro: %v\033[0m\n", err)
			}
			continue
		}
		
		if actionLower == "!npc" {
			if _, err := npcGenerator.GenerateNPC(campaignName); err != nil {
			fmt.Printf("\033[31mErro: %v\033[0m\n", err)
			}
			continue
		}

		if actionLower == "!bestiario" {
			showBestiary(campaignName)
			continue
		}

		if actionLower == "!macros" {
			manageMacros(campaignName, currentPlayer.Name)
			continue
		}

		// Check for macros by key (1-9)
		if len(action) == 1 && action >= "1" && action <= "9" {
			if macro, err := playerManager.GetMacroBy(campaignName, currentPlayer.Name, action); err == nil {
				action = macro.Command
				fmt.Printf("\033[36mExecutando macro: %s - %s\033[0m\n", macro.Name, macro.Description)
			}
		}

		// Check for quick actions (no dice roll)
		isQuickAction := false
		for _, qa := range currentPlayer.QuickActions {
			if strings.EqualFold(action, qa) {
				isQuickAction = true
				break
			}
		}

		// Set dice rolling based on action type
		if isQuickAction {
			dice.SetShouldRoll(false)
		} else {
			dice.SetShouldRoll(true)
		}

		var roll int
		var outcome string
		if isQuickAction {
			roll = 0 // No roll for quick actions
			outcome = fmt.Sprintf("O jogador %s realiza ação rápida: %s", currentPlayer.Name, action)
		} else {
			roll = dice.Roll(20)
			outcome = fmt.Sprintf("O jogador %s tenta: %s (d20: %d)", currentPlayer.Name, action, roll)
		}

		// Get master response based on campaign settings
		var context map[string]interface{}
		jsonManager.LoadInto(campaignName, ".context.json", &context)
		masterMode := models.DefaultMasterResponse
		if mode, ok := context["modo_resposta_mestre"].(string); ok {
			masterMode = mode
		}

		// Get master response
		masterResponse, err := chapterManager.GetMasterResponse(campaignName, action, outcome, masterMode)
		if err == nil && masterResponse != "" {
			fmt.Printf("\033[35mMestre:\033[0m %s\n", masterResponse)

			// For manual mode, allow editing
			if masterMode == models.MasterResponseManual {
				fmt.Print("\033[33mEditar resposta? (s/n): \033[0m")
				scanner.Scan()
				editInput := strings.ToLower(strings.TrimSpace(scanner.Text()))
				if editInput == "s" || editInput == "sim" {
					fmt.Print("\033[36mNova resposta: \033[0m")
					scanner.Scan()
					newResponse := scanner.Text()
					if newResponse != "" {
						masterResponse = newResponse
					}
				}
			}
		}

		// Annotate chapter with master response
		if _, err := chapterManager.Annotate(campaignName, outcome, masterResponse); err != nil {
			fmt.Printf("\033[31mErro ao anotar: %v\033[0m\n", err)
		}

		// Award XP for successful actions (not quick actions)
		if !isQuickAction && roll >= 15 {
			xpReward := 10 + roll // More XP for better rolls
			if _, newLevel, leveledUp, err := playerManager.AddXP(campaignName, currentPlayer.Name, xpReward); err == nil {
				if leveledUp {
					fmt.Printf("\033[32m🎉 %s subiu para o nível %d! 🎉\033[0m\n", currentPlayer.Name, newLevel)
				} else {
					fmt.Printf("\033[36m+%d XP ganho\033[0m\n", xpReward)
				}
			}
		}

		// Reload players
		players, _ = playerManager.List(campaignName)
		turnIndex++
	}
}

func displayHUD(campaignName string, players []*models.Player) {
	var context map[string]interface{}
	jsonManager.LoadInto(campaignName, ".context.json", &context)
	
	objective, _ := context["objetivo_principal"].(string)
	if objective == "" {
		objective = "?"
	}
	
	fmt.Printf("\n\033[36m%s\033[0m\n", strings.Repeat("=", 60))
	fmt.Printf(" CAMPANHA: %s | OBJETIVO: %s\n", strings.ToUpper(campaignName), objective)
	fmt.Printf("\033[36m%s\033[0m\n", strings.Repeat("=", 60))
	
	var playerStrings []string
	for _, p := range players {
		class := p.Class
		if class == "" {
			class = "Avent"
		}
		levelProgress := playerManager.GetLevelProgress(p)
		playerStrings = append(playerStrings, fmt.Sprintf("%s (%s) [HP: %d] %s", p.Name, class, p.Life, levelProgress))
	}
	fmt.Printf("\033[32m | \033[0m%s\n", strings.Join(playerStrings, " | "))
}

func listChapters(campaignName string) {
	chapters, err := chapterManager.List(campaignName)
	if err != nil {
		fmt.Printf("\033[31mErro: %v\033[0m\n", err)
		return
	}
	
	for _, chapter := range chapters {
		fmt.Printf("\n\033[33m%s\033[0m (%s)\n%s\n", chapter.Title, chapter.ClosedAt, chapter.Summary)
	}
	
	current, err := chapterManager.Current(campaignName)
	if err != nil {
		return
	}
	
	fmt.Printf("\nEm andamento: capítulo %d — %d cenas\n", current.Number, len(current.Events))
}

func showBestiary(campaignName string) {
	bestiary, err := codexManager.LoadBestiary(campaignName)
	if err != nil {
		fmt.Printf("\033[31mErro: %v\033[0m\n", err)
		return
	}
	
	fmt.Println("\n\033[35m=== BESTIÁRIO DA CAMPANHA ===\033[0m")
	
	if len(bestiary) == 0 {
		fmt.Println("\033[33mNenhum monstro registrado no Bestiário ainda.\033[0m")
		return
	}
	
	i := 1
	for _, monster := range bestiary {
		fmt.Printf("[%d] %s (ND: %s)\n", i, monster.Name, monster.Challenge)
		i++
	}
}

func manageRules() {
	name := chooseCampaign()
	if name == "" {
		return
	}
	
	var rules map[string]interface{}
	jsonManager.LoadInto(name, ".regras.json", &rules)
	
	fmt.Printf("Regras atuais: %v\n", rules)
	
	fmt.Print("Qual regra alterar/criar? ")
	scanner.Scan()
	key := scanner.Text()
	
	if key != "" {
		fmt.Printf("Novo valor para '%s': ", key)
		scanner.Scan()
		value := scanner.Text()
		
		if rules == nil {
			rules = make(map[string]interface{})
		}
		rules[key] = value
		
		if err := jsonManager.Save(name, ".regras.json", rules); err != nil {
			fmt.Printf("\033[31mErro: %v\033[0m\n", err)
		}
	}
}

func exportCampaign() {
	name := chooseCampaign()
	if name == "" {
		return
	}
	
	fmt.Print("Pasta de destino (Enter = exportacoes/): ")
	scanner.Scan()
	destination := strings.TrimSpace(scanner.Text())
	if destination == "" {
		destination = ""
	}
	
	if path, err := jsonManager.ExportCampaign(name, destination); err != nil {
		fmt.Printf("\033[31mErro: %v\033[0m\n", err)
	} else {
		fmt.Printf("\033[32mExportado: %s\033[0m\n", path)
	}
}

func deleteCampaign() {
	name := chooseCampaign()
	if name == "" {
		return
	}

	fmt.Printf("\033[31mExcluir '%s'? (s/n): \033[0m", name)
	scanner.Scan()
	confirm := strings.ToLower(scanner.Text())

	if confirm == "s" {
		if err := jsonManager.DeleteCampaign(name); err != nil {
			fmt.Printf("\033[31mErro: %v\033[0m\n", err)
		} else {
			fmt.Printf("\033[32mCampanha '%s' excluída.\033[0m\n", name)
		}
	}
}

func manageMacros(campaignName, playerName string) {
	for {
		fmt.Println("\n\033[33m=== GERENCIAR MACROS ===\033[0m")
		fmt.Println("\033[32m1.\033[0m Listar Macros")
		fmt.Println("\033[32m2.\033[0m Adicionar Macro")
		fmt.Println("\033[32m3.\033[0m Remover Macro")
		fmt.Println("\033[32m4.\033[0m Voltar")

		fmt.Print("\033[36mEscolha: \033[0m")
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			listMacros(campaignName, playerName)
		case "2":
			addMacro(campaignName, playerName)
		case "3":
			removeMacro(campaignName, playerName)
		case "4":
			return
		default:
			fmt.Println("\033[31mOpção inválida.\033[0m")
		}
	}
}

func listMacros(campaignName, playerName string) {
	players, err := playerManager.List(campaignName)
	if err != nil {
		fmt.Printf("\033[31mErro: %v\033[0m\n", err)
		return
	}

	for _, player := range players {
		if player.Name == playerName {
			if len(player.Macros) == 0 {
				fmt.Println("\033[33mNenhuma macro cadastrada.\033[0m")
				return
			}

			fmt.Println("\n\033[36mMacros do jogador:\033[0m")
			for _, macro := range player.Macros {
				fmt.Printf("\033[32m[%s]\033[0m %s: %s\n", macro.Key, macro.Name, macro.Description)
				fmt.Printf("    Comando: %s\n", macro.Command)
			}
			return
		}
	}

	fmt.Printf("\033[31mJogador '%s' não encontrado.\033[0m\n", playerName)
}

func addMacro(campaignName, playerName string) {
	fmt.Print("Nome da macro: ")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())

	fmt.Print("Descrição da macro: ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())

	fmt.Print("Comando da macro: ")
	scanner.Scan()
	command := strings.TrimSpace(scanner.Text)

	fmt.Print("Atalho (1-9): ")
	scanner.Scan()
	key := strings.TrimSpace(scanner.Text)

	if key < "1" || key > "9" {
		fmt.Println("\033[31mAtalho deve ser entre 1 e 9.\033[0m")
		return
	}

	macro := models.Macro{
		Name:        name,
		Description: description,
		Command:     command,
		Key:         key,
	}

	if err := playerManager.AddMacro(campaignName, playerName, macro); err != nil {
		fmt.Printf("\033[31mErro ao adicionar macro: %v\033[0m\n", err)
	} else {
		fmt.Printf("\033[32mMacro '%s' adicionada com atalho '%s'.\033[0m\n", name, key)
	}
}

func removeMacro(campaignName, playerName string) {
	fmt.Print("Atalho da macro a remover (1-9): ")
	scanner.Scan()
	key := strings.TrimSpace(scanner.Text)

	if err := playerManager.RemoveMacro(campaignName, playerName, key); err != nil {
		fmt.Printf("\033[31mErro ao remover macro: %v\033[0m\n", err)
	} else {
		fmt.Printf("\033[32mMacro com atalho '%s' removida.\033[0m\n", key)
	}
}
