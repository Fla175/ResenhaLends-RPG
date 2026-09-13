package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"resenhalends/internal/ai"
	"resenhalends/internal/generators"
	"resenhalends/internal/managers"
	"resenhalends/internal/models"
	"resenhalends/internal/utils"
	"strconv"
	"strings"
)

var (
	jsonManager     *managers.JSONManager
	playerManager   *managers.PlayerManager
	chapterManager  *managers.ChapterManager
	codexManager    *managers.CodexManager
	npcGenerator    *generators.NPCGenerator
	encounterGen    *generators.EncounterGenerator
	dice            *utils.Dice
	aiClient        managers.AIClient
	currentCampaign string
)

func init() {
	jsonManager = managers.NewJSONManager()
	dice = utils.NewDice(false) // No console output in GUI
	
	// Try to connect to Python AI backend, fall back to mock
	pythonClient := ai.NewPythonClient("http://localhost:8080/api")
	if connected, _ := pythonClient.CheckConnection(); connected {
		aiClient = pythonClient
	} else {
		aiClient = ai.NewMockClient()
	}
	
	playerManager = managers.NewPlayerManager(jsonManager, dice)
	chapterManager = managers.NewChapterManager(jsonManager, aiClient)
	codexManager = managers.NewCodexManager(jsonManager)
	npcGenerator = generators.NewNPCGenerator(jsonManager)
	encounterGen = generators.NewEncounterGenerator(jsonManager)
}

func main() {
	a := app.New()
	a.SetIcon(nil) // You can set an icon here
	
	w := a.NewWindow("ResenhaLends RPG Master")
	w.SetMainMenu(createMainMenu(a))
	w.SetContent(createMainMenuContent(a))
	w.Resize(fyne.NewSize(1024, 768))
	
	w.ShowAndRun()
}

func createMainMenu(a fyne.App) *fyne.MainMenu {
	newCampaignItem := fyne.NewMenuItem("Nova Campanha", func() {
		showNewCampaignDialog(a)
	})
	
	loadCampaignItem := fyne.NewMenuItem("Carregar Campanha", func() {
		showLoadCampaignDialog(a)
	})
	
	manageRulesItem := fyne.NewMenuItem("Gerenciar Regras", func() {
		showManageRulesDialog(a)
	})
	
	exportCampaignItem := fyne.NewMenuItem("Exportar Campanha", func() {
		showExportCampaignDialog(a)
	})
	
	deleteCampaignItem := fyne.NewMenuItem("Excluir Campanha", func() {
		showDeleteCampaignDialog(a)
	})
	
	fileMenu := fyne.NewMenu("Arquivo", newCampaignItem, loadCampaignItem, manageRulesItem, exportCampaignItem, deleteCampaignItem)
	
	quitItem := fyne.NewMenuItem("Sair", func() {
		a.Quit()
	})
	
	fileMenu = fyne.NewMenu("Arquivo", newCampaignItem, loadCampaignItem, manageRulesItem, exportCampaignItem, deleteCampaignItem, fyne.NewMenuItemSeparator(), quitItem)
	
	return fyne.NewMainMenu(fileMenu)
}

func createMainMenuContent(a fyne.App) *fyne.Container {
	logo := widget.NewLabel(`
    =================================================
 ____                      _             _                  _
|  _ \ ___  ___  ___ _ __ | |__   ____  | |    ___ _ __   __| |____
| |_) / _ \/ __|/ _ \ '_ \| '_ \ | '_ \ | |   / _ \ '_ \ / _' / ___|
|  _ <  __/\__ \  __/ | | | | | || |_| || |__|  __/ | | | (_| \__  \
|_| \_\___||___/\___|_| |_|_| |_||_| |_||____\___|__| |_|\__,_|____/
                                                     
        RPG SELF-HOSTED - MESTRE DE MESA PARA RPG
    =================================================
	`)
	
	welcomeLabel := widget.NewLabel("Bem-vindo ao ResenhaLends RPG Master!")
	welcomeLabel.Alignment = fyne.TextAlignCenter
	
	buttons := container.NewVBox(
		widget.NewButton("Nova Campanha", func() {
			showNewCampaignDialog(a)
		}),
		widget.NewButton("Carregar Campanha", func() {
			showLoadCampaignDialog(a)
		}),
		widget.NewButton("Gerenciar Regras", func() {
			showManageRulesDialog(a)
		}),
		widget.NewButton("Exportar Campanha", func() {
			showExportCampaignDialog(a)
		}),
		widget.NewButton("Excluir Campanha", func() {
			showDeleteCampaignDialog(a)
		}),
	)
	
	content := container.NewVBox(
		logo,
		widget.NewSeparator(),
		welcomeLabel,
		widget.NewSeparator(),
		buttons,
	)
	
	return container.NewPadded(content)
}

