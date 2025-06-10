package user

import (
	"context"
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
	GamerTag      *string `db:"gamer_tag"`
	Email         *string `db:"email"`
	EmailVerified *string `db:"email_verified"`
	Image         *string `db:"image"`
	ContentTagID  *string `db:"content_tag_id"`
}

func (s *SqlUserStore) UserByID(ctx context.Context, userID string) (*sharedpb.User, error) {
	var row userRow

	err := pgxscan.Get(ctx, s.db, &row, `
		SELECT id, name, gamer_tag, email, email_verified, image, content_tag_id
		FROM users
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
		SELECT id, name, gamer_tag, email, email_verified, image, content_tag_id
		FROM users
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