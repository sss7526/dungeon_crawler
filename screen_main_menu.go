package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)

type MainMenuScreen struct {
	list list.Model // Bubble tea list for menu rendering
}

func NewMainMenuScreen(m *model) *MainMenuScreen {
	items := mainMenuOptions(m)
	menuList := list.New(items, list.NewDefaultDelegate(), 20, 20)
	menuList.Title = "Main Menu"
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	return &MainMenuScreen{list: menuList}
}

func (s *MainMenuScreen) Init() tea.Cmd {
	return nil
}

func (s *MainMenuScreen) Update(msg tea.Msg, m *model) tea.Cmd {
	var cmd tea.Cmd
	var listCmd tea.Cmd
	s.list, listCmd = s.list.Update(msg)
	cmd = tea.Batch(cmd, listCmd)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.switchScreen(menuQuitPrompt)
		case tea.KeyEnter:
			if sel, ok := s.list.SelectedItem().(item); ok && sel.handler != nil {
				return sel.handler()
			}
		}
	}
	return cmd
}

func (s *MainMenuScreen) View(m *model) string {
	return "\n" + s.list.View()
}

func (m *model) handleStartNewGame() tea.Cmd { return m.switchScreen(menuGame) }
func (m *model) handleLoadGame() tea.Cmd     { return m.switchScreen(menuLoadGameScreen) }
func (m *model) handleQuit() tea.Cmd         { return m.switchScreen(menuQuitPrompt) }

func mainMenuOptions(m *model) []list.Item {
	return []list.Item{
		newItem("Start New Game", "", m.handleStartNewGame),
		newItem("Load Game", "", m.handleLoadGame),
		newItem("Quit", "", m.handleQuit),
	}
}