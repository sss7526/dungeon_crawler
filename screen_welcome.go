package main

import (
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type WelcomeScreen struct {
	message         string
	animatedMessage string
	animationStep   int
}

func NewWelcomeScreen() *WelcomeScreen {
	return &WelcomeScreen{
		message:         "Welcome to the Dungeon!",
		animatedMessage: "",
		animationStep:   0,
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
		m.theme.WelcomeStyle.Foreground(m.theme.Secondary).Render("\nPress ENTER to Continue"),
	)

	border := m.theme.BorderStyle.Render(content)
	return border
}