func showNewCampaignDialog(a fyne.App) {
	w := a.NewWindow("Nova Campanha")
	
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nome da campanha")
	
	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Descrição da campanha")
	
	objEntry := widget.NewEntry()
	objEntry.SetPlaceHolder("Objetivo principal")
	
	difficultySelect := widget.NewSelect([]string{"Fácil", "Normal", "Hardcore"}, nil)
	difficultySelect.SetSelected("Normal")
	
	sizeSelect := widget.NewSelect([]string{"Pequena", "Média", "Extensa"}, nil)
	sizeSelect.SetSelected("Média")

	masterModeSelect := widget.NewSelect([]string{
		"Automático (IA responde sempre)",
		"Semi-automático (IA sugere, você aprova)",
		"Manual (vê sugestão e edita)",
	}, nil)
	masterModeSelect.SetSelected("Semi-automático (IA sugere, você aprova)")

	form := container.NewVBox(
		widget.NewLabel("Nome:"), nameEntry,
		widget.NewLabel("Descrição:"), descEntry,
		widget.NewLabel("Objetivo:"), objEntry,
		widget.NewLabel("Dificuldade:"), difficultySelect,
		widget.NewLabel("Tamanho:"), sizeSelect,
		widget.NewLabel("Modo do Mestre:"), masterModeSelect,
	)
	
	createButton := widget.NewButton("Criar", func() {
		name := nameEntry.Text
		if name == "" {
			return
		}
		
		description := descEntry.Text
		objective := objEntry.Text
		difficulty := difficultySelect.Selected
		size := "media"
		switch sizeSelect.Selected {
		case "Pequena":
			size = "pequena"
		case "Extensa":
			size = "extensa"
		}
		
		if err := jsonManager.SetupCampaign(name, difficulty, description, objective, size); err != nil {
			showErrorDialog(a, "Erro ao criar campanha", err.Error())
			return
		}

		// Save master response mode
		modeMapping := map[string]string{
			"Automático (IA responde sempre)":            models.MasterResponseAuto,
			"Semi-automático (IA sugere, você aprova)":  models.MasterResponseSemi,
			"Manual (vê sugestão e edita)":              models.MasterResponseManual,
		}
		masterMode := modeMapping[masterModeSelect.Selected]

		var context map[string]interface{}
		jsonManager.LoadInto(name, ".context.json", &context)
		if context == nil {
			context = make(map[string]interface{})
		}
		context["modo_resposta_mestre"] = masterMode
		jsonManager.Save(name, ".context.json", context)
		
		currentCampaign = name
		chapterManager.EnsureOpening(name)
		showGameWindow(a, name)
		w.Close()
	})
	
	cancelButton := widget.NewButton("Cancelar", func() {
		w.Close()
	})
	
	buttons := container.NewHBox(createButton, cancelButton)
	
	content := container.NewVBox(form, widget.NewSeparator(), buttons)
	w.SetContent(container.NewPadded(content))
	w.Resize(fyne.NewSize(400, 500))
	w.Show()
}

