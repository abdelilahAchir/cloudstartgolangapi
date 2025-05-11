package main

import (
	"systementor.se/cloudgolangapi/data"
)

// Setup mock functions for data package
func setupMocks() func() {
	// Save original functions to restore later
	originalInitDB := data.InitDatabase
	originalSaveGame := data.SaveGame
	originalStats := data.Stats

	// Replace with mock implementations
	data.InitDatabase = func(file, server, database, username, password string, port int) {
		// Do nothing - mock implementation
	}

	data.SaveGame = func(yourSelection, mySelection, winner string) {
		// Do nothing - mock implementation
	}

	data.Stats = func() (int, int) {
		return 10, 5 // Mock stats: 10 total games, 5 wins
	}

	// Return cleanup function
	return func() {
		data.InitDatabase = originalInitDB
		data.SaveGame = originalSaveGame
		data.Stats = originalStats
	}
}
