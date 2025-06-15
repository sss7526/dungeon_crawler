package main

import (
	"fmt"
	"strings"
	"math"

	tea "github.com/charmbracelet/bubbletea"
)

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
			m.health = math.Min(maxHealth, m.health+healthRegen) // Regen health
		}
		if m.health <= minHealth {
			return m.switchScreen(menuGameOver) // game over if health runs out
		}
		return doTick()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			return m.saveGameState()
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
			fmt.Fprint(&b, m.theme.ToolbarSelected.Render(item.label)+" ")
		} else {
			fmt.Fprint(&b, m.theme.ToolbarStyle.Render(item.label)+" ")
		}
	}
	return b.String() + "\n\nHealth:\n" + m.theme.ProgressBar.ViewAs(float64(m.health)/maxHealth)
}