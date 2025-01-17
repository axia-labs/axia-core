package actions

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrustGraph(t *testing.T) {
	// Create temporary directory for test database
	tmpDir, err := os.MkdirTemp("", "trust-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create test graph
	graph := &TrustGraph{
		Claims: []TrustClaim{},
		dbPath: filepath.Join(tmpDir, "trust.json"),
	}

	// Test adding claims
	err = graph.Claim("Alice", "friend", "Bob")
	assert.NoError(t, err)
	err = graph.Claim("Bob", "colleague", "Charlie")
	assert.NoError(t, err)

	// Test getting claims
	claims, err := graph.Get("Bob")
	assert.NoError(t, err)
	assert.Len(t, claims, 2)

	// Test searching
	results, err := graph.Search("friend")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Alice", results[0].Subject)

	// Test stats
	stats := graph.Stats()
	assert.Equal(t, 2, stats["total_claims"])
	assert.Equal(t, 3, stats["unique_entities"])
	assert.Equal(t, 2, stats["unique_predicates"])
} 