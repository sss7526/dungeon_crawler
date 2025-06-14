package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type GameState struct {
	Health 		float64				`json:"health"`
	Inventory 	[]string 			`json:"inventory"`
	Stats 		map[string]int 		`json:"stats"`
	Timestamp	time.Time			`json:"timestamp"`
}

// getSaveDir returns the directory for saving game files, creating if necessary
func getSaveDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	saveDir := filepath.Join(dir, "dungeon_crawler") // Create a subdirectory for the game
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return "", err
	}
	return saveDir, nil
}

// listSavedGames returns a list of saved game files
func listSavedGames() ([]GameState, error) {
	saveDir, err := getSaveDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(saveDir)
	if err != nil {
		return nil, err
	}

	var saves []GameState
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		file, err := os.Open(filepath.Join(saveDir, entry.Name()))
		if err != nil {
			continue
		}
		defer file.Close()

		var gameState GameState
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&gameState); err != nil {
			continue
		}
		saves = append(saves, gameState)
	}

	return saves, nil
}

// saveGameState saves the current state of the game to a file.
func (m *model) saveGameState() tea.Cmd {
	saveDir, err := getSaveDir()
	if err != nil {
		return func() tea.Msg {
			return fmt.Sprintf("Failed to get save directory: %v", err)
		}
	}

	fileName := fmt.Sprintf("save_%d.json", time.Now().Unix()) // Unique filename with timestamp
	savePath := filepath.Join(saveDir, fileName)

	gameState := GameState{
		Health: 	m.health,
		Inventory: 	m.inventory,
		Stats:		m.stats,
		Timestamp:	time.Now(),
	}

	file, err := os.Create(savePath)
	if err != nil {
		return func() tea.Msg {
			return fmt.Sprintf("Failed to save game: %v", err)
		}
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&gameState); err != nil {
		return func() tea.Msg {
			return fmt.Sprintf("Failed to encode game state: %v", err)
		}
	}

	return func() tea.Msg {
		return "Game saved successfully"
	}
}

// loadGameState loads a selected game state from a file.
func (m *model) loadGameState(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		// return func() tea.Msg {
		// 	return fmt.Sprintf("Failed to load game: %v", err)
		// }
		return
	}
	defer file.Close()

	var gameState GameState
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&gameState); err != nil {
		// return func() tea.Msg {
			// return fmt.Sprintf("Failed to decode save file: %v", err)
		// }
		return
	}

	// Apply the loaded state to the model
	m.health = gameState.Health
	m.inventory = gameState.Inventory
	m.stats = gameState.Stats

	// return func() tea.Msg {
	// 	return "Game loaded successfully!"
	// }
}

