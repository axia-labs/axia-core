package database

import (
	"context"
	"github.com/google/uuid"
	"axia/internal/social/twitter"
)

// StoreTwitterReport stores a Twitter report and its associated claim
func (db *DB) StoreTwitterReport(ctx context.Context, report *twitter.TweetReport, claimID uuid.UUID) error {
	_, err := db.pool.Exec(ctx,
		`INSERT INTO twitter_reports 
		(tweet_id, author_id, subject_handle, action, project, amount, claim_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		report.TweetID,
		report.AuthorID,
		report.Subject,
		report.Action,
		report.Project,
		report.Amount,
		claimID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to store twitter report: %w", err)
	}
	
	return nil
}

// GetTwitterReportsByProject retrieves all reports for a specific project
func (db *DB) GetTwitterReportsByProject(ctx context.Context, project string) ([]*twitter.TweetReport, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT tweet_id, author_id, subject_handle, action, project, amount
		 FROM twitter_reports
		 WHERE project = $1
		 ORDER BY created_at DESC`,
		project,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query twitter reports: %w", err)
	}
	defer rows.Close()

	var reports []*twitter.TweetReport
	for rows.Next() {
		report := &twitter.TweetReport{}
		err := rows.Scan(
			&report.TweetID,
			&report.AuthorID,
			&report.Subject,
			&report.Action,
			&report.Project,
			&report.Amount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan twitter report: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, nil
} 