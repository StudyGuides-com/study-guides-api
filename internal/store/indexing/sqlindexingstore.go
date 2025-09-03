package indexing

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lucsky/cuid"
	"github.com/studyguides-com/study-guides-api/internal/store/tag"
)

// SqlIndexingStore implements IndexingStore using PostgreSQL
type SqlIndexingStore struct {
	db         *sql.DB
	pool       *pgxpool.Pool
	tagStore   tag.TagStore
	algolia    *search.Client
	tagIndex   *search.Index
}

// NewSqlIndexingStore creates a new SQL-based indexing store
func NewSqlIndexingStore(ctx context.Context, dbURL string, pool *pgxpool.Pool, tagStore tag.TagStore, algoliaAppID, algoliaAPIKey string) (IndexingStore, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Test the connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	// Initialize Algolia client
	algoliaClient := search.NewClient(algoliaAppID, algoliaAPIKey)
	tagIndex := algoliaClient.InitIndex("tags")
	
	return &SqlIndexingStore{
		db:       db,
		pool:     pool,
		tagStore: tagStore,
		algolia:  algoliaClient,
		tagIndex: tagIndex,
	}, nil
}

// StartIndexingJob starts a background indexing job
func (s *SqlIndexingStore) StartIndexingJob(ctx context.Context, objectType string, force bool) (string, error) {
	jobID := cuid.New()
	now := time.Now()
	
	// Create metadata
	metadata := map[string]interface{}{
		"objectType": objectType,
		"force":      force,
	}
	metadataJSON, _ := json.Marshal(metadata)
	
	// Insert job record
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO "Job" (id, type, status, description, "startedAt", metadata, "createdAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, jobID, "Index", "Running", fmt.Sprintf("Index %s (force=%t)", objectType, force), 
	   now, string(metadataJSON), now, now)
	
	if err != nil {
		return "", fmt.Errorf("failed to create job record: %w", err)
	}
	
	// Queue all items for reindexing
	if err := s.QueueBatchForReindex(ctx, objectType); err != nil {
		// Update job as failed
		s.updateJobStatus(jobID, "Failed", err.Error(), 0)
		return jobID, err
	}
	
	// Start async processing
	go s.runIndexingAsync(jobID, objectType, force)
	
	return jobID, nil
}

// GetJobStatus retrieves the status of an indexing job
func (s *SqlIndexingStore) GetJobStatus(ctx context.Context, jobID string) (*JobStatus, error) {
	var job JobStatus
	var metadata sql.NullString
	var description sql.NullString
	var startedAt sql.NullTime
	var completedAt sql.NullTime
	var progress sql.NullInt64
	var durationSeconds sql.NullInt64
	var errorMessage sql.NullString
	
	err := s.db.QueryRowContext(ctx, `
		SELECT id, type, description, status, "startedAt", "completedAt", 
		       progress, "durationSeconds", "errorMessge", metadata
		FROM "Job" 
		WHERE id = $1 AND type = 'Index'
	`, jobID).Scan(
		&job.ID, &job.Type, &description, &job.Status,
		&startedAt, &completedAt, &progress, &durationSeconds, &errorMessage, &metadata,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job %s not found", jobID)
		}
		return nil, fmt.Errorf("failed to get job status: %w", err)
	}
	
	// Set optional fields
	if description.Valid {
		job.Description = description.String
	}
	if startedAt.Valid {
		job.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}
	if progress.Valid {
		job.Progress = int(progress.Int64)
	}
	if durationSeconds.Valid {
		duration := int(durationSeconds.Int64)
		job.DurationSeconds = &duration
	}
	if errorMessage.Valid {
		job.ErrorMessage = &errorMessage.String
	}
	if metadata.Valid && metadata.String != "" {
		var meta map[string]interface{}
		if err := json.Unmarshal([]byte(metadata.String), &meta); err == nil {
			job.Metadata = meta
		}
	}
	
	return &job, nil
}

// ListRecentJobs returns recent indexing jobs for a specific object type
func (s *SqlIndexingStore) ListRecentJobs(ctx context.Context, objectType string) ([]JobStatus, error) {
	query := `
		SELECT id, type, description, status, "startedAt", "completedAt", 
		       progress, "durationSeconds", "errorMessge", metadata
		FROM "Job" 
		WHERE type = 'Index' 
		  AND metadata::jsonb->>'objectType' = $1
		ORDER BY "startedAt" DESC
		LIMIT 20
	`
	
	return s.queryJobs(ctx, query, objectType)
}

// ListRunningJobs returns all currently running indexing jobs
func (s *SqlIndexingStore) ListRunningJobs(ctx context.Context) ([]JobStatus, error) {
	query := `
		SELECT id, type, description, status, "startedAt", "completedAt", 
		       progress, "durationSeconds", "errorMessge", metadata
		FROM "Job" 
		WHERE type = 'Index' AND status = 'Running'
		ORDER BY "startedAt" DESC
	`
	
	return s.queryJobs(ctx, query)
}

