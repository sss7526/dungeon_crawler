package main

import (
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

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
		m.theme.TitleStyle.Render("Are you sure you want to quit?\n"),
		m.theme.TitleStyle.Foreground(m.theme.Secondary).Render("ESC to Cancel"),
		m.theme.TitleStyle.Foreground(m.theme.Secondary).Render("ENTER to Confirm"),
	)
	border := m.theme.BorderStyle.Render(content)
	return border

}