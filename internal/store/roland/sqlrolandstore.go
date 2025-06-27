package roland

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type SqlRolandStore struct {
	db *pgxpool.Pool
}

func GeneratePetAlias(input string) string {
	hash := sha256.Sum256([]byte(input))
	seed := binary.BigEndian.Uint64(hash[:8])
	rand.New(rand.NewSource(int64(seed)))
	return petname.Generate(2, "-")
}

func (s *SqlRolandStore) HealthCheck(ctx context.Context) error {
	return nil
}

func (s *SqlRolandStore) SaveBundle(ctx context.Context, bundle *sharedpb.Bundle, force bool) (bool, error) {
	now := time.Now()
	if bundle.CreatedAt == nil {
		bundle.CreatedAt = timestamppb.New(now)
	}
	bundle.UpdatedAt = timestamppb.New(now)

	shortId := GeneratePetAlias(bundle.Title)

	query := `
		INSERT INTO bundles (
			id, short_id, parser_type, title, payload,
			exported_to_dev, exported_to_test, exported_to_prod,
			created_at, updated_at, assisted_at, gob_payload
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11, $12
		)
		ON CONFLICT (id) DO ` + func() string {
		if force {
			return `UPDATE SET
				parser_type = EXCLUDED.parser_type,
				title = EXCLUDED.title,
				payload = EXCLUDED.payload,
				exported_to_dev = EXCLUDED.exported_to_dev,
				exported_to_test = EXCLUDED.exported_to_test,
				exported_to_prod = EXCLUDED.exported_to_prod,
				updated_at = EXCLUDED.updated_at,
				assisted_at = EXCLUDED.assisted_at,
				gob_payload = EXCLUDED.gob_payload`
		}
		return `NOTHING`
	}()

	result, err := s.db.Exec(ctx, query,
		bundle.Id,
		shortId,
		bundle.ParserType,
		bundle.Title,
		bundle.Payload,
		bundle.ExportedToDev,
		bundle.ExportedToTest,
		bundle.ExportedToProd,
		bundle.CreatedAt,
		bundle.UpdatedAt,
		bundle.AssistedAt,
		bundle.GobPayload,
	)
	if err != nil {
		return false, status.Errorf(codes.Internal, "failed to save bundle: %v", err)
	}

	if force {
		return true, nil
	}

	rowsAffected := result.RowsAffected()
	return rowsAffected > 0, nil
}

func (s *SqlRolandStore) Bundles(ctx context.Context) ([]*sharedpb.Bundle, error) {
	query := `
		SELECT id, short_id, parser_type, title, payload,
			exported_to_dev, exported_to_test, exported_to_prod,
			created_at, updated_at, assisted_at, gob_payload
		FROM bundles
		ORDER BY created_at DESC
	`

	var bundles []*sharedpb.Bundle
	err := pgxscan.Select(ctx, s.db, &bundles, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query bundles: %v", err)
	}

	return bundles, nil
}

func (s *SqlRolandStore) BundlesByParserType(ctx context.Context, parserType sharedpb.ParserType) ([]*sharedpb.Bundle, error) {
	query := `
		SELECT id, parser_type, title, payload,
			exported_to_dev, exported_to_test, exported_to_prod,
			created_at, updated_at, assisted_at, gob_payload
		FROM bundles
		WHERE parser_type = $1
		ORDER BY created_at DESC
	`

	var bundles []*sharedpb.Bundle
	err := pgxscan.Select(ctx, s.db, &bundles, query, parserType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query bundles by parser type: %v", err)
	}

	return bundles, nil
}

// UpdateGob updates the gob_payload for a specific bundle
func (s *SqlRolandStore) UpdateGob(ctx context.Context, id string, gobPayload []byte, force bool) (bool, error) {
	query := `
		UPDATE bundles
		SET gob_payload = $1,
			updated_at = $2,
			assisted_at = $3
		WHERE id = $4
		AND ($5 = true OR gob_payload IS NULL)
	`

	now := time.Now()
	result, err := s.db.Exec(ctx, query, gobPayload, now, now, id, force)
	if err != nil {
		return false, status.Errorf(codes.Internal, "failed to update gob payload: %v", err)
	}

	rowsAffected := result.RowsAffected()
	return rowsAffected > 0, nil
}

// DeleteAllBundles deletes all bundles from the database
func (s *SqlRolandStore) DeleteAllBundles(ctx context.Context) error {
	query := `DELETE FROM bundles`

	_, err := s.db.Exec(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to delete all bundles: %v", err)
	}

	return nil
}

// DeleteBundleByID deletes a bundle by its ID
func (s *SqlRolandStore) DeleteBundleByID(ctx context.Context, id string) error {
	query := `DELETE FROM bundles WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to delete bundle: %v", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return status.Errorf(codes.NotFound, "no bundle found with ID: %s", id)
	}

	return nil
}

// DeleteBundlesByShortID deletes all bundles that start with the given prefix
func (s *SqlRolandStore) DeleteBundlesByShortID(ctx context.Context, prefix string) (int, error) {
	query := `DELETE FROM bundles WHERE short_id = $1`

	result, err := s.db.Exec(ctx, query, prefix)
	if err != nil {
		return 0, status.Errorf(codes.Internal, "failed to delete bundles: %v", err)
	}

	rowsAffected := int(result.RowsAffected())
	if rowsAffected == 0 {
		return 0, status.Errorf(codes.NotFound, "no bundles found with prefix: %s", prefix)
	}

	return rowsAffected, nil
}

func (s *SqlRolandStore) MarkBundleExported(ctx context.Context, id string, exportType sharedpb.ExportType) (bool, error) {
	var query string
	switch exportType {
	case sharedpb.ExportType_EXPORT_TYPE_DEV:
		query = `
			UPDATE bundles
			SET exported_to_dev = true
			WHERE id = $1
		`
	case sharedpb.ExportType_EXPORT_TYPE_TEST:
		query = `
			UPDATE bundles
			SET exported_to_test = true
			WHERE id = $1
		`
	case sharedpb.ExportType_EXPORT_TYPE_PROD:
		query = `
			UPDATE bundles
			SET exported_to_prod = true
			WHERE id = $1
		`
	}
	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return false, status.Errorf(codes.Internal, "failed to mark bundle export: %v", err)
	}
	rowsAffected := result.RowsAffected()
	return rowsAffected > 0, nil
}