// QueueIndexOperation adds an operation to the indexing outbox
func (s *SqlIndexingStore) QueueIndexOperation(ctx context.Context, objectType, objectID, action string) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO "IndexOutbox" ("objectType", "objectId", action, "queuedAt")
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT ("objectType", "objectId") DO UPDATE
		SET action = $3, "queuedAt" = NOW()
	`, objectType, objectID, action)
	
	if err != nil {
		return fmt.Errorf("failed to queue index operation: %w", err)
	}
	
	return nil
}

// GetPendingOperations retrieves pending operations from the outbox
func (s *SqlIndexingStore) GetPendingOperations(ctx context.Context, objectType string, limit int) ([]IndexOperation, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, "objectType", "objectId", action, "queuedAt"
		FROM "IndexOutbox"
		WHERE "objectType" = $1
		ORDER BY "queuedAt", id
		LIMIT $2
	`, objectType, limit)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get pending operations: %w", err)
	}
	defer rows.Close()
	
	var operations []IndexOperation
	for rows.Next() {
		var op IndexOperation
		if err := rows.Scan(&op.ID, &op.ObjectType, &op.ObjectID, &op.Action, &op.QueuedAt); err != nil {
			return nil, fmt.Errorf("failed to scan operation: %w", err)
		}
		operations = append(operations, op)
	}
	
	return operations, nil
}

// RemoveFromOutbox removes a processed operation from the outbox
func (s *SqlIndexingStore) RemoveFromOutbox(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM "IndexOutbox" WHERE id = $1
	`, id)
	
	if err != nil {
		return fmt.Errorf("failed to remove from outbox: %w", err)
	}
	
	return nil
}

// GetIndexState retrieves the current indexing state of an object
func (s *SqlIndexingStore) GetIndexState(ctx context.Context, objectType, objectID string) (*IndexState, error) {
	var state IndexState
	var lastIndexedAt sql.NullTime
	var lastIndexedHash []byte
	var lastAttemptAt sql.NullTime
	var lastError sql.NullString
	
	err := s.db.QueryRowContext(ctx, `
		SELECT "objectType", "objectId", "lastIndexedAt", "lastIndexedHash", 
		       "lastAttemptAt", "attemptCount", "lastError"
		FROM "SearchIndexState"
		WHERE "objectType" = $1 AND "objectId" = $2
	`, objectType, objectID).Scan(
		&state.ObjectType, &state.ObjectID, &lastIndexedAt, &lastIndexedHash,
		&lastAttemptAt, &state.AttemptCount, &lastError,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No state exists yet
		}
		return nil, fmt.Errorf("failed to get index state: %w", err)
	}
	
	// Set optional fields
	if lastIndexedAt.Valid {
		state.LastIndexedAt = &lastIndexedAt.Time
	}
	state.LastIndexedHash = lastIndexedHash
	if lastAttemptAt.Valid {
		state.LastAttemptAt = &lastAttemptAt.Time
	}
	if lastError.Valid {
		state.LastError = &lastError.String
	}
	
	return &state, nil
}

// UpdateIndexState updates the indexing state after successful indexing
func (s *SqlIndexingStore) UpdateIndexState(ctx context.Context, objectType, objectID string, hash []byte) error {
	now := time.Now()
	
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO "SearchIndexState" ("objectType", "objectId", "lastIndexedAt", "lastIndexedHash", "lastAttemptAt", "attemptCount")
		VALUES ($1, $2, $3, $4, $3, 0)
		ON CONFLICT ("objectType", "objectId") DO UPDATE
		SET "lastIndexedAt" = $3, "lastIndexedHash" = $4, "lastAttemptAt" = $3, 
		    "attemptCount" = 0, "lastError" = NULL
	`, objectType, objectID, now, hash)
	
	if err != nil {
		return fmt.Errorf("failed to update index state: %w", err)
	}
	
	return nil
}

