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
		return tea.Batch(func() tea.Msg { return m.loadGameState(filePath) }, m.switchScreen(menuGame))
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
		fullPath := fmt.Sprintf("%s%d%s", saveFilePrefix, save.Timestamp.Unix(), saveFileExt)

		items[i] = newItem(
			fmt.Sprintf("Save from %s", save.Timestamp.Format(saveTimestampFmt)),
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
	
	switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            if sel, ok := s.list.SelectedItem().(item); ok && sel.handler != nil {
                return tea.Batch(sel.handler(), func() tea.Msg {
                    // Provide user feedback after a selection
                    return "Loading game..."
                })
            }
        case tea.KeyEsc:
            return m.switchScreen(menuMain)
        }
    case string: // Handle feedback messages here
        s.list.Title = msg // Update the list title with the message
    }
    return cmd
}

func (s *LoadGameScreen) View(m *model) string {
	return s.list.View()
}