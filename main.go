package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/list"
)

const (
	startingHealth = 100.0
	maxHealth 		= 100.0
	minHealth		= 0
	healthRegen		= 0.05
	welcomeDuration = 2 * time.Second
)

type Theme struct {
    // Colors
    Primary   gloss.AdaptiveColor
    Secondary gloss.AdaptiveColor
    Selected  gloss.AdaptiveColor
    HealthLow gloss.Color
    HealthHigh gloss.Color

    // Styles
    TitleStyle       gloss.Style
    WelcomeStyle     gloss.Style
    MenuOptionStyle  gloss.Style
    ToolbarStyle     gloss.Style
    ToolbarSelected  gloss.Style
	BorderStyle		gloss.Style

    // UI components
    ProgressBar progress.Model
}

// newTheme initializes and returns a Theme instance.
func newTheme() Theme {
	primaryColor :=  gloss.AdaptiveColor{Light: "#FF5733", Dark: "#AE81FC"}
    return Theme{
        // Adaptive colors for light and dark modes
        Primary:  primaryColor,
        Secondary: gloss.AdaptiveColor{Light: "#FFD700", Dark: "#FF9700"},
        Selected:  gloss.AdaptiveColor{Light: "#00C9A7", Dark: "#1B998B"},
        HealthLow: gloss.Color("#FF3E41"),
        HealthHigh: gloss.Color("#00FF00"),

        // Styles
        TitleStyle: gloss.NewStyle().
            Align(gloss.Center).
            Foreground(gloss.AdaptiveColor{Light: "#FF5733", Dark: "#AE81FC"}).
            Bold(true).
            Width(50),

        WelcomeStyle: gloss.NewStyle().
            Foreground(gloss.AdaptiveColor{Light: "#00DFA2", Dark: "#3EC5F8"}).
            Align(gloss.Center).
            Width(50).
            Bold(true),

        MenuOptionStyle: gloss.NewStyle().
            PaddingLeft(4).
            Foreground(gloss.AdaptiveColor{Light: "#FFD700", Dark: "#FF9700"}),

        ToolbarStyle: gloss.NewStyle().
            Background(gloss.AdaptiveColor{Light: "#FF5733", Dark: "#AE81FC"}).
            Foreground(gloss.Color("#FFFFFF")).
            Padding(0, 1),

        ToolbarSelected: gloss.NewStyle().
            Background(gloss.AdaptiveColor{Light: "#00C9A7", Dark: "#1B998B"}).
            Underline(true).
            Bold(true),

		BorderStyle: gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			Align(gloss.Center),

        // Progress Bar
        ProgressBar: progress.New(progress.WithGradient("#FF3E41", "#00FF00")),
    }
}

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
)

type model struct {
	theme 				Theme 			// Visual configuration for the TUI
	screens 			map[menuChoice]Screen
	currentScreen		Screen
	health			 	float64	 		// Player's health
	quitting		 	bool		 	// Detect if player wants to quite
	progress			progress.Model 	// Progress bar model for health
	activeMenu			int 			// Currently selected toolbar menu in game UI
	toolbar				[]toolbarItem	// The toolbar items
	inventory 			[]string 		// Example of player inventory
	stats 				map[string]int 	// Example player stats
	width 				int				// Terminal width
	height  			int 			// Terminal height
}

type Screen interface {
	Init() tea.Cmd
	Update(msg tea.Msg, m *model) tea.Cmd
	View(m *model) string
}

func initialModel() *model {
	theme := newTheme()
	m := &model{
		theme:				theme,
		health:				100,
		quitting:			false,
		progress:			theme.ProgressBar,
		inventory: 			[]string{"Potion", "Sword", "Shield"},
		stats:				map[string]int{
								"Strength": 10,
								"Agility":	8,
								"Intellect": 5,
		},
	}
	m.screens = map[menuChoice]Screen{
		menuWelcome: NewWelcomeScreen(),
		menuMain: NewMainMenuScreen(m),
		menuGame: NewGameMenuScreen(),
		menuQuitPrompt: NewQuitPromptScreen(),
		menuGameOver: NewGameOverScreen(),
		menuStats: NewStatsScreen(),
	}
	m.currentScreen = m.screens[menuWelcome]
	m.toolbar = newToolbar(m)
	return m
}