func showLoadCampaignDialog(a fyne.App) {
	campaigns := jsonManager.ListCampaigns()
	
	if len(campaigns) == 0 {
		showInfoDialog(a, "Nenhuma Campanha", "Não há campanhas disponíveis.")
		return
	}
	
	w := a.NewWindow("Carregar Campanha")
	
	list := widget.NewList(
		func() int {
			return len(campaigns)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(campaigns[id])
		},
	)
	
	loadButton := widget.NewButton("Carregar", func() {
		selected := list.Selected()
		if selected >= 0 && selected < len(campaigns) {
			currentCampaign = campaigns[selected]
			chapterManager.EnsureOpening(currentCampaign)
			showGameWindow(a, currentCampaign)
			w.Close()
		}
	})
	
	cancelButton := widget.NewButton("Cancelar", func() {
		w.Close()
	})
	
	content := container.NewVBox(list, container.NewHBox(loadButton, cancelButton))
	w.SetContent(container.NewPadded(content))
	w.Resize(fyne.NewSize(400, 300))
	w.Show()
}

func showGameWindow(a fyne.App, campaignName string) {
	w := a.NewWindow(campaignName)

	// Load players
	players, _ := playerManager.List(campaignName)

	// Create HUD
	hudLabel := widget.NewLabel("")
	updateHUD(hudLabel, campaignName, players)

	// Add player button
	addPlayerButton := widget.NewButton("Adicionar Jogador", func() {
		showAddPlayerDialog(a, campaignName)
		// Reload players after adding
		updatedPlayers, _ := playerManager.List(campaignName)
		updateHUD(hudLabel, campaignName, updatedPlayers)
	})
	
	// Action input
	actionEntry := widget.NewEntry()
	actionEntry.SetPlaceHolder("Digite sua ação...")

	// Quick action buttons
	quickActionButtons := container.NewHBox(
		widget.NewButton("Sair", func() {
			actionEntry.SetText("Sair")
		}),
		widget.NewButton("Descansar", func() {
			actionEntry.SetText("Descansar")
		}),
		widget.NewButton("Procurar", func() {
			actionEntry.SetText("Procurar")
		}),
		widget.NewButton("Esconder", func() {
			actionEntry.SetText("Esconder")
		}),
	)

	// Action buttons
	actionButtons := container.NewHBox(
		addPlayerButton,
		widget.NewButton("Executar", func() {
			action := actionEntry.Text
			if action == "" {
				return
			}

			// Check if it's a quick action
			isQuickAction := false
			for _, p := range players {
				for _, qa := range p.QuickActions {
					if strings.EqualFold(action, qa) {
						isQuickAction = true
						break
					}
				}
				if isQuickAction {
					break
				}
			}

			var outcome string
			if isQuickAction {
				outcome = fmt.Sprintf("Ação rápida: %s", action)
			} else {
				roll := dice.Roll(20)
				outcome = fmt.Sprintf("Ação executada (d20: %d)", roll)
			}

			chapterManager.Annotate(campaignName, outcome, "")
			actionEntry.SetText("")

			players, _ = playerManager.List(campaignName)
			updateHUD(hudLabel, campaignName, players)
		}),
		widget.NewButton("Capítulos", func() {
			showChaptersDialog(a, campaignName)
		}),
		widget.NewButton("Encerrar Capítulo", func() {
			chapter, _ := chapterManager.Close(campaignName)
			if chapter != nil {
				showInfoDialog(a, "Capítulo Encerrado", chapter.Summary)
			}
		}),
		widget.NewButton("Gerar Encontro", func() {
			encounterGen.GenerateEncounter(campaignName)
		}),
		widget.NewButton("Gerar NPC", func() {
			npcGenerator.GenerateNPC(campaignName)
		}),
		widget.NewButton("Voltar", func() {
			w.Close()
		}),
	)
	
	content := container.NewVBox(
		hudLabel,
		widget.NewSeparator(),
		actionEntry,
		quickActionButtons,
		actionButtons,
	)
	
	w.SetContent(container.NewPadded(content))
	w.Resize(fyne.NewSize(800, 600))
	w.Show()
}

