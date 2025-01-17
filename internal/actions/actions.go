package actions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/axia/axia-cli/internal/graph"
)

// TrustClaim represents a trust relationship
type TrustClaim struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
	Timestamp string `json:"timestamp"`
}

// TrustGraph represents the main trust graph structure
type TrustGraph struct {
	Claims []TrustClaim `json:"claims"`
	dbPath string
}

// NewTrustGraph creates a new trust graph instance
func NewTrustGraph() *TrustGraph {
	homeDir, _ := os.UserHomeDir()
	dbPath := filepath.Join(homeDir, ".axia-cli", "trust.json")
	
	// Create directory if it doesn't exist
	os.MkdirAll(filepath.Dir(dbPath), 0755)
	
	t := &TrustGraph{
		Claims: []TrustClaim{},
		dbPath: dbPath,
	}
	
	// Load existing data if available
	t.load()
	return t
}

func (t *TrustGraph) load() error {
	data, err := os.ReadFile(t.dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read trust database: %w", err)
	}

	return json.Unmarshal(data, &t.Claims)
}

func (t *TrustGraph) save() error {
	data, err := json.MarshalIndent(t.Claims, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal trust claims: %w", err)
	}

	return os.WriteFile(t.dbPath, data, 0644)
}

// Map generates a trust map
func (t *TrustGraph) Map() (string, error) {
	if len(t.Claims) == 0 {
		return "No trust claims found.", nil
	}

	var result string
	for _, claim := range t.Claims {
		result += fmt.Sprintf("%s -[%s]-> %s\n", claim.Subject, claim.Predicate, claim.Object)
	}
	return result, nil
}

// Get retrieves trust information
func (t *TrustGraph) Get(id string) ([]TrustClaim, error) {
	var results []TrustClaim
	for _, claim := range t.Claims {
		if claim.Subject == id || claim.Object == id {
			results = append(results, claim)
		}
	}
	return results, nil
}

// Claim creates a new trust claim
func (t *TrustGraph) Claim(subject string, predicate string, object string) error {
	claim := TrustClaim{
		Subject:   subject,
		Predicate: predicate,
		Object:    object,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	t.Claims = append(t.Claims, claim)
	return t.save()
}

// MapDOT generates a DOT format representation of the trust graph
func (t *TrustGraph) MapDOT() (string, error) {
	viz := graph.NewVisualizer()
	for _, claim := range t.Claims {
		viz.AddEdge(claim.Subject, claim.Object, claim.Predicate)
	}
	return viz.GenerateDOT(), nil
}

// MapASCII generates an ASCII art representation of the trust graph
func (t *TrustGraph) MapASCII() (string, error) {
	viz := graph.NewVisualizer()
	for _, claim := range t.Claims {
		viz.AddEdge(claim.Subject, claim.Object, claim.Predicate)
	}
	return viz.GenerateASCII(), nil
}

// Search searches for trust claims matching the given criteria
func (t *TrustGraph) Search(query string) ([]TrustClaim, error) {
	query = strings.ToLower(query)
	var results []TrustClaim
	
	for _, claim := range t.Claims {
		if strings.Contains(strings.ToLower(claim.Subject), query) ||
			strings.Contains(strings.ToLower(claim.Predicate), query) ||
			strings.Contains(strings.ToLower(claim.Object), query) {
			results = append(results, claim)
		}
	}
	
	return results, nil
}

// Stats returns statistics about the trust graph
func (t *TrustGraph) Stats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	// Count unique entities
	entities := make(map[string]bool)
	predicates := make(map[string]int)
	
	for _, claim := range t.Claims {
		entities[claim.Subject] = true
		entities[claim.Object] = true
		predicates[claim.Predicate]++
	}
	
	stats["total_claims"] = len(t.Claims)
	stats["unique_entities"] = len(entities)
	stats["unique_predicates"] = len(predicates)
	
	// Get top predicates
	type predCount struct {
		Predicate string
		Count     int
	}
	
	topPreds := make([]predCount, 0, len(predicates))
	for pred, count := range predicates {
		topPreds = append(topPreds, predCount{pred, count})
	}
	
	sort.Slice(topPreds, func(i, j int) bool {
		return topPreds[i].Count > topPreds[j].Count
	})
	
	if len(topPreds) > 5 {
		topPreds = topPreds[:5]
	}
	stats["top_predicates"] = topPreds
	
	return stats
} 