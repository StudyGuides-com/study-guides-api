package question

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

type SqlQuestionStore struct {
	db *pgxpool.Pool
}

type questionRow struct {
	ID              string            `db:"id"`
	BatchID         *string           `db:"batchId"`
	QuestionText    string            `db:"questionText"`
	AnswerText      string            `db:"answerText"`
	Hash            string            `db:"hash"`
	LearnMore       *string           `db:"learnMore"`
	Distractors     []string          `db:"distractors"`
	VideoURL        *string           `db:"videoUrl"`
	ImageURL        *string           `db:"imageUrl"`
	Version         int32             `db:"version"`
	Public          bool              `db:"public"`
	Metadata        map[string]string `db:"metadata"`
	CreatedAt       time.Time         `db:"createdAt"`
	UpdatedAt       time.Time         `db:"updatedAt"`
	CorrectCount    *int32            `db:"correctCount"`
	DifficultyRatio *float64          `db:"difficultyRatio"`
	IncorrectCount  *int32            `db:"incorrectCount"`
	OwnerID         *string           `db:"ownerId"`
	PassageID       *string           `db:"passageId"`
}

func mapRowsToQuestions(rows []questionRow) []*sharedpb.Question {
	var questions []*sharedpb.Question
	for _, row := range rows {
		q := &sharedpb.Question{
			Id:              row.ID,
			BatchId:         row.BatchID,
			QuestionText:    row.QuestionText,
			AnswerText:      row.AnswerText,
			Hash:            row.Hash,
			LearnMore:       row.LearnMore,
			Distractors:     row.Distractors,
			VideoUrl:        row.VideoURL,
			ImageUrl:        row.ImageURL,
			Version:         row.Version,
			Public:          row.Public,
			Metadata:        row.Metadata,
			CreatedAt:       timestamppb.New(row.CreatedAt),
			UpdatedAt:       timestamppb.New(row.UpdatedAt),
			CorrectCount:    row.CorrectCount,
			DifficultyRatio: row.DifficultyRatio,
			IncorrectCount:  row.IncorrectCount,
			OwnerId:         row.OwnerID,
			PassageId:       row.PassageID,
		}
		questions = append(questions, q)
	}
	return questions
}

func (s *SqlQuestionStore) GetQuestionsByTagID(ctx context.Context, tagID string) ([]*sharedpb.Question, error) {
	var rows []questionRow

	err := pgxscan.Select(ctx, s.db, &rows, `
		SELECT 
			q.id, q."batchId", q."questionText", q."answerText", q.hash, q."learnMore",
			q.distractors, q."videoUrl", q."imageUrl", q.version, q.public, q.metadata,
			q."createdAt", q."updatedAt", q."correctCount", q."difficultyRatio",
			q."incorrectCount", q."ownerId", q."passageId"
		FROM "Question" q
		JOIN "QuestionTag" qt ON q.id = qt."questionId"
		WHERE qt."tagId" = $1
	`, tagID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query questions: %v", err)
	}

	return mapRowsToQuestions(rows), nil
}
