package axiom

import (
	"time"
	"github.com/google/uuid"
)

// Claim represents an axiomatic claim in the trust network
type Claim struct {
	Context    string    `json:"@context"`
	Type       string    `json:"type"`
	Issuer     string    `json:"issuer"`
	Issued     time.Time `json:"issued"`
	ClaimBody  Body      `json:"claim"`
	Proof      Proof     `json:"proof"`
}

// Body represents the main content of an axiomatic claim
type Body struct {
	Context  string       `json:"@context"`
	Type     string       `json:"type"`
	Subject  string       `json:"subject"`
	Agent    string       `json:"agent"`
	Tags     []string     `json:"tags"`
	Rating   AxiomRating `json:"axiomRating"`
}

// AxiomRating represents confidence scoring for an axiom
type AxiomRating struct {
	Context         string  `json:"@context"`
	Type           string  `json:"type"`
	MaxConfidence  float64 `json:"maxConfidence"`
	MinConfidence  float64 `json:"minConfidence"`
	ConfidenceValue float64 `json:"confidenceValue"`
	Axiom          string  `json:"axiom"`
}

// Proof represents cryptographic verification of the claim
type Proof struct {
	Type        string    `json:"type"`
	Created     time.Time `json:"created"`
	Verifier    Verifier  `json:"verifier"`
	Domain      string    `json:"domain"`
	ProofValue  string    `json:"proofValue"`
}

// Verifier identifies the entity verifying the claim
type Verifier struct {
	ID string `json:"id"`
} 