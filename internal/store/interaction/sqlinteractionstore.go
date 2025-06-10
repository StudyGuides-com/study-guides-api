package interaction

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucsky/cuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	interactionpb "github.com/studyguides-com/study-guides-api/api/v1/interaction"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type SqlInteractionStore struct {
	db *pgxpool.Pool
}

// calculateDifficultyRatio calculates the new difficulty ratio based on correct and incorrect counts
func calculateDifficultyRatio(currentCorrect, currentIncorrect int64, isCorrect bool) (int64, int64, float64) {
	newCorrect := currentCorrect
	newIncorrect := currentIncorrect

	if isCorrect {
		newCorrect++
	} else {
		newIncorrect++
	}

	totalAttempts := newCorrect + newIncorrect
	var newDifficultyRatio float64
	if totalAttempts > 0 {
		newDifficultyRatio = float64(newCorrect) / float64(totalAttempts)
	}

	return newCorrect, newIncorrect, newDifficultyRatio
}

func (s *SqlInteractionStore) AnswerCorrectly(ctx context.Context, req *interactionpb.InteractRequest) error {
	// Start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return status.Error(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Get current question stats
	var currentCorrect, currentIncorrect int64
	err = tx.QueryRow(ctx, `
		SELECT "correctCount", "incorrectCount" 
		FROM "Question" 
		WHERE id = $1
	`, req.QuestionId).Scan(&currentCorrect, &currentIncorrect)
	if err != nil {
		return status.Error(codes.Internal, "failed to fetch question stats")
	}

	// Calculate new difficulty metrics
	newCorrect, newIncorrect, newDifficultyRatio := calculateDifficultyRatio(currentCorrect, currentIncorrect, true)

	// Update question difficulty
	_, err = tx.Exec(ctx, `
		UPDATE "Question" 
		SET "correctCount" = $1, 
			"incorrectCount" = $2, 
			"difficultyRatio" = $3,
			"updatedAt" = $4
		WHERE id = $5
	`, newCorrect, newIncorrect, newDifficultyRatio, time.Now(), req.QuestionId)
	if err != nil {
		return status.Error(codes.Internal, "failed to update question difficulty")
	}

	// Create interaction record
	interactionId := cuid.New()
	metadata := map[string]interface{}{
		"studyMethod": req.StudyMethod,
		"strengthScore": 0.0,
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal metadata")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "UserQuestionInteraction" (
			id, "userId", "questionId", type, "studyMethod", 
			correct, "strengthScore", metadata, "occurredAt"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, interactionId, req.UserId, req.QuestionId, sharedpb.InteractionType_INTERACTION_TYPE_ANSWER_CORRECTLY, req.StudyMethod,
		true, 0.0, metadataBytes, time.Now())
	if err != nil {
		return status.Error(codes.Internal, "failed to create interaction record")
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return status.Error(codes.Internal, "failed to commit transaction")
	}

	return nil
}

func (s *SqlInteractionStore) AnswerIncorrectly(ctx context.Context, req *interactionpb.InteractRequest) error {
	// Start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return status.Error(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Get current question stats
	var currentCorrect, currentIncorrect int64
	err = tx.QueryRow(ctx, `
		SELECT "correctCount", "incorrectCount" 
		FROM "Question" 
		WHERE id = $1
	`, req.QuestionId).Scan(&currentCorrect, &currentIncorrect)
	if err != nil {
		return status.Error(codes.Internal, "failed to fetch question stats")
	}

	// Calculate new difficulty metrics
	newCorrect, newIncorrect, newDifficultyRatio := calculateDifficultyRatio(currentCorrect, currentIncorrect, false)

	// Update question difficulty
	_, err = tx.Exec(ctx, `
		UPDATE "Question" 
		SET "correctCount" = $1, 
			"incorrectCount" = $2, 
			"difficultyRatio" = $3,
			"updatedAt" = $4
		WHERE id = $5
	`, newCorrect, newIncorrect, newDifficultyRatio, time.Now(), req.QuestionId)
	if err != nil {
		return status.Error(codes.Internal, "failed to update question difficulty")
	}

	// Create interaction record
	interactionId := cuid.New()
	metadata := map[string]interface{}{
		"studyMethod": req.StudyMethod,
		"strengthScore": 0.0,
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal metadata")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "UserQuestionInteraction" (
			id, "userId", "questionId", type, "studyMethod", 
			correct, "strengthScore", metadata, "occurredAt"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, interactionId, req.UserId, req.QuestionId, sharedpb.InteractionType_INTERACTION_TYPE_ANSWER_INCORRECTLY, req.StudyMethod,
		false, 0.0, metadataBytes, time.Now())
	if err != nil {
		return status.Error(codes.Internal, "failed to create interaction record")
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return status.Error(codes.Internal, "failed to commit transaction")
	}

	return nil
}

