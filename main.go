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

    // UI components
    ProgressBar progress.Model
}

// newTheme initializes and returns a Theme instance.
func newTheme() Theme {
    return Theme{
        // Adaptive colors for light and dark modes
        Primary:   gloss.AdaptiveColor{Light: "#FF5733", Dark: "#AE81FC"},
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

        // Progress Bar
        ProgressBar: progress.New(progress.WithGradient("#FF3E41", "#00FF00")),
    }
}
// var (

// 	gradientPrimary = gloss.AdaptiveColor{Light: "#FF5733", Dark: "#AE81FC"}
// 	gradientAlt = gloss.AdaptiveColor{Light: "#FFD700", Dark: "#FF9700"}
// 	gradientSelected = gloss.AdaptiveColor{Light: "#00C9A7", Dark: "#1B998B"}


// 	styleTitle = gloss.NewStyle().
// 		Align(gloss.Center).
// 		Foreground(gradientPrimary).
// 		Bold(true).
// 		Width(50)

// 	styleWelcomeMessage = gloss.NewStyle().
// 		Foreground(gloss.AdaptiveColor{Light: "#00DFA2", Dark: "#3EC5F8"}).
// 		Align(gloss.Center).
// 		Width(50).
// 		Bold(true)

// 	styleWelcomeHelpPrompt = styleWelcomeMessage.Foreground(gradientAlt)
	
// 	styleMenuOption = gloss.NewStyle().
// 		PaddingLeft(4).
// 		Foreground(gradientAlt)

// 	// Regular toolbar item style
// 	toolbarStyle = gloss.NewStyle().
// 			Background(gradientPrimary).
// 			Foreground(gloss.Color("#FFFFFF")).
// 			Padding(0, 1) // Add some horizontal padding

// 	// Highlighted (active/selected) toolbar item style
// 	toolbarSelected = toolbarStyle.
// 			Background(gradientSelected).
// 			Underline(true).
// 			Bold(true)

// 	// styleSelectedItem = gloss.NewStyle().
// 	// 	Foreground(gradientPrimary).
// 	// 	Bold(true)

// 	progressBar = progress.New(progress.WithGradient("#FF3E41", "#00FF00"))
// )

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
	currentMenu		 	menuChoice  	// Current screen being displayed
	welcomeMessage	 	string		 	// Welcome text that fades in
	animatedMessage  	string 
	animationStep		int
	health			 	float64	 		// Player's health
	quitting		 	bool		 	// Detect if player wants to quite
	progress			progress.Model 	// Progress bar model for health
	list				*list.Model  	// Main menu option model
	activeMenu			int 			// Currently selected toolbar menu in game UI
	toolbar				[]toolbarItem	// The toolbar items
	inventory 			[]string 		// Example of player inventory
	stats 				map[string]int 	// Example player stats
}

func initialModel() *model {
	theme := newTheme()
	m := &model{
		theme:				theme,
		currentMenu:		menuWelcome,
		welcomeMessage: 	"Welcome to the Dungeon!",
		animatedMessage: 	"",
		animationStep: 		0,
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
	m.list = mainMenu(m)
	m.toolbar = newToolbar(m)
	return m
}

func mainMenu(m *model) *list.Model {
	options := mainMenuOptions(m)

	menuList := list.New(options, list.NewDefaultDelegate(), 20, 20)
	menuList.Title = "Main Menu"
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	return &menuList
}

type item struct {
	title 			string
	description		string
	handler			func()
}

func (i item) FilterValue() string { return i.title }
func (i item) Title() string { return i.title }
func (i item) Description() string { return i.description }

func newItem(title, description string, handler func()) item {
	return item{
		title:			title,
		description:	description,
		handler:		handler,
	}
}

func (m *model) handleStartNewGame() { m.currentMenu = menuGame }
func (m *model) handleLoadGame() { m.currentMenu = menuGame }
func (m *model) handleQuit() { m.currentMenu = menuQuitPrompt }

func mainMenuOptions(m *model) []list.Item {
	return []list.Item{
		newItem("Start New Game", "", m.handleStartNewGame),
		newItem("Load Game", "", m.handleLoadGame),
		newItem("Quit", "", m.handleQuit),
	}
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

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})

}