// UpdateIndexError updates the indexing state after a failed attempt
func (s *SqlIndexingStore) UpdateIndexError(ctx context.Context, objectType, objectID string, indexErr error) error {
	now := time.Now()
	errMsg := indexErr.Error()
	
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO "SearchIndexState" ("objectType", "objectId", "lastAttemptAt", "attemptCount", "lastError")
		VALUES ($1, $2, $3, 1, $4)
		ON CONFLICT ("objectType", "objectId") DO UPDATE
		SET "lastAttemptAt" = $3, "attemptCount" = "SearchIndexState"."attemptCount" + 1, "lastError" = $4
	`, objectType, objectID, now, errMsg)
	
	if err != nil {
		return fmt.Errorf("failed to update index error: %w", err)
	}
	
	return nil
}

// QueueBatchForReindex queues all objects of a type for reindexing
func (s *SqlIndexingStore) QueueBatchForReindex(ctx context.Context, objectType string) error {
	var query string
	
	switch objectType {
	case "Tag":
		query = `
			INSERT INTO "IndexOutbox" ("objectType", "objectId", action, "queuedAt")
			SELECT 'Tag', id, 'upsert', NOW()
			FROM "Tag"
			WHERE context IS NOT NULL
			ON CONFLICT ("objectType", "objectId") DO UPDATE
			SET action = 'upsert', "queuedAt" = NOW()
		`
	default:
		return fmt.Errorf("unsupported object type for batch reindex: %s", objectType)
	}
	
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to queue batch for reindex: %w", err)
	}
	
	return nil
}

// Helper function to query jobs
func (s *SqlIndexingStore) queryJobs(ctx context.Context, query string, args ...interface{}) ([]JobStatus, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %w", err)
	}
	defer rows.Close()
	
	var jobs []JobStatus
	for rows.Next() {
		var job JobStatus
		var metadata sql.NullString
		var description sql.NullString
		var startedAt sql.NullTime
		var completedAt sql.NullTime
		var progress sql.NullInt64
		var durationSeconds sql.NullInt64
		var errorMessage sql.NullString
		
		if err := rows.Scan(
			&job.ID, &job.Type, &description, &job.Status,
			&startedAt, &completedAt, &progress, &durationSeconds, 
			&errorMessage, &metadata,
		); err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		
		// Set optional fields
		if description.Valid {
			job.Description = description.String
		}
		if startedAt.Valid {
			job.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}
		if progress.Valid {
			job.Progress = int(progress.Int64)
		}
		if durationSeconds.Valid {
			duration := int(durationSeconds.Int64)
			job.DurationSeconds = &duration
		}
		if errorMessage.Valid {
			job.ErrorMessage = &errorMessage.String
		}
		if metadata.Valid && metadata.String != "" {
			var meta map[string]interface{}
			if err := json.Unmarshal([]byte(metadata.String), &meta); err == nil {
				job.Metadata = meta
			}
		}
		
		jobs = append(jobs, job)
	}
	
	return jobs, nil
}

// updateJobStatus updates a job's status in the database
func (s *SqlIndexingStore) updateJobStatus(jobID, status, errorMsg string, itemsProcessed int) {
	now := time.Now()
	
	if status == "Failed" && errorMsg != "" {
		s.db.ExecContext(context.Background(), `
			UPDATE "Job" 
			SET status = $1, "completedAt" = $2, "errorMessge" = $3, "updatedAt" = $4
			WHERE id = $5
		`, status, now, errorMsg, now, jobID)
	} else {
		// Get original metadata and update with items processed
		var originalMetadata sql.NullString
		s.db.QueryRowContext(context.Background(), 
			`SELECT metadata FROM "Job" WHERE id = $1`, jobID).Scan(&originalMetadata)
		
		var metadata map[string]interface{}
		if originalMetadata.Valid && originalMetadata.String != "" {
			json.Unmarshal([]byte(originalMetadata.String), &metadata)
		} else {
			metadata = make(map[string]interface{})
		}
		metadata["itemsProcessed"] = itemsProcessed
		
		metadataJSON, _ := json.Marshal(metadata)
		
		s.db.ExecContext(context.Background(), `
			UPDATE "Job" 
			SET status = $1, "completedAt" = $2, metadata = $3, progress = 100, "updatedAt" = $4
			WHERE id = $5
		`, status, now, string(metadataJSON), now, jobID)
	}
}

// runIndexingAsync processes indexing operations in the background
func (s *SqlIndexingStore) runIndexingAsync(jobID, objectType string, force bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	
	itemsProcessed := 0
	batchSize := 100
	
	for {
		// Get batch from outbox
		operations, err := s.GetPendingOperations(ctx, objectType, batchSize)
		if err != nil {
			s.updateJobStatus(jobID, "Failed", err.Error(), itemsProcessed)
			return
		}
		
		if len(operations) == 0 {
			break // No more operations
		}
		
		for _, op := range operations {
			if err := s.processOperation(ctx, op, force); err != nil {
				// Log error but continue processing
				fmt.Printf("Error processing %s %s: %v\n", op.ObjectType, op.ObjectID, err)
				s.UpdateIndexError(ctx, op.ObjectType, op.ObjectID, err)
				continue
			}
			
			// Remove from outbox after successful processing
			if err := s.RemoveFromOutbox(ctx, op.ID); err != nil {
				fmt.Printf("Error removing from outbox: %v\n", err)
			}
			
			itemsProcessed++
		}
	}
	
	// Update job as completed
	s.updateJobStatus(jobID, "Completed", "", itemsProcessed)
}

// processOperation processes a single indexing operation
func (s *SqlIndexingStore) processOperation(ctx context.Context, op IndexOperation, force bool) error {
	switch op.ObjectType {
	case "Tag":
		return s.processTagOperation(ctx, op, force)
	// Future: Add User, Contact, FAQ cases
	default:
		return fmt.Errorf("unsupported object type: %s", op.ObjectType)
	}
}