func (s *SqlInteractionStore) AnswerEasy(ctx context.Context, req *interactionpb.InteractRequest) error {
	// Start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return status.Error(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Create interaction record
	interactionId := cuid.New()
	metadata := map[string]interface{}{
		"studyMethod": req.StudyMethod,
		"strengthScore": 0.0,
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal metadata")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "UserQuestionInteraction" (
			id, "userId", "questionId", type, "studyMethod", 
			correct, "strengthScore", metadata, "occurredAt"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, interactionId, req.UserId, req.QuestionId, sharedpb.InteractionType_INTERACTION_TYPE_ANSWER_EASY, req.StudyMethod,
		true, 0.0, metadataBytes, time.Now())
	if err != nil {
		return status.Error(codes.Internal, "failed to create interaction record")
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return status.Error(codes.Internal, "failed to commit transaction")
	}

	return nil
}

func (s *SqlInteractionStore) AnswerHard(ctx context.Context, req *interactionpb.InteractRequest) error {
	// Start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return status.Error(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Create interaction record
	interactionId := cuid.New()
	metadata := map[string]interface{}{
		"studyMethod": req.StudyMethod,
		"strengthScore": 0.0,
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal metadata")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "UserQuestionInteraction" (
			id, "userId", "questionId", type, "studyMethod", 
			correct, "strengthScore", metadata, "occurredAt"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, interactionId, req.UserId, req.QuestionId, sharedpb.InteractionType_INTERACTION_TYPE_ANSWER_HARD, req.StudyMethod,
		false, 0.0, metadataBytes, time.Now())
	if err != nil {
		return status.Error(codes.Internal, "failed to create interaction record")
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return status.Error(codes.Internal, "failed to commit transaction")
	}

	return nil
}

func (s *SqlInteractionStore) Reveal(ctx context.Context, req *interactionpb.InteractRequest) error {
	// Start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return status.Error(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Create interaction record
	interactionId := cuid.New()
	metadata := map[string]interface{}{
		"studyMethod": req.StudyMethod,
		"strengthScore": 0.0,
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal metadata")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "UserQuestionInteraction" (
			id, "userId", "questionId", type, "studyMethod", 
			correct, "strengthScore", metadata, "occurredAt"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, interactionId, req.UserId, req.QuestionId, sharedpb.InteractionType_INTERACTION_TYPE_REVEAL, req.StudyMethod,
		nil, 0.0, metadataBytes, time.Now())
	if err != nil {
		return status.Error(codes.Internal, "failed to create interaction record")
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return status.Error(codes.Internal, "failed to commit transaction")
	}

	return nil
}

func (s *SqlInteractionStore) ViewLearnMore(ctx context.Context, req *interactionpb.InteractRequest) error {
	// Start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return status.Error(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Create interaction record
	interactionId := cuid.New()
	metadata := map[string]interface{}{
		"studyMethod": req.StudyMethod,
		"strengthScore": 0.0,
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal metadata")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "UserQuestionInteraction" (
			id, "userId", "questionId", type, "studyMethod", 
			correct, "strengthScore", metadata, "occurredAt"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, interactionId, req.UserId, req.QuestionId, sharedpb.InteractionType_INTERACTION_TYPE_VIEW_LEARN_MORE, req.StudyMethod,
		nil, 0.0, metadataBytes, time.Now())
	if err != nil {
		return status.Error(codes.Internal, "failed to create interaction record")
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return status.Error(codes.Internal, "failed to commit transaction")
	}

	return nil
}

