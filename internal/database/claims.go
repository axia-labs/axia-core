package database

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
	"axia/internal/axiom"
)

// StoreClaim stores a new claim in the database
func (db *DB) StoreClaim(ctx context.Context, claim *axiom.Claim) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var claimID uuid.UUID
	err = tx.QueryRow(ctx,
		`INSERT INTO claims (issuer, subject, axiom_text, confidence, proof_type, proof_value, proof_created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		claim.Issuer,
		claim.ClaimBody.Subject,
		claim.ClaimBody.Rating.Axiom,
		claim.ClaimBody.Rating.ConfidenceValue,
		claim.Proof.Type,
		claim.Proof.ProofValue,
		claim.Proof.Created,
	).Scan(&claimID)
	
	if err != nil {
		return fmt.Errorf("failed to insert claim: %w", err)
	}

	// Store tags
	for _, tag := range claim.ClaimBody.Tags {
		_, err = tx.Exec(ctx,
			`INSERT INTO claim_tags (claim_id, tag) VALUES ($1, $2)`,
			claimID, tag)
		if err != nil {
			return fmt.Errorf("failed to insert tag: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// QueryClaims retrieves claims based on filters
func (db *DB) QueryClaims(ctx context.Context, filters map[string]interface{}) ([]*axiom.Claim, error) {
	query := `
		SELECT DISTINCT c.id, c.issuer, c.subject, c.axiom_text, c.confidence,
		       c.proof_type, c.proof_value, c.proof_created_at
		FROM claims c
		LEFT JOIN claim_tags ct ON c.id = ct.claim_id
		WHERE 1=1
	`
	
	args := make([]interface{}, 0)
	argPos := 1

	if v, ok := filters["issuer"]; ok {
		query += fmt.Sprintf(" AND c.issuer = $%d", argPos)
		args = append(args, v)
		argPos++
	}

	if v, ok := filters["subject"]; ok {
		query += fmt.Sprintf(" AND c.subject = $%d", argPos)
		args = append(args, v)
		argPos++
	}

	if v, ok := filters["tag"]; ok {
		query += fmt.Sprintf(" AND ct.tag = $%d", argPos)
		args = append(args, v)
		argPos++
	}

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query claims: %w", err)
	}
	defer rows.Close()

	var claims []*axiom.Claim
	for rows.Next() {
		claim := &axiom.Claim{}
		err := rows.Scan(
			&claim.ID,
			&claim.Issuer,
			&claim.ClaimBody.Subject,
			&claim.ClaimBody.Rating.Axiom,
			&claim.ClaimBody.Rating.ConfidenceValue,
			&claim.Proof.Type,
			&claim.Proof.ProofValue,
			&claim.Proof.Created,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan claim: %w", err)
		}
		claims = append(claims, claim)
	}

	return claims, nil
} 