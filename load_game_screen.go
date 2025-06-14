package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)


type LoadGameScreen struct {
	list list.Model
}

func createLoadHandler(m *model, filePath string) func() tea.Cmd {
	return func() tea.Cmd {
		m.loadGameState(filePath)
		return m.switchScreen(menuGame)
	}
}

func NewLoadGameScreen(m *model) *LoadGameScreen {

	// List saved games
	saves, err := listSavedGames()
	if err != nil {
		fmt.Println("Failed to list saved games:", err)
		return nil
	}

	// Create list items for each save game
	items := make([]list.Item, len(saves))
	for i, save := range saves {
		fullPath := fmt.Sprintf("save_%v.json", save.Timestamp.Unix())

		items[i] = newItem(
			fmt.Sprintf("Save from %s", save.Timestamp.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("Health: %.0f", save.Health),
			createLoadHandler(m, fullPath),
		)
	}

	l := list.New(items, list.NewDefaultDelegate(), 20, 20)
	l.Title = "Select Saved Game"
	return &LoadGameScreen{list: l}
}

func (s *LoadGameScreen) Init() tea.Cmd {
	return nil
}

func (s *LoadGameScreen) Update(msg tea.Msg, m *model) tea.Cmd {
	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)

	if KeyMsg, ok := msg.(tea.KeyMsg); ok && KeyMsg.Type == tea.KeyEnter {
		if sel, ok := s.list.SelectedItem().(item); ok && sel.handler != nil {
			return sel.handler()
		}
	}

	if KeyMsg, ok := msg.(tea.KeyMsg); ok && KeyMsg.Type == tea.KeyEsc {
		return m.switchScreen(menuMain)
	}
	return cmd
}

func (s *LoadGameScreen) View(m *model) string {
	return s.list.View()
}