func (m *model) switchScreen(choice menuChoice) tea.Cmd {
	m.currentScreen = m.screens[choice]
	return m.currentScreen.Init()
}

type WelcomeScreen struct {
	message 		string
	animatedMessage string
	animationStep 	int 
}

func NewWelcomeScreen() *WelcomeScreen {
	return &WelcomeScreen{
		message: 			"Welcome to the Dungeon!",
		animatedMessage:	"",
		animationStep:		0,
	}
}

func (s *WelcomeScreen) Init() tea.Cmd {
	return doTick()
}

func (s *WelcomeScreen) Update(msg tea.Msg, m *model) tea.Cmd {
	switch msg := msg.(type) {
	case TickMsg:
		if s.animationStep < len(s.message) {
			s.animationStep++
			s.animatedMessage = s.message[:s.animationStep]
			return doTick()
		}
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			return m.switchScreen(menuMain)
		}
	}
	return nil
}

func (s *WelcomeScreen) View(m *model) string {
	content := gloss.JoinVertical(
		gloss.Center,
		m.theme.WelcomeStyle.Render(s.animatedMessage),
		m.theme.WelcomeStyle.Foreground(m.theme.Secondary).Render("Press ENTER to Continue"),
	)

	border := m.theme.BorderStyle.Render(content)
	return border
}

type MainMenuScreen struct {
	list list.Model	// Bubble tea list for menu rendering
}

func NewMainMenuScreen(m *model) *MainMenuScreen {
	items := mainMenuOptions(m)
	menuList := list.New(items, list.NewDefaultDelegate(), 20, 20)
	menuList.Title = "Main Menu"
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	return &MainMenuScreen{list: menuList}
}

func (s *MainMenuScreen) Init() tea.Cmd {
	return nil
}

func (s *MainMenuScreen) Update(msg tea.Msg, m *model) tea.Cmd {
	var cmd tea.Cmd
	var listCmd tea.Cmd
	s.list, listCmd = s.list.Update(msg)
	cmd = tea.Batch(cmd, listCmd)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.switchScreen(menuQuitPrompt)
		case tea.KeyEnter:
			if sel, ok := s.list.SelectedItem().(item); ok && sel.handler != nil {
				return sel.handler()
			}
		}
	}
	return cmd
}

func (s *MainMenuScreen) View(m *model) string {
	return "\n" + s.list.View()
}

type item struct {
	title 			string
	description		string
	handler			func() tea.Cmd
}

func (i item) FilterValue() string { return i.title }
func (i item) Title() string { return i.title }
func (i item) Description() string { return i.description }

func newItem(title, description string, handler func() tea.Cmd ) item {
	return item{
		title:			title,
		description:	description,
		handler:		handler,
	}
}

func (m *model) handleStartNewGame() tea.Cmd { return m.switchScreen(menuGame) }
func (m *model) handleLoadGame() tea.Cmd { return m.switchScreen(menuGame) }
func (m *model) handleQuit() tea.Cmd { return m.switchScreen(menuQuitPrompt) }

func mainMenuOptions(m *model) []list.Item {
	return []list.Item{
		newItem("Start New Game", "", m.handleStartNewGame),
		newItem("Load Game", "", m.handleLoadGame),
		newItem("Quit", "", m.handleQuit),
	}
}

type GameScreen struct{}

func NewGameMenuScreen() *GameScreen {
	return &GameScreen{}
}

func (s *GameScreen) Init() tea.Cmd {
	return doTick()
}

