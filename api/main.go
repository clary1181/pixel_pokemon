package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Pokemon represents a Pokemon entity
type Pokemon struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Types    []string `json:"types"`
	Level    int      `json:"level"`
	HP       int      `json:"hp"`
	MaxHP    int      `json:"maxHp"`
	Attack   int      `json:"attack"`
	Defense  int      `json:"defense"`
	Speed    int      `json:"speed"`
	Moves    []Move   `json:"moves"`
	ImageURL string   `json:"imageUrl"`
}

// Move represents a Pokemon move
type Move struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Power    int    `json:"power"`
	Accuracy int    `json:"accuracy"`
}

// BattleResult represents the outcome of a battle simulation
type BattleResult struct {
	Winner       string   `json:"winner"`
	Rounds       int      `json:"rounds"`
	BattleLog    []string `json:"battleLog"`
	PlayerHP     int      `json:"playerHp"`
	OpponentHP   int      `json:"opponentHp"`
	PlayerMoves  []string `json:"playerMoves"`
	OpponentMoves []string `json:"opponentMoves"`
}

// SaveRequest is the payload for saving game progress
type SaveRequest struct {
	PlayerID  string    `json:"playerId"`
	Pokemons  []Pokemon `json:"pokemons"`
	Position  Position  `json:"position"`
	Timestamp int64     `json:"timestamp"`
}

// Position represents player position in the game
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// SaveResponse is the response for save requests
type SaveResponse struct {
	Success   bool   `json:"success"`
	SaveID    string `json:"saveId,omitempty"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message,omitempty"`
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Global variables
var (
	pokedex = []Pokemon{
		{
			ID:       1,
			Name:     "Bulbasaur",
			Types:    []string{"grass", "poison"},
			Level:    5,
			HP:       45,
			MaxHP:    45,
			Attack:   49,
			Defense:  49,
			Speed:    45,
			ImageURL: "/assets/sprites/bulbasaur.png",
			Moves: []Move{
				{Name: "Tackle", Type: "normal", Power: 40, Accuracy: 100},
				{Name: "Growl", Type: "normal", Power: 0, Accuracy: 100},
				{Name: "Vine Whip", Type: "grass", Power: 45, Accuracy: 100},
			},
		},
		{
			ID:       4,
			Name:     "Charmander",
			Types:    []string{"fire"},
			Level:    5,
			HP:       39,
			MaxHP:    39,
			Attack:   52,
			Defense:  43,
			Speed:    65,
			ImageURL: "/assets/sprites/charmander.png",
			Moves: []Move{
				{Name: "Scratch", Type: "normal", Power: 40, Accuracy: 100},
				{Name: "Growl", Type: "normal", Power: 0, Accuracy: 100},
				{Name: "Ember", Type: "fire", Power: 40, Accuracy: 100},
			},
		},
		{
			ID:       7,
			Name:     "Squirtle",
			Types:    []string{"water"},
			Level:    5,
			HP:       44,
			MaxHP:    44,
			Attack:   48,
			Defense:  65,
			Speed:    43,
			ImageURL: "/assets/sprites/squirtle.png",
			Moves: []Move{
				{Name: "Tackle", Type: "normal", Power: 40, Accuracy: 100},
				{Name: "Tail Whip", Type: "normal", Power: 0, Accuracy: 100},
				{Name: "Water Gun", Type: "water", Power: 40, Accuracy: 100},
			},
		},
	}
	
	// Map to store saved games - in a real app, this would be a database
	savedGames   = make(map[string]SaveRequest)
	savedGamesMu sync.RWMutex // protege savedGames
)

// Limites para evitar agotamiento de recursos
const (
	maxBodyBytes = 1 << 20 // 1 MB por peticion
	maxSaves     = 10000   // partidas guardadas en memoria
	maxPokemons  = 100     // pokemon por partida
	maxRounds    = 50      // rondas por simulacion
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create a new ServeMux (router)
	mux := http.NewServeMux()

	// Register handlers - using standard pattern for compatibility with all Go versions
	mux.HandleFunc("/api/pokemon", getPokemonHandler)
	mux.HandleFunc("/api/pokemon/", getPokemonByIDHandler) // Trailing slash for wildcard matching
	mux.HandleFunc("/api/battle/simulate", simulateBattleHandler)
	mux.HandleFunc("/api/game/save", saveGameHandler)
	mux.HandleFunc("/api/game/load/", loadGameHandler) // Trailing slash for wildcard matching
	mux.HandleFunc("/api/health", healthCheckHandler)

	// Set up CORS middleware
	handler := corsMiddleware(mux)
	
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Starting Pokemon game API server on %s", serverAddr)
	srv := &http.Server{
		Addr:              serverAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// Middleware for CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Handle requests for all Pokemon
func getPokemonHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pokedex)
}

// Handle requests for a specific Pokemon by ID
func getPokemonByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	// Extract ID from the path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(w, "Invalid Pokemon ID", http.StatusBadRequest)
		return
	}

	// Find the Pokemon
	for _, pokemon := range pokedex {
		if pokemon.ID == id {
			json.NewEncoder(w).Encode(pokemon)
			return
		}
	}

	// If we get here, the Pokemon wasn't found
	sendErrorResponse(w, fmt.Sprintf("Pokemon with ID %d not found", id), http.StatusNotFound)
}

