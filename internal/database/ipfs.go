package database

import (
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
)

// IPFSRecord stores information about uploaded trust graphs
type IPFSRecord struct {
	ID        uuid.UUID
	IPFSID    string
	Type      string // e.g., "trust_graph", "claim_batch"
	Metadata  map[string]interface{}
	CreatedAt time.Time
}

// StoreIPFSRecord stores a record of an IPFS upload
func (db *DB) StoreIPFSRecord(ctx context.Context, record *IPFSRecord) error {
	_, err := db.pool.Exec(ctx,
		`INSERT INTO ipfs_records (id, ipfs_id, type, metadata, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		record.ID,
		record.IPFSID,
		record.Type,
		record.Metadata,
		record.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to store IPFS record: %w", err)
	}
	
	return nil
}

// GetIPFSRecords retrieves IPFS records by type
func (db *DB) GetIPFSRecords(ctx context.Context, recordType string) ([]*IPFSRecord, error) {
	rows, err := db.pool.Query(ctx,
		`SELECT id, ipfs_id, type, metadata, created_at
		 FROM ipfs_records
		 WHERE type = $1
		 ORDER BY created_at DESC`,
		recordType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query IPFS records: %w", err)
	}
	defer rows.Close()

	var records []*IPFSRecord
	for rows.Next() {
		record := &IPFSRecord{}
		err := rows.Scan(
			&record.ID,
			&record.IPFSID,
			&record.Type,
			&record.Metadata,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan IPFS record: %w", err)
		}
		records = append(records, record)
	}

	return records, nil
} 