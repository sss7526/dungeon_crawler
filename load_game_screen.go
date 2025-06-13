package main

import (
	"path/filepath"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)


type LoadGameScreen struct {
	list list.Model
}

func NewLoadGameScreen() *LoadGameScreen {
	// Get the save directory
	saveDir, err := getSaveDir()
	if err != nil {
		fmt.Println("Failed to get save dir:", err)
		return nil
	}

	// List saved games
	saves, _ := listSavedGames() // Ignore errors for now
	items := make([]list.Item, len(saves))
	for i, save := range saves {
		fileName := fmt.Sprintf("save_%v.json", save.Timestamp.Unix())
		fullPath := filepath.Join(saveDir, fileName)

		items[i] = newItem(
			fmt.Sprintf("Save from %s", save.Timestamp.Format("2006-01-02 15:04:05")),
			fmt.Sprintf("Health: %.0f", save.Health),
			func(filePath string) func() tea.Cmd {
				return func() tea.Cmd {
					return initialModel().loadGameState(filePath)
				}
			}(fullPath), // filepath handler
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
	return cmd
}

func (s *LoadGameScreen) View(m *model) string {
	return s.list.View()
}