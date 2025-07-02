package user

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type SqlUserStore struct {
	db *pgxpool.Pool
}

type userRow struct {
	ID            string  `db:"id"`
	Name          *string `db:"name"`
	GamerTag      *string `db:"gamerTag"`
	Email         *string `db:"email"`
	EmailVerified *string `db:"emailVerified"`
	Image         *string `db:"image"`
	ContentTagID  *string `db:"contentTagId"`
}

func (s *SqlUserStore) UserByID(ctx context.Context, userID string) (*sharedpb.User, error) {
	var row userRow

	err := pgxscan.Get(ctx, s.db, &row, `
		SELECT id, name, "gamerTag", email, "emailVerified", image, "contentTagId"
		FROM public."User"
		WHERE id = $1
	`, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "get user by id: "+err.Error())
	}

	user := &sharedpb.User{
		Id: row.ID,
	}
	if row.Name != nil {
		user.Name = row.Name
	}
	if row.GamerTag != nil {
		user.GamerTag = row.GamerTag
	}
	if row.Email != nil {
		user.Email = row.Email
	}
	if row.EmailVerified != nil {
		parsedTime, err := time.Parse(time.RFC3339, *row.EmailVerified)
		if err != nil {
			return nil, status.Error(codes.Internal, "parse email verified timestamp: "+err.Error())
		}
		user.EmailVerified = timestamppb.New(parsedTime)
	}
	if row.Image != nil {
		user.Image = row.Image
	}
	if row.ContentTagID != nil {
		user.ContentTagId = row.ContentTagID
	}

	return user, nil
}

func (s *SqlUserStore) UserByEmail(ctx context.Context, email string) (*sharedpb.User, error) {
	var row userRow

	err := pgxscan.Get(ctx, s.db, &row, `
		SELECT id, name, "gamerTag", email, "emailVerified", image, "contentTagId"
		FROM public."User"
		WHERE email = $1
	`, email)
	if err != nil {
		return nil, status.Error(codes.Internal, "get user by email: "+err.Error())
	}

	user := &sharedpb.User{
		Id: row.ID,
	}
	if row.Name != nil {
		user.Name = row.Name
	}
	if row.GamerTag != nil {
		user.GamerTag = row.GamerTag
	}
	if row.Email != nil {
		user.Email = row.Email
	}
	if row.EmailVerified != nil {
		parsedTime, err := time.Parse(time.RFC3339, *row.EmailVerified)
		if err != nil {
			return nil, status.Error(codes.Internal, "parse email verified timestamp: "+err.Error())
		}
		user.EmailVerified = timestamppb.New(parsedTime)
	}
	if row.Image != nil {
		user.Image = row.Image
	}
	if row.ContentTagID != nil {
		user.ContentTagId = row.ContentTagID
	}

	return user, nil
}

func (s *SqlUserStore) Profile(ctx context.Context, userID string) (*sharedpb.User, error) {
	// For now, we'll just return the user data as the profile
	return s.UserByID(ctx, userID)
}

func (s *SqlUserStore) UserCount(ctx context.Context, params map[string]string) (int64, error) {
	query := `SELECT COUNT(*) FROM public."User" WHERE 1=1`
	var args []interface{}
	argIndex := 1

	// Handle time-based filters
	if since, ok := params["since"]; ok && since != "" {
		query += fmt.Sprintf(` AND "createdAt" >= $%d`, argIndex)
		args = append(args, since)
		argIndex++
	}

	if until, ok := params["until"]; ok && until != "" {
		query += fmt.Sprintf(` AND "createdAt" <= $%d`, argIndex)
		args = append(args, until)
		argIndex++
	}

	if days, ok := params["days"]; ok && days != "" {
		query += fmt.Sprintf(` AND "createdAt" >= NOW() - INTERVAL '%s days'`, days)
	}

	if months, ok := params["months"]; ok && months != "" {
		query += fmt.Sprintf(` AND "createdAt" >= NOW() - INTERVAL '%s months'`, months)
	}

	if years, ok := params["years"]; ok && years != "" {
		query += fmt.Sprintf(` AND "createdAt" >= NOW() - INTERVAL '%s years'`, years)
	}

	if month, ok := params["month"]; ok && month != "" {
		query += fmt.Sprintf(` AND EXTRACT(MONTH FROM "createdAt") = $%d`, argIndex)
		args = append(args, month)
		argIndex++
	}

	if year, ok := params["year"]; ok && year != "" {
		query += fmt.Sprintf(` AND EXTRACT(YEAR FROM "createdAt") = $%d`, argIndex)
		args = append(args, year)
		argIndex++
	}

	var count int64
	err := s.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, status.Error(codes.Internal, "count users: "+err.Error())
	}

	return count, nil
}

func (s *SqlUserStore) KillUser(ctx context.Context, email string) (bool, error) {
	result, err := s.db.Exec(ctx, `
		DELETE FROM public."User"
		WHERE email = $1
	`, email)
	if err != nil {
		return false, status.Error(codes.Internal, "kill user: "+err.Error())
	}

	return result.RowsAffected() > 0, nil
}
