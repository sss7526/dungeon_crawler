package main

import (
	// "fmt"
	// "strings"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/bubbles/table"
	gloss "github.com/charmbracelet/lipgloss"
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

// func (s *StatsScreen) View(m *model) string {
// 	var b strings.Builder
// 	fmt.Fprintln(&b, m.theme.TitleStyle.Render("Player Stats"))
// 	for stat, value := range m.stats {
// 		// fmt.Fprintf(&b, "%s: %d\n", m.theme.MenuOptionStyle.Render(stat), value)
// 		fmt.Fprintf(&b, "%s:", m.theme.MenuOptionStyle.Align(gloss.Left).Render(stat))
// 		fmt.Fprintf(&b, "%d\n", m.theme.MenuOptionStyle.Align(gloss.Right).Render(value))
// 	}
// 	return b.String()
// }

// func (s *StatsScreen) View(m *model) string {
// 	player := newTestPlayer()

// 	// Style def
// 	border := m.theme.BorderStyle
// 	header := m.theme.TitleStyle.
// 		PaddingBottom(1).
// 		Underline(true)
// 	sectionTitle := m.theme.MenuOptionStyle.
// 		Bold(true).
// 		MarginBottom(1)
// 	// attributeStyle := gloss.NewStyle().
// 	// 	Width(20).
// 	// 	PaddingRight(2)
// 	progressBar := m.theme.ProgressBar

// 	name := sectionTitle.Render("Name:") + " " + player.Info.Name
// 	class := sectionTitle.Render("Class:") + " " + player.Info.Class
// 	level := sectionTitle.Render("Level:") + fmt.Sprintf(" %d", player.Info.Level)
// 	coreInfo := gloss.JoinVertical(gloss.Left, name, class, level)

// 	// Attributes section
// 	strength := fmt.Sprintf("Strength:  %s [ %d / %d ]",
// 		progressBar.ViewAs(float64(player.Attributes.Strength.Current)/float64(player.Attributes.Strength.Max)),
// 		player.Attributes.Strength.Current, player.Attributes.Strength.Max)

// 	agility := fmt.Sprintf("Agility:   %s [ %d / %d ]",
// 		progressBar.ViewAs(float64(player.Attributes.Agility.Current)/float64(player.Attributes.Agility.Max)),
// 		player.Attributes.Agility.Current, player.Attributes.Agility.Max)

// 	intellect := fmt.Sprintf("Intellect: %s [ %d / %d ]",
// 		progressBar.ViewAs(float64(player.Attributes.Intelect.Current)/float64(player.Attributes.Intelect.Max)),
// 		player.Attributes.Intelect.Current, player.Attributes.Intelect.Max)

// 	endurance := fmt.Sprintf("Endurance: %s [ %d / %d ]",
// 		progressBar.ViewAs(float64(player.Attributes.Endurance.Current)/float64(player.Attributes.Endurance.Max)),
// 		player.Attributes.Endurance.Current, player.Attributes.Endurance.Max)

// 	attributes := gloss.JoinVertical(gloss.Left,
// 		sectionTitle.Render("Attributes"),
// 		gloss.JoinVertical(gloss.Left, strength, agility, intellect, endurance),
// 	)

// 	// Equipment section
// 	equipment := gloss.JoinVertical(gloss.Left,
// 		sectionTitle.Render("Equipment"),
// 		gloss.NewStyle().Render("- " + player.Equipment.Weapon),
// 		gloss.NewStyle().Render("- " + player.Equipment.Armor),
// 		gloss.JoinVertical(gloss.Left, player.Equipment.Trinkets...),
// 	)

// 	// Combine all sections
// 	content := gloss.JoinVertical(gloss.Left,
// 		header.Render("Player Stats"),
// 		border.Render(coreInfo),
// 		border.Render(attributes),
// 		border.Render(equipment),
// 	)

// 	return gloss.NewStyle().
// 		Width(m.terminalWidth).
// 		Height(m.terminalHeight).
// 		Align(gloss.Center).
// 		Render(content)
// }

func (s *StatsScreen) View(m *model) string {
	player := newTestPlayer()
	content := player.View(m.theme)

	return gloss.NewStyle().
		Width(m.terminalWidth).
		Height(m.terminalHeight).
		Align(gloss.Center).
		Render(content)
}