func (s *GameScreen) Update(msg tea.Msg, m *model) tea.Cmd {
	switch msg := msg.(type) {
	case TickMsg:
		if m.health > minHealth {
			m.health = math.Min(maxHealth, m.health + healthRegen)	// Regen health
		}
		if m.health <= minHealth {
			return m.switchScreen(menuGameOver) 		// game over if health runs out
		}
		return doTick()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m.switchScreen(menuQuitPrompt)
		case tea.KeyLeft:
			m.activeMenu = max(m.activeMenu-1, 0)
		case tea.KeyRight:
			m.activeMenu = min(m.activeMenu+1, len(m.toolbar)-1)
		case tea.KeyEnter:
			selected := m.toolbar[m.activeMenu]
			if selected.handler != nil {
				return selected.handler(m)
			}
			m.switchScreen(selected.menuChoice)
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "h":
				m.health = math.Max(0, m.health-10)
			case "r":
				m.health = math.Min(maxHealth, m.health+10)
			}
		}
	}
	return nil
}

func (s *GameScreen) View(m *model) string {
	var b strings.Builder
	for i, item := range m.toolbar {
		if i == m.activeMenu {
			fmt.Fprint(&b, m.theme.ToolbarSelected.Render(item.label) + " ")
		} else {
			fmt.Fprint(&b, m.theme.ToolbarStyle.Render(item.label) + " ")
		}
	}
	return b.String() + "\n\nHealth:\n" + m.theme.ProgressBar.ViewAs(float64(m.health)/maxHealth)
}

type toolbarItem struct {
	label		string
	menuChoice	menuChoice
	handler		func(m *model) tea.Cmd
}

func newToolbarItem(label string, menuChoice menuChoice, handler func(m *model) tea.Cmd) toolbarItem {
	return toolbarItem{
		label:		label,
		menuChoice:	menuChoice,
		handler:	handler,
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

type QuitPromptScreen struct{}

func NewQuitPromptScreen() *QuitPromptScreen {
	return &QuitPromptScreen{}
}

func (s *QuitPromptScreen) Init() tea.Cmd {
	return nil
}

func (s *QuitPromptScreen) Update(msg tea.Msg, m *model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return tea.Quit
		case tea.KeyEsc:
			return m.switchScreen(menuMain)
		}
	}
	return nil
}

func (s *QuitPromptScreen) View(m *model) string {
	content := gloss.JoinVertical(
		gloss.Center,
		m.theme.TitleStyle.Render("Are you sure you want to quit?"),
		m.theme.TitleStyle.Foreground(m.theme.Secondary).Render("ESC to Cancel"),
		m.theme.TitleStyle.Foreground(m.theme.Secondary).Render("ENTER to Confirm"),
	)
	border := m.theme.BorderStyle.Render(content)
	return border

}

type GameOverScreen struct{}

func NewGameOverScreen() *GameOverScreen {
	return &GameOverScreen{}
}

func (s *GameOverScreen) Init() tea.Cmd {
	return nil
}

func (s *GameOverScreen) Update(msg tea.Msg, m *model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m.switchScreen(menuMain)
		}
	}
	return nil
}

func (s *GameOverScreen) View(m *model) string {
	return m.theme.TitleStyle.Render("YOU DIED") + "\n\nPress ENTER to return to the main menu."
}

type StatsScreen struct{}

func NewStatsScreen() *StatsScreen {
	return &StatsScreen{}
}

func (s *StatsScreen) Init() tea.Cmd {
	return nil
}

func (s *StatsScreen) Update(msg tea.Msg, m *model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m.switchScreen(menuGame)
		}
	}
	return nil
}

func (s *StatsScreen) View(m *model) string {
	var b strings.Builder
	fmt.Fprintln(&b, m.theme.TitleStyle.Render("Player Stats"))
	for stat, value := range m.stats {
		fmt.Fprintf(&b, "%s: %d\n", m.theme.MenuOptionStyle.Render(stat), value)
	}
	return b.String()
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
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