// Handle battle simulation requests
func simulateBattleHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Parse the request body
	var requestData struct {
		PlayerPokemon   Pokemon `json:"playerPokemon"`
		OpponentPokemon Pokemon `json:"opponentPokemon"`
		Rounds          int     `json:"rounds"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, p := range []Pokemon{requestData.PlayerPokemon, requestData.OpponentPokemon} {
		if err := validatePokemon(p); err != nil {
			sendErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Validate request data
	if requestData.PlayerPokemon.ID == 0 || requestData.OpponentPokemon.ID == 0 {
		sendErrorResponse(w, "Both player and opponent Pokemon are required", http.StatusBadRequest)
		return
	}

	// Set a default number of rounds if not specified
	rounds := requestData.Rounds
	if rounds <= 0 {
		rounds = 5
	}
	if rounds > maxRounds {
		rounds = maxRounds
	}

	// Simulate the battle
	result := simulateBattle(requestData.PlayerPokemon, requestData.OpponentPokemon, rounds)
	
	// Return the result
	json.NewEncoder(w).Encode(result)
}

// Handle game save requests
func saveGameHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Parse the request body (con limite de tamano)
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var saveRequest SaveRequest
	err := json.NewDecoder(r.Body).Decode(&saveRequest)
	if err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request data
	if saveRequest.PlayerID == "" || len(saveRequest.PlayerID) > 64 || len(saveRequest.Pokemons) > maxPokemons {
		sendErrorResponse(w, "Player ID is required", http.StatusBadRequest)
		return
	}

	// Set timestamp if not provided
	if saveRequest.Timestamp == 0 {
		saveRequest.Timestamp = time.Now().Unix()
	}

	// Save the game state
	savedGamesMu.Lock()
	if _, exists := savedGames[saveRequest.PlayerID]; !exists && len(savedGames) >= maxSaves {
		savedGamesMu.Unlock()
		sendErrorResponse(w, "Save storage is full", http.StatusServiceUnavailable)
		return
	}
	savedGames[saveRequest.PlayerID] = saveRequest
	savedGamesMu.Unlock()

	// Return success response
	response := SaveResponse{
		Success:   true,
		SaveID:    saveRequest.PlayerID,
		Timestamp: saveRequest.Timestamp,
		Message:   "Game saved successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// Handle game load requests
func loadGameHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Extract playerID from the path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	playerID := parts[len(parts)-1]

	// Retrieve the saved game
	savedGamesMu.RLock()
	savedGame, exists := savedGames[playerID]
	savedGamesMu.RUnlock()
	if !exists {
		sendErrorResponse(w, "No saved game found for this player", http.StatusNotFound)
		return
	}

	// Return the saved game
	json.NewEncoder(w).Encode(savedGame)
}

// Health check endpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status":  "UP",
		"message": "Pokemon Game API is running",
		"time":    time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}

// Helper function to send error responses
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	errorResponse := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	}
	json.NewEncoder(w).Encode(errorResponse)
}

// Function to simulate a battle between two Pokemon
func simulateBattle(player, opponent Pokemon, maxRounds int) BattleResult {
	result := BattleResult{
		BattleLog:      make([]string, 0),
		PlayerMoves:    make([]string, 0),
		OpponentMoves:  make([]string, 0),
	}

	// Make copies to avoid modifying the originals
	playerPokemon := player
	opponentPokemon := opponent

	// Add starting message
	result.BattleLog = append(result.BattleLog, 
		fmt.Sprintf("Battle started: %s (Lv.%d) vs %s (Lv.%d)", 
			playerPokemon.Name, playerPokemon.Level, 
			opponentPokemon.Name, opponentPokemon.Level))

	// Determine who goes first based on speed
	playerFirst := playerPokemon.Speed >= opponentPokemon.Speed

	// Battle loop
	round := 0
	for round < maxRounds && playerPokemon.HP > 0 && opponentPokemon.HP > 0 {
		round++
		result.BattleLog = append(result.BattleLog, fmt.Sprintf("Round %d:", round))

		// Player's turn if they go first, otherwise opponent's turn
		if playerFirst {
			// Player attacks
			attackResult := executeAttack(&playerPokemon, &opponentPokemon)
			result.BattleLog = append(result.BattleLog, attackResult)
			result.PlayerMoves = append(result.PlayerMoves, attackResult)

			// Check if opponent fainted
			if opponentPokemon.HP <= 0 {
				result.BattleLog = append(result.BattleLog, 
					fmt.Sprintf("%s fainted!", opponentPokemon.Name))
				break
			}

			// Opponent attacks
			attackResult = executeAttack(&opponentPokemon, &playerPokemon)
			result.BattleLog = append(result.BattleLog, attackResult)
			result.OpponentMoves = append(result.OpponentMoves, attackResult)

			// Check if player fainted
			if playerPokemon.HP <= 0 {
				result.BattleLog = append(result.BattleLog, 
					fmt.Sprintf("%s fainted!", playerPokemon.Name))
				break
			}
		} else {
			// Opponent attacks first
			attackResult := executeAttack(&opponentPokemon, &playerPokemon)
			result.BattleLog = append(result.BattleLog, attackResult)
			result.OpponentMoves = append(result.OpponentMoves, attackResult)

			// Check if player fainted
			if playerPokemon.HP <= 0 {
				result.BattleLog = append(result.BattleLog, 
					fmt.Sprintf("%s fainted!", playerPokemon.Name))
				break
			}

			// Player attacks
			attackResult = executeAttack(&playerPokemon, &opponentPokemon)
			result.BattleLog = append(result.BattleLog, attackResult)
			result.PlayerMoves = append(result.PlayerMoves, attackResult)

			// Check if opponent fainted
			if opponentPokemon.HP <= 0 {
				result.BattleLog = append(result.BattleLog, 
					fmt.Sprintf("%s fainted!", opponentPokemon.Name))
				break
			}
		}
	}

	// Determine winner
	if playerPokemon.HP <= 0 {
		result.Winner = "opponent"
		result.BattleLog = append(result.BattleLog, fmt.Sprintf("%s wins the battle!", opponentPokemon.Name))
	} else if opponentPokemon.HP <= 0 {
		result.Winner = "player"
		result.BattleLog = append(result.BattleLog, fmt.Sprintf("%s wins the battle!", playerPokemon.Name))
	} else {
		result.Winner = "draw"
		result.BattleLog = append(result.BattleLog, "The battle ended in a draw!")
	}

	// Set final HP values
	result.PlayerHP = playerPokemon.HP
	result.OpponentHP = opponentPokemon.HP
	result.Rounds = round

	return result
}

// Helper function to execute an attack between two Pokemon
func executeAttack(attacker, defender *Pokemon) string {
	// Select a random move
	moveIndex := rand.Intn(len(attacker.Moves))
	move := attacker.Moves[moveIndex]

	// Check if attack hits based on accuracy
	if rand.Intn(100) >= move.Accuracy {
		return fmt.Sprintf("%s's %s missed!", attacker.Name, move.Name)
	}

	// Calculate damage
	// Simple formula: ((2 * Level / 5 + 2) * Power * Attack / Defense) / 50 + 2
	// Simplified for this example
	damage := 0
	if move.Power > 0 {
		baseDamage := ((2 * attacker.Level / 5 + 2) * move.Power * attacker.Attack / defender.Defense) / 50 + 2
		
		// Add some randomness (85-100% of calculated damage)
		randomFactor := 85 + rand.Intn(16)
		damage = (baseDamage * randomFactor) / 100
		
		// Apply type effectiveness (simplified)
		// Not implementing full type chart here
		typeMultiplier := 1.0
		for _, defenderType := range defender.Types {
			if (move.Type == "water" && defenderType == "fire") ||
			   (move.Type == "fire" && defenderType == "grass") ||
			   (move.Type == "grass" && defenderType == "water") {
				typeMultiplier = 2.0
				break
			} else if (move.Type == "water" && defenderType == "grass") ||
					  (move.Type == "fire" && defenderType == "water") ||
					  (move.Type == "grass" && defenderType == "fire") {
				typeMultiplier = 0.5
				break
			}
		}
		
		damage = int(float64(damage) * typeMultiplier)
		
		// Ensure minimum damage
		if damage < 1 {
			damage = 1
		}
		
		// Apply damage
		defender.HP -= damage
		if defender.HP < 0 {
			defender.HP = 0
		}
		
		// Return attack message
		effectivenessMsg := ""
		if typeMultiplier > 1.0 {
			effectivenessMsg = "It's super effective!"
		} else if typeMultiplier < 1.0 {
			effectivenessMsg = "It's not very effective..."
		}
		
		return fmt.Sprintf("%s used %s! %s %d damage! %s", 
			attacker.Name, move.Name, 
			effectivenessMsg, damage,
			fmt.Sprintf("%s HP: %d/%d", defender.Name, defender.HP, defender.MaxHP))
	}
	
	// Handle non-damaging moves (not implemented in detail)
	return fmt.Sprintf("%s used %s!", attacker.Name, move.Name)
} 

// validatePokemon rechaza datos que causarian panics o abuso de recursos
func validatePokemon(p Pokemon) error {
	switch {
	case p.Level < 1 || p.Level > 100:
		return fmt.Errorf("level must be between 1 and 100")
	case p.HP < 1 || p.HP > 999 || p.MaxHP < 1 || p.MaxHP > 999:
		return fmt.Errorf("hp and maxHp must be between 1 and 999")
	case p.Attack < 0 || p.Attack > 999 || p.Defense < 1 || p.Defense > 999 || p.Speed < 0 || p.Speed > 999:
		return fmt.Errorf("invalid stats")
	case len(p.Moves) == 0 || len(p.Moves) > 4:
		return fmt.Errorf("a pokemon needs between 1 and 4 moves")
	}
	for _, m := range p.Moves {
		if m.Power < 0 || m.Power > 250 || m.Accuracy < 0 || m.Accuracy > 100 {
			return fmt.Errorf("invalid move values")
		}
	}
	return nil
}