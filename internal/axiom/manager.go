package axiom

import (
	"github.com/sirupsen/logrus"
	"axia/internal/crypto"
	"axia/internal/state"
)

// Manager handles creation and verification of axiomatic claims
type Manager struct {
	proofGen *crypto.ProofGenerator
	state    *state.StateManager
	logger   *logrus.Logger
}

// NewManager creates a new axiom manager
func NewManager(logger *logrus.Logger) *Manager {
	return &Manager{
		proofGen: crypto.NewProofGenerator(),
		state:    state.NewNodeStateManager(logger),
		logger:   logger,
	}
}

// CreateClaim creates a new axiomatic claim
func (m *Manager) CreateClaim(agent, subject, axiom string, confidence float64, tags []string) (*Claim, error) {
	m.logger.WithFields(logrus.Fields{
		"agent":      agent,
		"subject":    subject,
		"confidence": confidence,
	}).Info("Creating new axiomatic claim")

	claim := &Claim{
		Context: "https://schema.axios.ai/AxiomaticClaim.jsonld",
		Type:    "AxiomaticClaim",
		Issuer:  agent,
		Issued:  time.Now().UTC(),
		ClaimBody: Body{
			Context: "https://schema.axios.ai/",
			Type:    "Axiom",
			Subject: subject,
			Agent:   agent,
			Tags:    tags,
			Rating: AxiomRating{
				Context:         "https://schema.axios.ai/",
				Type:           "Confidence",
				MaxConfidence:  1.0,
				MinConfidence:  0.0,
				ConfidenceValue: confidence,
				Axiom:          axiom,
			},
		},
	}

	// Generate cryptographic proof
	proof, err := m.proofGen.GenerateProof(claim)
	if err != nil {
		m.logger.WithError(err).Error("Failed to generate proof for claim")
		return nil, err
	}

	claim.Proof = Proof{
		Type:       "AxiomaticVerification2024",
		Created:    time.Now().UTC(),
		Domain:     "axios.ai",
		ProofValue: proof.Hash,
		Verifier: Verifier{
			ID: "Axiomatic-key:" + proof.Hash[:64],
		},
	}

	return claim, nil
} 