func (m *model) Init() tea.Cmd {
	return doTick()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	if m.health <= 0 && m.currentMenu == menuGame {
		m.currentMenu = menuGameOver
		return m, nil
	}

	switch msg := msg.(type) {
	case TickMsg:
		switch m.currentMenu {
		case menuWelcome:
			if m.animationStep < len(m.welcomeMessage) {
				m.animationStep++
				m.animatedMessage = m.welcomeMessage[:m.animationStep]
			}
		case menuGame:
			m.health = math.Min(maxHealth, m.health + 0.05)
		}
		return m, doTick()
	
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch m.currentMenu {
		case menuWelcome:
			if msg.Type == tea.KeyEnter {
				m.currentMenu = menuMain
			}
		case menuMain:
			var listCmd tea.Cmd
			*m.list, listCmd = m.list.Update(msg)
			cmd = tea.Batch(cmd, listCmd)

			switch msg.Type {
			case tea.KeyEsc:
				m.currentMenu = menuQuitPrompt
			case tea.KeyEnter:
				if sel, ok := m.list.SelectedItem().(item); ok && sel.handler != nil {
					sel.handler()
				}
			}
		case menuGame:
			switch msg.Type {
			case tea.KeyEsc:
				m.currentMenu = menuQuitPrompt
			case tea.KeyLeft:
				m.activeMenu = max(m.activeMenu-1, 0)
			case tea.KeyRight:
				m.activeMenu = min(m.activeMenu+1, len(m.toolbar)-1)
			case tea.KeyEnter:
				selected := m.toolbar[m.activeMenu]
				if selected.handler != nil {
					return m, selected.handler(m)
				}
				m.currentMenu = selected.menuChoice
			case tea.KeyRunes:	// h and r keypresses are for testing
				switch string(msg.Runes) {
				case "h":
					m.health = math.Max(0, m.health - 10) // Reduce health
				case "r":
					m.health = math.Min(maxHealth, m.health + 10) // Restore health
				}
			}
		case menuStats:
			switch msg.Type {
			case tea.KeyEsc:
				m.currentMenu = menuGame
			}
		case menuQuitPrompt:
			switch msg.Type {
			case tea.KeyEsc:
				m.currentMenu = menuMain
			case tea.KeyEnter:
				return m, tea.Quit
			}
		case menuGameOver:
			switch msg.Type {
			case tea.KeyEnter:
				m.currentMenu = menuMain
			}
		}
	}
	return m, cmd
}

func (m *model) View() string {
	switch m.currentMenu {
	case menuWelcome:
		return renderWelcomeScreen(m)
	case menuMain:
		return renderMainMenu(m)
	case menuGame:
		return renderGameScreen(m)
	case menuStats:
		return renderStats(m)
	case menuQuitPrompt:
		return renderQuitPrompt(m)
	case menuGameOver:
		return renderGameOver(m)
	// case menuFile:
	default:
		return "Unknown menu"
	}
}

func renderWelcomeScreen(m *model) string {
	return m.theme.WelcomeStyle.Render(m.animatedMessage) +
		"\n\n" +
		m.theme.WelcomeStyle.Foreground(m.theme.Secondary).Render("Press ENTER to Continue")
}

func renderQuitPrompt(m *model) string {
	return m.theme.TitleStyle.Render("Are you sure you want to quit? (ESC to cancel, ENTER to confirm)")
}

func renderGameOver(m *model) string {
	return m.theme.TitleStyle.Render("YOU DIED") + "\n\nPress ENTER to return to the main menu."
}

func renderMainMenu(m *model) string {
	return "\n" + m.list.View()
}

func renderGameScreen(m *model) string {
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

func renderStats(m *model) string {
	var b strings.Builder
	fmt.Fprintln(&b, m.theme.TitleStyle.Render("Player Stats"))
	for stat, value := range m.stats {
		fmt.Fprintf(&b, "%s: %d\n", m.theme.MenuOptionStyle.Render(stat), value)
	}
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting application:", err)
		os.Exit(1)
	}
}