func (s *SqlInteractionStore) ViewPassage(ctx context.Context, req *interactionpb.InteractRequest) error {
	// Start a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return status.Error(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Create interaction record
	interactionId := cuid.New()
	metadata := map[string]interface{}{
		"studyMethod": req.StudyMethod,
		"strengthScore": 0.0,
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal metadata")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO "UserQuestionInteraction" (
			id, "userId", "questionId", type, "studyMethod", 
			correct, "strengthScore", metadata, "occurredAt"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, interactionId, req.UserId, req.QuestionId, sharedpb.InteractionType_INTERACTION_TYPE_VIEW_PASSAGE, req.StudyMethod,
		nil, 0.0, metadataBytes, time.Now())
	if err != nil {
		return status.Error(codes.Internal, "failed to create interaction record")
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return status.Error(codes.Internal, "failed to commit transaction")
	}

	return nil
}





/*
const result = await prisma.userQuestionInteraction.create({
      data: {
        userId,
        questionId,
        type: interactionType,
        studyMethod: studyMethod,
        strengthScore: scoreResult?.strengthScore || 0,
        metadata: { ...data, strengthScore: scoreResult },
      },
    });
*/


/*
import { prisma } from "@/lib/database";

export const updateQuestionDifficulty = async (
  questionId: string,
  isCorrect: boolean
): Promise<void> => {
  try {
    // Fetch the current correct and incorrect counts for the question
    const question = await prisma.question.findUnique({
      where: { id: questionId },
      select: { correctCount: true, incorrectCount: true },
    });

    if (!question) {
      console.error(`Question with id ${questionId} not found.`);
      return;
    }

    // Initialize counts if they are null
    const currentCorrectCount = question.correctCount ?? 0;
    const currentIncorrectCount = question.incorrectCount ?? 0;

    // Increment either correctCount or incorrectCount based on isCorrect
    const newCorrectCount = isCorrect ? currentCorrectCount + 1 : currentCorrectCount;
    const newIncorrectCount = !isCorrect ? currentIncorrectCount + 1 : currentIncorrectCount;

    // Calculate the new difficulty ratio
    const totalAttempts = newCorrectCount + newIncorrectCount;
    const newDifficultyRatio = totalAttempts > 0 ? newCorrectCount / totalAttempts : 0;

    // Update the question with new counts and difficulty ratio
    await prisma.question.update({
      where: { id: questionId },
      data: {
        correctCount: newCorrectCount,
        incorrectCount: newIncorrectCount,
        difficultyRatio: newDifficultyRatio,
        updatedAt: new Date(),
      },
    });

  } catch (error) {
    console.error(
      `Failed to update difficulty metrics for question with id ${questionId}. ${error.message}`
    );
  }
};
*/


/*
-- public."UserQuestionInteraction" definition

-- Drop table

-- DROP TABLE public."UserQuestionInteraction";

CREATE TABLE public."UserQuestionInteraction" (
	id text NOT NULL,
	"occurredAt" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	"studyMethod" public."StudyMethod" DEFAULT 'None'::"StudyMethod" NOT NULL,
	"type" public."InteractionType" NOT NULL,
	correct bool NULL,
	"strengthScore" float8 DEFAULT 0.0 NOT NULL,
	metadata jsonb NULL,
	"userId" text NOT NULL,
	"questionId" text NOT NULL,
	"createdAt" timestamp(3) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT "UserQuestionInteraction_pkey" PRIMARY KEY (id)
);
CREATE INDEX "UserQuestionInteraction_questionId_createdAt_idx" ON public."UserQuestionInteraction" USING btree ("questionId", "createdAt");
CREATE INDEX "UserQuestionInteraction_questionId_idx" ON public."UserQuestionInteraction" USING btree ("questionId");
CREATE INDEX "UserQuestionInteraction_questionId_userId_idx" ON public."UserQuestionInteraction" USING btree ("questionId", "userId");
CREATE INDEX "UserQuestionInteraction_userId_occurredAt_idx" ON public."UserQuestionInteraction" USING btree ("userId", "occurredAt");
CREATE INDEX idx_interaction_user_question ON public."UserQuestionInteraction" USING btree ("userId", "questionId");


-- public."UserQuestionInteraction" foreign keys

ALTER TABLE public."UserQuestionInteraction" ADD CONSTRAINT "UserQuestionInteraction_questionId_fkey" FOREIGN KEY ("questionId") REFERENCES public."Question"(id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE public."UserQuestionInteraction" ADD CONSTRAINT "UserQuestionInteraction_userId_fkey" FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE ON UPDATE CASCADE;
*/
