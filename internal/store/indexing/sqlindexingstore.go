package indexing

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	"github.com/lucsky/cuid"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
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
	
	// Queue items based on force flag
	if force {
		// Force mode: Queue ALL items for complete rebuild
		if err := s.QueueBatchForReindex(ctx, objectType); err != nil {
			// Update job as failed
			s.updateJobStatus(jobID, "Failed", err.Error(), 0)
			return jobID, err
		}
		fmt.Printf("Force mode: Queued all %s items for complete reindexing\n", objectType)
	} else {
		// Incremental mode: Only queue changed items
		if err := s.QueueChangedForIndex(ctx, objectType); err != nil {
			// Update job as failed
			s.updateJobStatus(jobID, "Failed", err.Error(), 0)
			return jobID, err
		}
		// Note: QueueChangedForIndex already logs how many items were queued
	}
	
	// Start async processing
	go s.runIndexingAsync(jobID, objectType, force)
	
	return jobID, nil
}

// StartIndexingJobWithFilters starts a background indexing job with optional filters
func (s *SqlIndexingStore) StartIndexingJobWithFilters(ctx context.Context, objectType string, force bool, tagTypes []sharedpb.TagType, contextTypes []sharedpb.ContextType) (string, error) {
	jobID := cuid.New()
	now := time.Now()

	// Create metadata including filter information
	metadata := map[string]interface{}{
		"objectType": objectType,
		"force":      force,
	}
	if len(tagTypes) > 0 {
		tagTypeStrings := make([]string, len(tagTypes))
		for i, tagType := range tagTypes {
			tagTypeStrings[i] = tagType.String()
		}
		metadata["tagTypes"] = tagTypeStrings
	}
	if len(contextTypes) > 0 {
		contextTypeStrings := make([]string, len(contextTypes))
		for i, contextType := range contextTypes {
			contextTypeStrings[i] = contextType.String()
		}
		metadata["contextTypes"] = contextTypeStrings
	}
	metadataJSON, _ := json.Marshal(metadata)

	// Insert job record
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO "Job" (id, type, status, description, "startedAt", metadata, "createdAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, jobID, "Index", "Running", fmt.Sprintf("Index %s with filters (force=%t)", objectType, force),
	   now, string(metadataJSON), now, now)

	if err != nil {
		return "", fmt.Errorf("failed to create job record: %w", err)
	}

	// Queue items based on force flag with filters
	if force {
		// Force mode: Queue ALL matching items for complete rebuild
		if err := s.QueueBatchForReindexWithFilters(ctx, objectType, tagTypes, contextTypes); err != nil {
			// Update job as failed
			s.updateJobStatus(jobID, "Failed", err.Error(), 0)
			return jobID, err
		}
		fmt.Printf("Force mode: Queued all matching %s items for complete reindexing\n", objectType)
	} else {
		// Incremental mode: Only queue changed matching items
		if err := s.QueueChangedForIndexWithFilters(ctx, objectType, tagTypes, contextTypes); err != nil {
			// Update job as failed
			s.updateJobStatus(jobID, "Failed", err.Error(), 0)
			return jobID, err
		}
	}

	// Start async processing
	go s.runIndexingAsync(jobID, objectType, force)

	return jobID, nil
}

