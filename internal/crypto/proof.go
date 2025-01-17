package crypto

import (
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/sha3"
)

// Proof represents a cryptographic proof for a graph element
type Proof struct {
	Hash      string            `json:"hash"`
	Timestamp int64             `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
}

// ProofGenerator handles creation of cryptographic proofs
type ProofGenerator struct {
	algorithm string
}

// NewProofGenerator creates a new proof generator
func NewProofGenerator() *ProofGenerator {
	return &ProofGenerator{
		algorithm: "sha3-256",
	}
}

// GenerateProof creates a cryptographic proof for any data
func (pg *ProofGenerator) GenerateProof(data interface{}) (*Proof, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	hash := sha3.New256()
	hash.Write(jsonData)
	
	return &Proof{
		Hash:      hex.EncodeToString(hash.Sum(nil)),
		Timestamp: time.Now().Unix(),
		Metadata: map[string]string{
			"algorithm": pg.algorithm,
			"type":     "sha3-256",
		},
	}, nil
} 