package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/progress"
)

const (
	startingHealth  = 100.0
	maxHealth       = 100.0
	minHealth       = 0
	healthRegen     = 0.05
	welcomeDuration = 2 * time.Second
)

type menuChoice int

const (
	menuWelcome menuChoice = iota
	menuMain
	menuGame
	menuFile
	menuStats
	menuInventory
	menuHelp
	menuQuitPrompt
	menuGameOver
	menuLoadGameScreen
)

type model struct {
	theme         Theme // Visual configuration for the TUI
	screens       map[menuChoice]Screen
	currentScreen Screen
	health        float64        // Player's health
	quitting      bool           // Detect if player wants to quite
	progress      progress.Model // Progress bar model for health
	activeMenu    int            // Currently selected toolbar menu in game UI
	toolbar       []toolbarItem  // The toolbar items
	inventory     []string       // Example of player inventory
	stats         map[string]int // Example player stats
}

type Screen interface {
	Init() tea.Cmd
	Update(msg tea.Msg, m *model) tea.Cmd
	View(m *model) string
}

func initialModel() *model {
	theme := newTheme()
	m := &model{
		theme:     theme,
		health:    100,
		quitting:  false,
		progress:  theme.ProgressBar,
		inventory: []string{"Potion", "Sword", "Shield"},
		stats: map[string]int{
			"Strength":  10,
			"Agility":   8,
			"Intellect": 5,
		},
	}
	m.screens = map[menuChoice]Screen{
		menuWelcome:    	NewWelcomeScreen(),
		menuMain:       	NewMainMenuScreen(m),
		menuGame:       	NewGameMenuScreen(),
		menuQuitPrompt: 	NewQuitPromptScreen(),
		menuGameOver:   	NewGameOverScreen(),
		menuStats:      	NewStatsScreen(),
		menuLoadGameScreen: NewLoadGameScreen(m),
	}
	m.currentScreen = m.screens[menuWelcome]
	m.toolbar = newToolbar(m)
	return m
}

func (m *model) switchScreen(choice menuChoice) tea.Cmd {
	m.currentScreen = m.screens[choice]
	return m.currentScreen.Init()
}

type item struct {
	title       string
	description string
	handler     func() tea.Cmd
}

func (i item) FilterValue() string { return i.title }
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }

func newItem(title, description string, handler func() tea.Cmd) item {
	return item{
		title:       title,
		description: description,
		handler:     handler,
	}
}

type toolbarItem struct {
	label      string
	menuChoice menuChoice
	handler    func(m *model) tea.Cmd
}

func newToolbarItem(label string, menuChoice menuChoice, handler func(m *model) tea.Cmd) toolbarItem {
	return toolbarItem{
		label:      label,
		menuChoice: menuChoice,
		handler:    handler,
	}
}

func newToolbar(_ *model) []toolbarItem {
	return []toolbarItem{
		newToolbarItem("File", menuFile, nil),
		newToolbarItem("Stats", menuStats, nil),
		newToolbarItem("Inventory", menuInventory, nil),
		newToolbarItem("Help", menuHelp, nil),
	}
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})

}

func (m *model) Init() tea.Cmd {
	// Start the game clock
	return doTick()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Delegate updates to the current screen's Update method
	cmd := m.currentScreen.Update(msg, m)
	return m, cmd
}

func (m *model) View() string {
	// Delegate rendering to the current screen's View method
	return m.currentScreen.View(m)
}

func main() {
	for {
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Println("Error starting application:", err)
			continue
		}

		os.Exit(0)
	}

}
