package twitter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"github.com/sirupsen/logrus"
	"axia/internal/axiom"
	"axia/internal/trust"
	"axia/internal/auth"
)

type TweetData struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	AuthorID  string `json:"author_id"`
	CreatedAt string `json:"created_at"`
}

type WebhookPayload struct {
	TweetCreateEvents []TweetData `json:"tweet_create_events"`
}

type Handler struct {
	manager *axiom.Manager
	network *trust.Network
	logger  *logrus.Logger
	auth    *auth.Authenticator
}

func NewHandler(manager *axiom.Manager, network *trust.Network, logger *logrus.Logger) *Handler {
	return &Handler{
		manager: manager,
		network: network,
		logger:  logger,
	}
}

// HandleWebhook processes incoming Twitter webhook events
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("X-Axia-Key")
	if !h.auth.ValidateKey(key) {
		h.logger.Warn("Unauthorized webhook request")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.logger.WithError(err).Error("Failed to decode webhook payload")
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	for _, tweet := range payload.TweetCreateEvents {
		if err := h.processTweet(tweet); err != nil {
			h.logger.WithError(err).Error("Failed to process tweet")
			continue
		}
	}

	w.WriteHeader(http.StatusOK)
}

// processTweet parses tweet content and creates corresponding trust claims
func (h *Handler) processTweet(tweet TweetData) error {
	h.logger.WithField("tweet_id", tweet.ID).Info("Processing tweet")

	// Parse tweet format: @axia_terminal #report @xyz rugged $arc at 50m mc
	report, err := h.parseTweetReport(tweet.Text)
	if err != nil {
		return fmt.Errorf("failed to parse tweet: %w", err)
	}

	// Create axiomatic claim
	claim, err := h.manager.CreateClaim(
		fmt.Sprintf("twitter:%s", tweet.AuthorID),
		fmt.Sprintf("project:%s", report.Project),
		report.generateAxiom(),
		report.calculateConfidence(),
		[]string{"twitter", "report", report.Action},
	)
	if err != nil {
		return fmt.Errorf("failed to create claim: %w", err)
	}

	// Add claim to trust network
	if err := h.network.AddClaim(claim); err != nil {
		return fmt.Errorf("failed to add claim to network: %w", err)
	}

	return nil
}

type TweetReport struct {
	Reporter string  // Twitter handle of reporter
	Subject  string  // Twitter handle being reported
	Action   string  // Action being reported (e.g., "rugged")
	Project  string  // Project symbol/name
	Amount   string  // Amount involved
	Context  string  // Additional context
}

func (h *Handler) parseTweetReport(text string) (*TweetReport, error) {
	// Regular expression to match tweet format
	pattern := regexp.MustCompile(`@axia_terminal\s+#report\s+@(\w+)\s+(\w+)\s+\$(\w+)(?:\s+at\s+(\d+[kmb])\s+mc)?`)
	
	matches := pattern.FindStringSubmatch(text)
	if matches == nil {
		return nil, fmt.Errorf("invalid tweet format")
	}

	report := &TweetReport{
		Subject:  matches[1],
		Action:   matches[2],
		Project:  matches[3],
		Amount:   matches[4],
	}

	return report, nil
}

func (r *TweetReport) generateAxiom() string {
	axiom := fmt.Sprintf("%s performed '%s' action on project %s", r.Subject, r.Action, r.Project)
	if r.Amount != "" {
		axiom += fmt.Sprintf(" with market cap of %s", r.Amount)
	}
	return axiom
}

func (r *TweetReport) calculateConfidence() float64 {
	// Simple confidence calculation - can be made more sophisticated
	// Base confidence of 0.7 for all reports
	confidence := 0.7

	// Adjust based on context provided
	if r.Amount != "" {
		confidence += 0.1 // More confidence if amount is specified
	}

	return confidence
} 