func updateHUD(label *widget.Label, campaignName string, players []*models.Player) {
	var context map[string]interface{}
	jsonManager.LoadInto(campaignName, ".context.json", &context)
	
	objective, _ := context["objetivo_principal"].(string)
	if objective == "" {
		objective = "?"
	}
	
	hudText := "CAMPANHA: " + campaignName + " | OBJETIVO: " + objective + "\n\n"
	
	for _, p := range players {
		class := p.Class
		if class == "" {
			class = "Avent"
		}
		levelProgress := playerManager.GetLevelProgress(p)
		hudText += p.Name + " (" + class + ") [HP: " + fmt.Sprintf("%d", p.Life) + "] " + levelProgress + "\n"
	}
	
	label.SetText(hudText)
}

func showChaptersDialog(a fyne.App, campaignName string) {
	w := a.NewWindow("Capítulos")
	
	chapters, _ := chapterManager.List(campaignName)
	
	var chapterText string
	for _, chapter := range chapters {
		chapterText += chapter.Title + "\n"
		chapterText += chapter.Summary + "\n\n"
	}
	
	text := widget.NewLabel(chapterText)
	text.Wrapping = fyne.TextWrapWord
	
	scroll := container.NewScroll(text)
	
	closeButton := widget.NewButton("Fechar", func() {
		w.Close()
	})
	
	content := container.NewVBox(scroll, closeButton)
	w.SetContent(container.NewPadded(content))
	w.Resize(fyne.NewSize(600, 400))
	w.Show()
}

func showManageRulesDialog(a fyne.App) {
	w := a.NewWindow("Gerenciar Regras")
	
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nome da campanha")
	
	ruleEntry := widget.NewEntry()
	ruleEntry.SetPlaceHolder("Regra")
	
	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder("Valor")
	
	form := container.NewVBox(
		widget.NewLabel("Campanha:"), nameEntry,
		widget.NewLabel("Regra:"), ruleEntry,
		widget.NewLabel("Valor:"), valueEntry,
	)
	
	saveButton := widget.NewButton("Salvar", func() {
		campaignName := nameEntry.Text
		if campaignName == "" {
			return
		}
		
		var rules map[string]interface{}
		jsonManager.LoadInto(campaignName, ".regras.json", &rules)
		
		if rules == nil {
			rules = make(map[string]interface{})
		}
		
		rules[ruleEntry.Text] = valueEntry.Text
		jsonManager.Save(campaignName, ".regras.json", rules)
		
		w.Close()
	})
	
	cancelButton := widget.NewButton("Cancelar", func() {
		w.Close()
	})
	
	buttons := container.NewHBox(saveButton, cancelButton)
	
	content := container.NewVBox(form, widget.NewSeparator(), buttons)
	w.SetContent(container.NewPadded(content))
	w.Resize(fyne.NewSize(400, 300))
	w.Show()
}

func showExportCampaignDialog(a fyne.App) {
	w := a.NewWindow("Exportar Campanha")
	
	campaigns := jsonManager.ListCampaigns()
	
	if len(campaigns) == 0 {
		showInfoDialog(a, "Nenhuma Campanha", "Não há campanhas disponíveis.")
		w.Close()
		return
	}
	
	list := widget.NewList(
		func() int {
			return len(campaigns)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(campaigns[id])
		},
	)
	
	exportButton := widget.NewButton("Exportar", func() {
		selected := list.Selected()
		if selected >= 0 && selected < len(campaigns) {
			path, err := jsonManager.ExportCampaign(campaigns[selected], "")
			if err != nil {
				showErrorDialog(a, "Erro", err.Error())
			} else {
				showInfoDialog(a, "Exportado", "Campanha exportada para: "+path)
			}
			w.Close()
		}
	})
	
	cancelButton := widget.NewButton("Cancelar", func() {
		w.Close()
	})
	
	content := container.NewVBox(list, container.NewHBox(exportButton, cancelButton))
	w.SetContent(container.NewPadded(content))
	w.Resize(fyne.NewSize(400, 300))
	w.Show()
}