// StartSingleIndexingJob starts a background indexing job for a single specific item
func (s *SqlIndexingStore) StartSingleIndexingJob(ctx context.Context, objectType, objectID string, force bool) (string, error) {
	jobID := cuid.New()
	now := time.Now()

	// Create metadata
	metadata := map[string]interface{}{
		"objectType": objectType,
		"objectID":   objectID,
		"force":      force,
	}
	metadataJSON, _ := json.Marshal(metadata)

	// Insert job record
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO "Job" (id, type, status, description, "startedAt", metadata, "createdAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, jobID, "Index", "Running", fmt.Sprintf("Index single %s (ID=%s, force=%t)", objectType, objectID, force),
	   now, string(metadataJSON), now, now)

	if err != nil {
		return "", fmt.Errorf("failed to create job record: %w", err)
	}

	// Queue the specific item for indexing
	if err := s.QueueIndexOperation(ctx, objectType, objectID, "upsert"); err != nil {
		// Update job as failed
		s.updateJobStatus(jobID, "Failed", err.Error(), 0)
		return jobID, err
	}

	fmt.Printf("Queued single %s item (ID=%s) for indexing\n", objectType, objectID)

	// Start async processing for this single item
	go s.runSingleIndexingAsync(jobID, objectType, objectID, force)

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

// QueueChangedForIndex queues only changed objects for indexing (incremental mode)
func (s *SqlIndexingStore) QueueChangedForIndex(ctx context.Context, objectType string) error {
	var query string
	
	switch objectType {
	case "Tag":
		// This query identifies tags that need indexing by comparing current state with indexed state
		// It queues tags that either:
		// 1. Don't exist in SearchIndexState (never indexed)
		// 2. Have been updated since last indexing
		// 3. Have had their access permissions changed
		// 4. Have had ancestor tags modified (parent changed)
		query = `
			WITH changed_tags AS (
				SELECT DISTINCT t.id
				FROM "Tag" t
				LEFT JOIN "SearchIndexState" s ON s."objectType" = 'Tag' AND s."objectId" = t.id
				WHERE t.context IS NOT NULL
				AND (
					-- Never indexed
					s."objectId" IS NULL
					-- Or tag itself was updated
					OR t."updatedAt" > COALESCE(s."lastIndexedAt", '1970-01-01'::timestamp)
				)
				
				-- NOTE: TagAccess change detection removed - TagAccess table has no updatedAt field
				-- Future enhancement: implement proper permission change tracking
				
				UNION
				
				-- Tags whose ancestors were modified (affects ancestry chain)
				SELECT DISTINCT child.id
				FROM "Tag" child
				INNER JOIN "Tag" parent ON child."parentTagId" = parent.id
				LEFT JOIN "SearchIndexState" s ON s."objectType" = 'Tag' AND s."objectId" = child.id
				WHERE child.context IS NOT NULL
				AND parent."updatedAt" > COALESCE(s."lastIndexedAt", '1970-01-01'::timestamp)
			)
			INSERT INTO "IndexOutbox" ("objectType", "objectId", action, "queuedAt")
			SELECT 'Tag', id, 'upsert', NOW()
			FROM changed_tags
			ON CONFLICT ("objectType", "objectId") DO UPDATE
			SET action = 'upsert', "queuedAt" = NOW()
		`
	default:
		return fmt.Errorf("unsupported object type for changed index: %s", objectType)
	}
	
	result, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to queue changed items for index: %w", err)
	}
	
	// Log how many items were queued
	if rows, err := result.RowsAffected(); err == nil {
		fmt.Printf("Queued %d changed %s items for indexing\n", rows, objectType)
	}

	return nil
}

// QueueBatchForReindexWithFilters queues all objects of a type for reindexing with optional filters
func (s *SqlIndexingStore) QueueBatchForReindexWithFilters(ctx context.Context, objectType string, tagTypes []sharedpb.TagType, contextTypes []sharedpb.ContextType) error {
	var query string
	var args []interface{}
	argIndex := 1

	switch objectType {
	case "Tag":
		query = `
			INSERT INTO "IndexOutbox" ("objectType", "objectId", action, "queuedAt")
			SELECT 'Tag', id, 'upsert', NOW()
			FROM "Tag"
			WHERE context IS NOT NULL`

		// Add TagType filter if provided
		if len(tagTypes) > 0 {
			query += ` AND type = ANY($` + fmt.Sprintf("%d", argIndex) + `)`

			// Convert enum values to strings
			tagTypeStrings := make([]string, len(tagTypes))
			for i, tagType := range tagTypes {
				tagTypeStrings[i] = tagType.String()
			}
			args = append(args, pq.Array(tagTypeStrings))
			argIndex++
		}

		// Add ContextType filter if provided
		if len(contextTypes) > 0 {
			query += ` AND context = ANY($` + fmt.Sprintf("%d", argIndex) + `)`

			// Convert enum values to strings
			contextTypeStrings := make([]string, len(contextTypes))
			for i, contextType := range contextTypes {
				contextTypeStrings[i] = contextType.String()
			}
			args = append(args, pq.Array(contextTypeStrings))
			argIndex++
		}

		query += `
			ON CONFLICT ("objectType", "objectId") DO UPDATE
			SET action = 'upsert', "queuedAt" = NOW()`

	default:
		return fmt.Errorf("unsupported object type for batch reindex: %s", objectType)
	}

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to queue filtered batch for reindex: %w", err)
	}

	return nil
}

