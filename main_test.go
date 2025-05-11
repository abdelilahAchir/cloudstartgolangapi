package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"systementor.se/cloudgolangapi/data"
)

// Setup test router with the routes we want to test
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/", start)
	router.GET("/api/play", apiPlay)
	router.GET("/api/stats", apiStats)
	return router
}

// Test the start endpoint
func TestStartEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Tjena", w.Body.String())
}

// Test the stats endpoint
func TestStatsEndpoint(t *testing.T) {
	// Mock data.Stats before setting up the router
	originalStats := data.Stats
	data.Stats = func() (int, int) {
		return 42, 21 // mock values
	}
	defer func() { data.Stats = originalStats }()

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/stats", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify the response has the expected keys
	_, hasTotalGames := response["totalGames"]
	_, hasWins := response["wins"]
	assert.True(t, hasTotalGames)
	assert.True(t, hasWins)
}

// Test the play endpoint
func TestPlayEndpoint(t *testing.T) {
	// Initialize random with fixed seed for deterministic tests
	theRandom = rand.New(rand.NewSource(42))

	// Mock SaveGame function
	originalSaveGame := data.SaveGame
	data.SaveGame = func(yourSelection, mySelection, winner string) {
		// Do nothing - mock implementation
	}
	defer func() { data.SaveGame = originalSaveGame }()

	router := setupRouter()

	// Test case 1: User plays STONE
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/play?yourSelection=STONE", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response has the expected keys
	_, hasWinner := response["winner"]
	_, hasYourSelection := response["yourSelection"]
	_, hasComputerSelection := response["computerSelection"]
	assert.True(t, hasWinner)
	assert.True(t, hasYourSelection)
	assert.True(t, hasComputerSelection)
	assert.Equal(t, "STONE", response["yourSelection"])
}

// Test the randomizeSelection function
func TestRandomizeSelection(t *testing.T) {
	// Use a fixed seed for deterministic tests
	theRandom = rand.New(rand.NewSource(42))

	// Get several selections to make sure we get different values
	selection1 := randomizeSelection()
	selection2 := randomizeSelection()
	selection3 := randomizeSelection()

	// Verify they are valid selections
	validSelections := map[string]bool{
		"STONE":   true,
		"SCISSOR": true,
		"BAG":     true,
	}

	assert.True(t, validSelections[selection1])
	assert.True(t, validSelections[selection2])
	assert.True(t, validSelections[selection3])
}
