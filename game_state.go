package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
	// "fmt"

	// tea "github.com/charmbracelet/bubbletea"
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