// QueueChangedForIndexWithFilters queues only changed objects for indexing with optional filters
func (s *SqlIndexingStore) QueueChangedForIndexWithFilters(ctx context.Context, objectType string, tagTypes []sharedpb.TagType, contextTypes []sharedpb.ContextType) error {
	var query string
	var args []interface{}
	argIndex := 1

	switch objectType {
	case "Tag":
		query = `
			WITH changed_tags AS (
				SELECT DISTINCT t.id
				FROM "Tag" t
				LEFT JOIN "SearchIndexState" s ON s."objectType" = 'Tag' AND s."objectId" = t.id
				WHERE t.context IS NOT NULL`

		// Add TagType filter if provided
		if len(tagTypes) > 0 {
			query += ` AND t.type = ANY($` + fmt.Sprintf("%d", argIndex) + `)`

			// Convert enum values to strings
			tagTypeStrings := make([]string, len(tagTypes))
			for i, tagType := range tagTypes {
				tagTypeStrings[i] = tagType.String()
			}
			args = append(args, pq.Array(tagTypeStrings))
			argIndex++
		}

		// Add ContextType filter if provided
		if len(contextTypes) > 0 {
			query += ` AND t.context = ANY($` + fmt.Sprintf("%d", argIndex) + `)`

			// Convert enum values to strings
			contextTypeStrings := make([]string, len(contextTypes))
			for i, contextType := range contextTypes {
				contextTypeStrings[i] = contextType.String()
			}
			args = append(args, pq.Array(contextTypeStrings))
			argIndex++
		}

		query += `
				AND (
					-- Never indexed
					s."objectId" IS NULL
					-- Or tag itself was updated
					OR t."updatedAt" > COALESCE(s."lastIndexedAt", '1970-01-01'::timestamp)
				)

				UNION

				-- Tags whose ancestors were modified (affects ancestry chain)
				SELECT DISTINCT child.id
				FROM "Tag" child
				INNER JOIN "Tag" parent ON child."parentTagId" = parent.id
				LEFT JOIN "SearchIndexState" s ON s."objectType" = 'Tag' AND s."objectId" = child.id
				WHERE child.context IS NOT NULL`

		// Apply same filters to the ancestor check
		if len(tagTypes) > 0 {
			query += ` AND child.type = ANY($` + fmt.Sprintf("%d", len(args)) + `)`
		}
		if len(contextTypes) > 0 {
			query += ` AND child.context = ANY($` + fmt.Sprintf("%d", len(args)+1) + `)`
		}

		query += `
				AND parent."updatedAt" > COALESCE(s."lastIndexedAt", '1970-01-01'::timestamp)
			)
			INSERT INTO "IndexOutbox" ("objectType", "objectId", action, "queuedAt")
			SELECT 'Tag', id, 'upsert', NOW()
			FROM changed_tags
			ON CONFLICT ("objectType", "objectId") DO UPDATE
			SET action = 'upsert', "queuedAt" = NOW()`

	default:
		return fmt.Errorf("unsupported object type for changed index: %s", objectType)
	}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to queue filtered changed items for index: %w", err)
	}

	// Log how many items were queued
	if rows, err := result.RowsAffected(); err == nil {
		fmt.Printf("Queued %d changed %s items for indexing (with filters)\n", rows, objectType)
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
		_, err := s.db.ExecContext(context.Background(), `
			UPDATE "Job"
			SET status = $1, "completedAt" = $2, "errorMessge" = $3, "updatedAt" = $4
			WHERE id = $5
		`, status, now, errorMsg, now, jobID)
		if err != nil {
			fmt.Printf("Failed to update job %s status to Failed: %v\n", jobID, err)
		}
	} else {
		// Get original metadata and update with items processed
		var originalMetadata sql.NullString
		err := s.db.QueryRowContext(context.Background(),
			`SELECT metadata FROM "Job" WHERE id = $1`, jobID).Scan(&originalMetadata)
		if err != nil {
			fmt.Printf("Failed to get metadata for job %s: %v\n", jobID, err)
			// Don't return early - still try to update the job status
			originalMetadata.Valid = false
		}

		var metadata map[string]interface{}
		if originalMetadata.Valid && originalMetadata.String != "" {
			if err := json.Unmarshal([]byte(originalMetadata.String), &metadata); err != nil {
				fmt.Printf("Failed to unmarshal metadata for job %s: %v\n", jobID, err)
				metadata = make(map[string]interface{})
			}
		} else {
			metadata = make(map[string]interface{})
		}
		metadata["itemsProcessed"] = itemsProcessed

		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			fmt.Printf("Failed to marshal metadata for job %s: %v\n", jobID, err)
			return
		}

		_, err = s.db.ExecContext(context.Background(), `
			UPDATE "Job"
			SET status = $1, "completedAt" = $2, metadata = $3, progress = 100, "updatedAt" = $4
			WHERE id = $5
		`, status, now, string(metadataJSON), now, jobID)
		if err != nil {
			fmt.Printf("Failed to update job %s status to %s: %v\n", jobID, status, err)
		} else {
			fmt.Printf("Successfully updated job %s status to %s (items processed: %d)\n", jobID, status, itemsProcessed)
		}
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

				// Check if this is a "not found" error or if we've exceeded retry limit
				shouldRemove := false
				if strings.Contains(err.Error(), "no rows in result set") || strings.Contains(err.Error(), "failed to get tag") {
					// Item doesn't exist, remove from queue
					fmt.Printf("Item %s %s not found, removing from queue\n", op.ObjectType, op.ObjectID)
					shouldRemove = true
				} else {
					// Check attempt count
					state, _ := s.GetIndexState(ctx, op.ObjectType, op.ObjectID)
					if state != nil && state.AttemptCount >= 10 {
						fmt.Printf("Item %s %s exceeded retry limit (%d attempts), removing from queue\n", op.ObjectType, op.ObjectID, state.AttemptCount)
						shouldRemove = true
					}
				}

				s.UpdateIndexError(ctx, op.ObjectType, op.ObjectID, err)

				if shouldRemove {
					// Remove from outbox to prevent infinite retries
					if err := s.RemoveFromOutbox(ctx, op.ID); err != nil {
						fmt.Printf("Error removing failed item from outbox: %v\n", err)
					}
				}
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

