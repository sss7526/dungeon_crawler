package main

import (
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

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
	content := gloss.JoinVertical(
		gloss.Center,
		m.theme.TitleStyle.Render("YOU DIED\n"),
		m.theme.TitleStyle.Foreground(m.theme.Secondary).Render("Press ENTER to return to the Main Menu"),
	)
	border := m.theme.BorderStyle.Render(content)
	return border
}