func showDeleteCampaignDialog(a fyne.App) {
	w := a.NewWindow("Excluir Campanha")
	
	campaigns := jsonManager.ListCampaigns()
	
	if len(campaigns) == 0 {
		showInfoDialog(a, "Nenhuma Campanha", "Não há campanhas disponíveis.")
		w.Close()
		return
	}
	
	list := widget.NewList(
		func() int {
			return len(campaigns)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(campaigns[id])
		},
	)
	
	deleteButton := widget.NewButton("Excluir", func() {
		selected := list.Selected()
		if selected >= 0 && selected < len(campaigns) {
			err := jsonManager.DeleteCampaign(campaigns[selected])
			if err != nil {
				showErrorDialog(a, "Erro", err.Error())
			} else {
				showInfoDialog(a, "Excluído", "Campanha excluída com sucesso.")
			}
			w.Close()
		}
	})
	
	cancelButton := widget.NewButton("Cancelar", func() {
		w.Close()
	})
	
	content := container.NewVBox(list, container.NewHBox(deleteButton, cancelButton))
	w.SetContent(container.NewPadded(content))
	w.Resize(fyne.NewSize(400, 300))
	w.Show()
}

func showErrorDialog(a fyne.App, title, message string) {
	w := a.NewWindow(title)
	label := widget.NewLabel(message)
	label.Wrapping = fyne.TextWrapWord
	closeButton := widget.NewButton("Fechar", func() {
		w.Close()
	})
	content := container.NewVBox(label, closeButton)
	w.SetContent(container.NewPadded(content))
	w.Show()
}

func showInfoDialog(a fyne.App, title, message string) {
	w := a.NewWindow(title)
	label := widget.NewLabel(message)
	label.Wrapping = fyne.TextWrapWord
	closeButton := widget.NewButton("Fechar", func() {
		w.Close()
	})
	content := container.NewVBox(label, closeButton)
	w.SetContent(container.NewPadded(content))
	w.Show()
}

func showAddPlayerDialog(a fyne.App, campaignName string) {
	w := a.NewWindow("Adicionar Jogador")

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nome do jogador")

	classEntry := widget.NewEntry()
	classEntry.SetPlaceHolder("Classe")

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetPlaceHolder("Descrição")

	raceEntry := widget.NewEntry()
	raceEntry.SetPlaceHolder("Raça")
	raceEntry.SetText("Humano")

	lifeEntry := widget.NewEntry()
	lifeEntry.SetPlaceHolder("Vida (HP)")
	lifeEntry.SetText("20")

	form := container.NewVBox(
		widget.NewLabel("Nome:"), nameEntry,
		widget.NewLabel("Classe:"), classEntry,
		widget.NewLabel("Descrição:"), descEntry,
		widget.NewLabel("Raça:"), raceEntry,
		widget.NewLabel("Vida:"), lifeEntry,
	)

	createButton := widget.NewButton("Criar", func() {
		name := nameEntry.Text
		if name == "" {
			showErrorDialog(a, "Erro", "Nome é obrigatório")
			return
		}

		player := models.NewPlayer(name)
		player.Class = classEntry.Text
		player.Description = descEntry.Text
		player.Race = raceEntry.Text
		if raceEntry.Text == "" {
			player.Race = "Humano"
		}

		// Parse life
		var life int
		if lifeEntry.Text != "" {
			if l, err := strconv.Atoi(lifeEntry.Text); err == nil {
				life = l
			}
		}
		if life > 0 {
			player.Life = life
		}

		if _, err := playerManager.Add(campaignName, player, true); err != nil {
			showErrorDialog(a, "Erro", err.Error())
			return
		}

		// Show formatted attributes
		players, _ := playerManager.List(campaignName)
		for _, p := range players {
			if p.Name == name {
				showInfoDialog(a, "Jogador Criado", playerManager.FormatAttributes(p))
				break
			}
		}

		w.Close()
	})

	cancelButton := widget.NewButton("Cancelar", func() {
		w.Close()
	})

	buttons := container.NewHBox(createButton, cancelButton)

	content := container.NewVBox(form, widget.NewSeparator(), buttons)
	w.SetContent(container.NewPadded(content))
	w.Resize(fyne.NewSize(400, 500))
	w.Show()
}