// runSingleIndexingAsync processes a single item indexing job asynchronously
func (s *SqlIndexingStore) runSingleIndexingAsync(jobID, objectType, objectID string, force bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	itemsProcessed := 0

	// Get the specific operation from outbox for this objectID
	operations, err := s.GetPendingOperations(ctx, objectType, 1000) // Get more to find our specific item
	if err != nil {
		s.updateJobStatus(jobID, "Failed", err.Error(), itemsProcessed)
		return
	}

	// Find the operation for our specific objectID
	var targetOperation *IndexOperation
	for _, op := range operations {
		if op.ObjectID == objectID {
			targetOperation = &op
			break
		}
	}

	if targetOperation == nil {
		s.updateJobStatus(jobID, "Failed", fmt.Sprintf("Operation for %s ID %s not found in queue", objectType, objectID), itemsProcessed)
		return
	}

	// Process the single operation
	if err := s.processOperation(ctx, *targetOperation, force); err != nil {
		fmt.Printf("Error processing single %s %s: %v\n", targetOperation.ObjectType, targetOperation.ObjectID, err)
		s.UpdateIndexError(ctx, targetOperation.ObjectType, targetOperation.ObjectID, err)
		s.updateJobStatus(jobID, "Failed", err.Error(), itemsProcessed)
		return
	}

	// Remove from outbox after successful processing
	if err := s.RemoveFromOutbox(ctx, targetOperation.ID); err != nil {
		fmt.Printf("Error removing from outbox: %v\n", err)
	}

	itemsProcessed = 1

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

// StartPruningJob starts a background job to prune orphaned index objects
func (s *SqlIndexingStore) StartPruningJob(ctx context.Context, objectType string, tagTypes []sharedpb.TagType, contextTypes []sharedpb.ContextType) (string, error) {
	jobID := cuid.New()
	now := time.Now()

	// Create metadata
	metadata := map[string]interface{}{
		"objectType": objectType,
		"pruneCount": 0,
	}
	if len(tagTypes) > 0 {
		metadata["tagTypes"] = tagTypes
	}
	if len(contextTypes) > 0 {
		metadata["contextTypes"] = contextTypes
	}
	metadataJSON, _ := json.Marshal(metadata)

	// Insert job record
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO "Job" (id, type, status, description, "startedAt", metadata, "createdAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, jobID, "Prune", "Running", "Prune orphaned index objects",
	   now, string(metadataJSON), now, now)

	if err != nil {
		return "", fmt.Errorf("failed to create job record: %w", err)
	}

	// Start async pruning
	go s.runPruningAsync(jobID, objectType, tagTypes, contextTypes)

	return jobID, nil
}

// runPruningAsync performs the actual pruning operation asynchronously
func (s *SqlIndexingStore) runPruningAsync(jobID, objectType string, tagTypes []sharedpb.TagType, contextTypes []sharedpb.ContextType) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	startTime := time.Now()
	deletedCount := 0

	defer func() {
		// Always update job status on completion
		duration := int(time.Since(startTime).Seconds())
		if r := recover(); r != nil {
			s.updateJobStatus(jobID, "Failed", fmt.Sprintf("panic: %v", r), duration)
		}
	}()

	// Currently only support Tag pruning
	if objectType != "Tag" {
		s.updateJobStatus(jobID, "Failed", fmt.Sprintf("unsupported object type: %s", objectType), 0)
		return
	}

	// Prepare existence check query with optional filters
	checkQuery := `SELECT 1 FROM "Tag" WHERE id = $1`
	filterArgs := []interface{}{}
	argCount := 1

	if len(tagTypes) > 0 {
		argCount++
		checkQuery += fmt.Sprintf(" AND type = ANY($%d)", argCount)
		filterArgs = append(filterArgs, pq.Array(tagTypes))
	}

	if len(contextTypes) > 0 {
		argCount++
		checkQuery += fmt.Sprintf(" AND context = ANY($%d)", argCount)
		filterArgs = append(filterArgs, pq.Array(contextTypes))
	}

	// Prepare statement for efficient repeated queries
	stmt, err := s.db.PrepareContext(ctx, checkQuery)
	if err != nil {
		s.updateJobStatus(jobID, "Failed", fmt.Sprintf("failed to prepare query: %v", err), 0)
		return
	}
	defer stmt.Close() // Moved immediately after successful preparation

	// Stream through Algolia objects
	it, err := s.tagIndex.BrowseObjects()
	if err != nil {
		s.updateJobStatus(jobID, "Failed", fmt.Sprintf("failed to browse Algolia index: %v", err), 0)
		return
	}

	// Batch for deletion (accumulate up to 1000 IDs)
	deleteBuffer := make([]string, 0, 1000)
	processedCount := 0

	for {
		// Check for context cancellation (timeout or cancellation)
		select {
		case <-ctx.Done():
			fmt.Printf("Pruning job %s cancelled/timed out after processing %d items: %v\n", jobID, processedCount, ctx.Err())
			goto cleanup
		default:
			// Continue processing
		}

		// Get next object from Algolia
		var record struct {
			ObjectID string `json:"objectID"`
		}

		_, err := it.Next(&record)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Warning: Failed to parse Algolia record in job %s: %v\n", jobID, err)
			continue // Skip bad records but log the error
		}

		processedCount++

		// Check if exists in database
		var exists int
		queryArgs := append([]interface{}{record.ObjectID}, filterArgs...)
		err = stmt.QueryRowContext(ctx, queryArgs...).Scan(&exists)

		if err == sql.ErrNoRows {
			// Object doesn't exist in DB, queue for deletion
			deleteBuffer = append(deleteBuffer, record.ObjectID)

			// Delete when buffer is full
			if len(deleteBuffer) >= 1000 {
				_, err := s.tagIndex.DeleteObjects(deleteBuffer)
				if err == nil {
					deletedCount += len(deleteBuffer)
				} else {
					fmt.Printf("Failed to delete batch from Algolia: %v\n", err)
				}
				deleteBuffer = deleteBuffer[:0] // Clear buffer

				// Periodic progress update
				s.updatePruneProgress(jobID, deletedCount)
			}
		}
	}

