package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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