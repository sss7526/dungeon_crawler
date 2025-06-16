package main

import (
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type ErrorScreen struct{}

func NewErrorScreen() *ErrorScreen {
	return &ErrorScreen{}
}

func (s *ErrorScreen) Init() tea.Cmd {
	return nil
}

func (s *ErrorScreen) Update(msg tea.Msg, m *model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m.switchScreen(menuMain)
		case tea.KeyEsc:
			return m.switchScreen(menuQuitPrompt)
		}
	}
	return nil
}

func (s *ErrorScreen) View(m *model) string {
	content := gloss.JoinVertical(
		gloss.Center,
		m.theme.ErrorStyle.Render("An error ocurred."),
		m.theme.ErrorStyle.Render("ESC to Quit"),
		m.theme.ErrorStyle.Render("Press ENTER for Main Menu"),
	)
	border := m.theme.ErrorBorder.Render(content)
	return border
}