cleanup:
	// Delete remaining items in buffer
	if len(deleteBuffer) > 0 {
		_, err := s.tagIndex.DeleteObjects(deleteBuffer)
		if err == nil {
			deletedCount += len(deleteBuffer)
		} else {
			fmt.Printf("Failed to delete final batch from Algolia: %v\n", err)
		}
	}

	// Update job as completed
	duration := int(time.Since(startTime).Seconds())
	s.updateJobStatus(jobID, "Completed", "", duration)
	s.updatePruneProgress(jobID, deletedCount)

	fmt.Printf("Pruning job %s completed: processed %d items, deleted %d orphaned objects in %d seconds\n",
		jobID, processedCount, deletedCount, duration)
}

// updatePruneProgress updates the pruning job metadata with deletion count
func (s *SqlIndexingStore) updatePruneProgress(jobID string, deletedCount int) {
	metadata := map[string]interface{}{
		"objectType": "Tag",
		"pruneCount": deletedCount,
	}
	metadataJSON, _ := json.Marshal(metadata)

	_, err := s.db.Exec(`
		UPDATE "Job"
		SET metadata = $1, "updatedAt" = $2
		WHERE id = $3
	`, string(metadataJSON), time.Now(), jobID)

	if err != nil {
		fmt.Printf("Failed to update prune progress for job %s: %v\n", jobID, err)
	}
}