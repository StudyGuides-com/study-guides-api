package admin

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/utils"
)

type SqlAdminStore struct {
	db *pgxpool.Pool
}

const maxTagDepth = 5 // Maximum depth for tag hierarchy traversal

// NewTag creates a new tag
func NewTag(id string, name string, hash string, tagType sharedpb.TagType, parentTagId *string, contentRating sharedpb.ContentRating, contentDescriptors []sharedpb.ContentDescriptorType, metaTags []string, parserType sharedpb.ParserType, metadata *sharedpb.Metadata) *sharedpb.Tag {
	if contentRating == sharedpb.ContentRating_Unspecified {
		contentRating = sharedpb.ContentRating_RatingPending
	}

	// Convert ParserType to ContextType
	contextType, _ := utils.GetContextTypeForParser(parserType)

	return &sharedpb.Tag{
		Id:                 id,
		Name:               name,
		Description:        &name,
		Hash:               hash,
		Type:               tagType,
		ParentTagId:        parentTagId,
		ContentRating:      contentRating,
		ContentDescriptors: contentDescriptors,
		MetaTags:           metaTags,
		Context:            contextType,
		Public:             true,
		Metadata:           metadata,
	}
}

// NewPassage creates a new passage
func NewPassage(id string, title string, body string, tagId string, metadata *sharedpb.Metadata) *sharedpb.Passage {
	return &sharedpb.Passage{
		Id:       id,
		Title:    title,
		Body:     body,
		TagId:    tagId,
		Metadata: metadata,
	}
}

// NewQuestion creates a new question
func NewQuestion(id string, passageId *string, hash string, questionText string, answerText string, learnMore *string, distractors *[]string, metadata *sharedpb.Metadata) *sharedpb.Question {
	return &sharedpb.Question{
		Id:           id,
		PassageId:    passageId,
		QuestionText: questionText,
		AnswerText:   answerText,
		Hash:         hash,
		LearnMore:    learnMore,
		Distractors:  *distractors,
		Public:       true,
		Version:      1,
		Metadata:     metadata,
	}
}

// NewQuestionTag creates a new question tag
func NewQuestionTag(questionId string, tagId string) *sharedpb.QuestionTag {
	return &sharedpb.QuestionTag{
		QuestionId: questionId,
		TagId:      tagId,
	}
}

// UpsertTag saves or updates a tag in the database
func (s *SqlAdminStore) UpsertTag(ctx context.Context, tag *sharedpb.Tag) (*sharedpb.Tag, error) {
	now := time.Now()
	if tag.CreatedAt == nil {
		tag.CreatedAt = timestamppb.New(now)
	}
	tag.UpdatedAt = timestamppb.New(now)

	if tag.ContentDescriptors == nil {
		tag.ContentDescriptors = []sharedpb.ContentDescriptorType{}
	}
	if tag.MetaTags == nil {
		tag.MetaTags = []string{}
	}

	query := `
		INSERT INTO public."Tag" (
			id, "batchId", hash, name, description, type, context,
			"parentTagId", "contentRating", "contentDescriptors", "metaTags",
			public, "accessCount", metadata, "createdAt", "updatedAt",
			"hasQuestions", "hasChildren", "ownerId"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		)
		ON CONFLICT (hash) DO UPDATE SET
			"batchId" = EXCLUDED."batchId",
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			type = EXCLUDED.type,
			context = EXCLUDED.context,
			"parentTagId" = EXCLUDED."parentTagId",
			"contentRating" = EXCLUDED."contentRating",
			"contentDescriptors" = EXCLUDED."contentDescriptors",
			"metaTags" = EXCLUDED."metaTags",
			public = EXCLUDED.public,
			"accessCount" = EXCLUDED."accessCount",
			metadata = EXCLUDED.metadata,
			"updatedAt" = EXCLUDED."updatedAt",
			"hasQuestions" = EXCLUDED."hasQuestions",
			"hasChildren" = EXCLUDED."hasChildren",
			"ownerId" = EXCLUDED."ownerId"
		RETURNING *
	`

	var updated sharedpb.Tag
	err := pgxscan.Get(ctx, s.db, &updated, query,
		tag.Id, tag.BatchId, tag.Hash, tag.Name, tag.Description, tag.Type, tag.Context,
		tag.ParentTagId, tag.ContentRating, tag.ContentDescriptors, tag.MetaTags,
		tag.Public, tag.AccessCount, tag.Metadata, tag.CreatedAt, tag.UpdatedAt,
		tag.HasQuestions, tag.HasChildren, tag.OwnerId,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to upsert tag")
	}

	return &updated, nil
}

// UpsertPassage saves or updates a passage in the database
func (s *SqlAdminStore) UpsertPassage(ctx context.Context, passage *sharedpb.Passage) (*sharedpb.Passage, error) {
	now := time.Now()
	if passage.CreatedAt == nil {
		passage.CreatedAt = timestamppb.New(now)
	}
	passage.UpdatedAt = timestamppb.New(now)

	query := `
		INSERT INTO passages (
			id, title, body, hash, tagId, metadata, createdAt, updatedAt
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			body = EXCLUDED.body,
			hash = EXCLUDED.hash,
			tagId = EXCLUDED.tagId,
			metadata = EXCLUDED.metadata,
			updatedAt = EXCLUDED.updatedAt
		RETURNING id, title, body, hash, tagId, metadata, createdAt, updatedAt
	`

	var updated sharedpb.Passage
	err := pgxscan.Get(ctx, s.db, &updated, query,
		passage.Id,
		passage.Title,
		passage.Body,
		passage.Hash,
		passage.TagId,
		passage.Metadata,
		passage.CreatedAt,
		passage.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

// UpsertQuestion saves or updates a question in the database
func (s *SqlAdminStore) UpsertQuestion(ctx context.Context, question *sharedpb.Question) (*sharedpb.Question, error) {
	// Set timestamps if not already set
	now := time.Now()
	if question.CreatedAt == nil {
		question.CreatedAt = timestamppb.New(now)
	}
	question.UpdatedAt = timestamppb.New(now)

	query := `
		INSERT INTO public."Question" (
			id, "batchId", "questionText", "answerText", "hash", "learnMore",
			"distractors", "videoUrl", "imageUrl", "version", "public",
			"metadata", "createdAt", "updatedAt", "passageId"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
		ON CONFLICT (hash) DO UPDATE SET
			"batchId" = EXCLUDED."batchId",
			"questionText" = EXCLUDED."questionText",
			"answerText" = EXCLUDED."answerText",
			"learnMore" = EXCLUDED."learnMore",
			"distractors" = EXCLUDED."distractors",
			"videoUrl" = EXCLUDED."videoUrl",
			"imageUrl" = EXCLUDED."imageUrl",
			"version" = EXCLUDED."version",
			"public" = EXCLUDED."public",
			"metadata" = EXCLUDED."metadata",
			"updatedAt" = EXCLUDED."updatedAt",
			"passageId" = EXCLUDED."passageId"
		RETURNING id, "batchId", "questionText", "answerText", "hash", "learnMore",
			"distractors", "videoUrl", "imageUrl", "version", "public",
			"metadata", "createdAt", "updatedAt", "passageId"
	`

	var updated sharedpb.Question
	err := pgxscan.Get(ctx, s.db, &updated, query,
		question.Id,
		question.BatchId,
		question.QuestionText,
		question.AnswerText,
		question.Hash,
		question.LearnMore,
		question.Distractors,
		question.VideoUrl,
		question.ImageUrl,
		question.Version,
		question.Public,
		question.Metadata,
		question.CreatedAt,
		question.UpdatedAt,
		question.PassageId,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to upsert question")
	}

	return &updated, nil
}

// UpsertQuestionTag saves a question tag if it doesn't exist
func (s *SqlAdminStore) UpsertQuestionTag(ctx context.Context, questionTag *sharedpb.QuestionTag) error {
	now := time.Now()
	// Set timestamp if not already set
	if questionTag.CreatedAt == nil {
		questionTag.CreatedAt = timestamppb.New(now)
	}

	query := `
		INSERT INTO public."QuestionTag" (
			"questionId", "tagId", "createdAt"
		) VALUES (
			$1, $2, $3
		)
		ON CONFLICT ("questionId", "tagId") DO NOTHING
	`

	_, err := s.db.Exec(ctx, query,
		questionTag.QuestionId,
		questionTag.TagId,
		questionTag.CreatedAt,
	)

	if err != nil {
		return status.Error(codes.Internal, "failed to upsert question tag")
	}

	return nil
}

// TagsByContextType retrieves tags by context type
func (s *SqlAdminStore) TagsByContextType(ctx context.Context, contextType sharedpb.ContextType) ([]*sharedpb.Tag, error) {
	query := `
		SELECT *  
		FROM public."Tag" 
		WHERE context = $1
	`

	rows, err := s.db.Query(ctx, query, contextType)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get tags by context type")
	}
	defer rows.Close()

	tags := []*sharedpb.Tag{}
	for rows.Next() {
		var tag sharedpb.Tag
		err = rows.Scan(
			&tag.Id,
			&tag.BatchId,
			&tag.Hash,
			&tag.Name,
			&tag.Description,
			&tag.Type,
			&tag.ParentTagId,
			&tag.ContentRating,
			&tag.ContentDescriptors,
			&tag.MetaTags,
			&tag.Public,
			&tag.Metadata,
			&tag.CreatedAt,
			&tag.UpdatedAt,
			&tag.AccessCount,
			&tag.Context,
			&tag.OwnerId,
			&tag.HasQuestions,
			&tag.HasChildren,
		)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to scan tag")
		}
		tags = append(tags, &tag)
	}

	return tags, nil
}

func (s *SqlAdminStore) TagsForIndexing(ctx context.Context, contextType sharedpb.ContextType) ([]*sharedpb.TagIndexResult, error) {
	query := `
		SELECT t.*
		FROM public."Tag" t
		LEFT JOIN _index_cache ic ON ic.id = t.id
		WHERE t.context = $1 AND ic.id IS NULL
		ORDER BY t."updatedAt" DESC
		LIMIT 5000
	`

	var tags []*sharedpb.Tag
	err := pgxscan.Select(ctx, s.db, &tags, query, contextType)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get tags for indexing")
	}

	results := []*sharedpb.TagIndexResult{}
	for _, tag := range tags {
		// Get TagInfo for this tag
		tagInfos, err := s.TagInfos(ctx, tag.Id)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to get tag infos")
		}

		result := &sharedpb.TagIndexResult{
			Tag:      tag,
			TagInfos: tagInfos,
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *SqlAdminStore) Dump(ctx context.Context, contextType sharedpb.ContextType) ([]*sharedpb.Node, error) {
	query := `
		WITH RECURSIVE ancestry AS (
  -- Get leaf tags (hasQuestions = true, context = $1)
  SELECT
    t.id AS leaf_id,
    NULL::jsonb AS node,
    t."parentTagId" AS parent_id,
    1 AS level  -- root will end up with highest level
  FROM public."Tag" t
  WHERE t."hasQuestions" = true
    AND t."context" = $1
	AND (
    t."metadata" IS NULL OR NOT (t."metadata" ? 'bmr_export')
  )

  UNION ALL

  -- Walk up the parent chain, building parent nodes
  SELECT
    a.leaf_id,
    jsonb_build_object(
      'id', pt.id,
      'name', pt.name,
      'type', pt.type,
      'level', a.level
    ) AS node,
    pt."parentTagId",
    a.level + 1
  FROM ancestry a
  JOIN public."Tag" pt ON pt.id = a.parent_id
),
-- Get each leaf tag's ancestry nodes
parent_nodes AS (
  SELECT
    leaf_id,
    ARRAY_AGG(node ORDER BY level DESC) FILTER (WHERE node IS NOT NULL) AS nodes
  FROM ancestry
  GROUP BY leaf_id
),
-- Build the nested tree
tree_built AS (
  SELECT
    leaf_id,
    build_tree(nodes::jsonb[]) AS parent
  FROM parent_nodes
),
-- Get questions for each leaf tag
questions_agg AS (
  SELECT
    qt."tagId" AS tag_id,
    jsonb_agg(jsonb_build_object(
      'id', q.id,
      'questionText', q."questionText",
      'answerText', q."answerText",
      'learnMore', q."learnMore",
      'distractors', q."distractors"
    ) ORDER BY q."createdAt") AS questions
  FROM public."QuestionTag" qt
  JOIN public."Question" q ON q.id = qt."questionId"
  GROUP BY qt."tagId"
),
-- Get leaf tag info (not part of the parent chain!)
leaf_tags AS (
  SELECT
    t.id,
    t.name,
    t.type
  FROM public."Tag" t
  WHERE t."hasQuestions" = true
    AND t."context" = $1
	AND (
    t."metadata" IS NULL OR NOT (t."metadata" ? 'bmr_export')
  )
)
-- Final output
SELECT jsonb_build_object(
  'id', lt.id,
  'name', lt.name,
  'type', lt.type,
  'level', 0,
  'questions', COALESCE(qa.questions, '[]'::jsonb),
  'parent', tb.parent
) AS tag_json
FROM leaf_tags lt
LEFT JOIN tree_built tb ON tb.leaf_id = lt.id
LEFT JOIN questions_agg qa ON qa.tag_id = lt.id
ORDER BY lt.id;
	`

	rows, err := s.db.Query(ctx, query, contextType)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to execute dump query")
	}
	defer rows.Close()

	var nodes []*sharedpb.Node
	for rows.Next() {
		var tagJSON []byte
		if err := rows.Scan(&tagJSON); err != nil {
			return nil, status.Error(codes.Internal, "failed to scan row")
		}

		var node sharedpb.Node
		if err := json.Unmarshal(tagJSON, &node); err != nil {
			return nil, status.Error(codes.Internal, "failed to unmarshal node")
		}

		nodes = append(nodes, &node)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Internal, "error iterating rows")
	}

	return nodes, nil
}

func (s *SqlAdminStore) DumpId(ctx context.Context, id string) (*sharedpb.Node, error) {
	query := `
		WITH RECURSIVE ancestry AS (
			-- Get the specific tag by ID
			SELECT
				t.id AS leaf_id,
				NULL::jsonb AS node,
				t."parentTagId" AS parent_id,
				1 AS level
			FROM public."Tag" t
			WHERE t.id = $1

			UNION ALL

			-- Walk up the parent chain, building parent nodes
			SELECT
				a.leaf_id,
				jsonb_build_object(
					'id', pt.id,
					'name', pt.name,
					'type', pt.type,
					'level', a.level
				) AS node,
				pt."parentTagId",
				a.level + 1
			FROM ancestry a
			JOIN public."Tag" pt ON pt.id = a.parent_id
		),
		-- Get the tag's ancestry nodes
		parent_nodes AS (
			SELECT
				leaf_id,
				ARRAY_AGG(node ORDER BY level DESC) FILTER (WHERE node IS NOT NULL) AS nodes
			FROM ancestry
			GROUP BY leaf_id
		),
		-- Build the nested tree
		tree_built AS (
			SELECT
				leaf_id,
				build_tree(nodes::jsonb[]) AS parent
			FROM parent_nodes
		),
		-- Get questions for the tag
		questions_agg AS (
			SELECT
				qt."tagId" AS tag_id,
				jsonb_agg(jsonb_build_object(
					'id', q.id,
					'questionText', q."questionText",
					'answerText', q."answerText",
					'learnMore', q."learnMore",
					'distractors', q."distractors"
				) ORDER BY q."createdAt") AS questions
			FROM public."QuestionTag" qt
			JOIN public."Question" q ON q.id = qt."questionId"
			WHERE qt."tagId" = $1
			GROUP BY qt."tagId"
		),
		-- Get tag info
		tag_info AS (
			SELECT
				t.id,
				t.name,
				t.type
			FROM public."Tag" t
			WHERE t.id = $1
		)
		-- Final output
		SELECT jsonb_build_object(
			'id', ti.id,
			'name', ti.name,
			'type', ti.type,
			'level', 0,
			'questions', COALESCE(qa.questions, '[]'::jsonb),
			'parent', tb.parent
		) AS tag_json
		FROM tag_info ti
		LEFT JOIN tree_built tb ON tb.leaf_id = ti.id
		LEFT JOIN questions_agg qa ON qa.tag_id = ti.id;
	`

	var tagJSON []byte
	err := s.db.QueryRow(ctx, query, id).Scan(&tagJSON)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to execute dump query")
	}

	var node sharedpb.Node
	if err := json.Unmarshal(tagJSON, &node); err != nil {
		return nil, status.Error(codes.Internal, "failed to unmarshal node")
	}

	return &node, nil
}

func (s *SqlAdminStore) Tree(ctx context.Context, id string) (*sharedpb.TagNode, error) {
	query := `
		WITH RECURSIVE tag_tree AS (
			SELECT
				t.id,
				t."parentTagId",
				t.name,
				t.type,
				t.public,
				t.context,
				t.description,
				t."hasChildren",
				t."hasQuestions",
				0 AS level
			FROM public."Tag" t
			WHERE t.id = $1

			UNION ALL

			SELECT
				c.id,
				c."parentTagId",
				c.name,
				c.type,
				c.public,
				c.context,
				c.description,
				c."hasChildren",
				c."hasQuestions",
				tt.level + 1 AS level
			FROM public."Tag" c
			JOIN tag_tree tt ON c."parentTagId" = tt.id
		)
		SELECT
			id,
			COALESCE("parentTagId", '') AS "parent_tag_id",
			name,
			CASE
				WHEN type = 'Category' THEN 0
				WHEN type = 'SubCategory' THEN 1
				WHEN type = 'University' THEN 2
				WHEN type = 'Region' THEN 3
				WHEN type = 'Department' THEN 4
				WHEN type = 'Course' THEN 5
				WHEN type = 'Topic' THEN 6
				WHEN type = 'UserStudyGuide' THEN 7
				WHEN type = 'UserContent' THEN 8
				WHEN type = 'UserFolder' THEN 9
				WHEN type = 'UserTopic' THEN 10
				WHEN type = 'Organization' THEN 11
				WHEN type = 'Certifying_Agency' THEN 12
				WHEN type = 'Certification' THEN 13
				WHEN type = 'Module' THEN 14
				WHEN type = 'Domain' THEN 15
				WHEN type = 'Entrance_Exam' THEN 16
				WHEN type = 'AP_Exam' THEN 17
				WHEN type = 'Branch' THEN 18
				WHEN type = 'Instruction_Type' THEN 19
				WHEN type = 'Instruction_Group' THEN 20
				WHEN type = 'Instruction' THEN 21
				WHEN type = 'Chapter' THEN 22
				WHEN type = 'Section' THEN 23
				WHEN type = 'Part' THEN 24
				WHEN type = 'Volume' THEN 25
				WHEN type = 'Range' THEN 26
				ELSE 0
			END AS type,
			public,
			CASE
				WHEN context = 'Colleges' THEN 0
				WHEN context = 'Certifications' THEN 1
				WHEN context = 'EntranceExams' THEN 2
				WHEN context = 'APExams' THEN 3
				WHEN context = 'UserGeneratedContent' THEN 4
				WHEN context = 'DoD' THEN 5
				WHEN context = 'Encyclopedia' THEN 6
				ELSE 0
			END AS context,
			COALESCE(description, '') AS description,
			"hasChildren" AS "has_children",
			"hasQuestions" AS "has_questions",
			level
		FROM tag_tree
		ORDER BY level, id;
	`

	var tagRows []*sharedpb.TagRow
	err := pgxscan.Select(ctx, s.db, &tagRows, query, id)
	if err != nil {
		log.Printf("Failed to scan tag rows for id %s: %v", id, err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to scan tag rows: %v", err))
	}

	// Build a map of nodes by ID for quick lookup
	nodeMap := make(map[string]*sharedpb.TagNode)
	var rootNode *sharedpb.TagNode

	// First pass: Create all nodes
	for _, row := range tagRows {
		node := &sharedpb.TagNode{
			TagRow:   row,
			Children: []*sharedpb.TagNode{},
		}
		nodeMap[row.Id] = node
		if row.Level == 0 {
			rootNode = node
		}
	}

	// Second pass: Build parent-child relationships
	for _, row := range tagRows {
		if row.ParentTagId != "" {
			if parentNode, exists := nodeMap[row.ParentTagId]; exists {
				parentNode.Children = append(parentNode.Children, nodeMap[row.Id])
			}
		}
	}

	if rootNode == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("tag with id '%s' not found or has no accessible tree structure", id))
	}

	return rootNode, nil
}

func (s *SqlAdminStore) KillTree(ctx context.Context, id string) ([]string, error) {
	tree, err := s.Tree(ctx, id)
	if err != nil {
		return nil, err
	}

	// Collect all IDs from the tree
	ids := collectNodeIDs(tree)

	// Recursively delete the tree starting from leaf nodes
	if err := s.deleteTreeRecursively(ctx, tree); err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *SqlAdminStore) UserByEmail(ctx context.Context, email string) (*sharedpb.User, error) {
	query := `
		SELECT id, name, email, image, "emailVerified", "stripeCustomerId", "gamerTag"
		FROM public."User" WHERE email = $1
	`

	var user sharedpb.User
	err := pgxscan.Get(ctx, s.db, &user, query, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("user not found with email: %s", email))
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &user, nil
}

// collectNodeIDs recursively collects all node IDs from a tree
func collectNodeIDs(node *sharedpb.TagNode) []string {
	ids := []string{node.TagRow.Id}
	for _, child := range node.Children {
		ids = append(ids, collectNodeIDs(child)...)
	}
	return ids
}

// deleteTreeRecursively recursively deletes a tree of tags, starting from leaf nodes
func (s *SqlAdminStore) deleteTreeRecursively(ctx context.Context, node *sharedpb.TagNode) error {
	// First delete all children recursively
	for _, child := range node.Children {
		if err := s.deleteTreeRecursively(ctx, child); err != nil {
			return err
		}
	}

	// Then delete this node and its references
	if err := s.deleteTagAndReferences(ctx, node.TagRow.Id); err != nil {
		return fmt.Errorf("failed to delete tag %s (%s): %w", node.TagRow.Id, node.TagRow.Name, err)
	}

	return nil
}

// deleteTagAndReferences deletes a tag and all its references from related tables
func (s *SqlAdminStore) deleteTagAndReferences(ctx context.Context, tagID string) error {
	query := `
		WITH deleted_question_tags AS (
			DELETE FROM public."QuestionTag" WHERE "tagId" = $1
		),
		deleted_passages AS (
			DELETE FROM public."Passage" WHERE "tagId" = $1
		),
		deleted_ratings AS (
			DELETE FROM public."UserTagRating" WHERE "tagId" = $1
		),
		deleted_test_questions AS (
			DELETE FROM public."TestQuestion"
			WHERE "sessionId" IN (SELECT id FROM public."TestSession" WHERE "tagId" = $1)
		),
		deleted_test_sessions AS (
			DELETE FROM public."TestSession" WHERE "tagId" = $1
		),
		deleted_reports AS (
			DELETE FROM public."UserTagReport" WHERE "tagId" = $1
		),
		deleted_recent_tags AS (
			DELETE FROM public."UserTagRecent" WHERE "tagId" = $1
		),
		deleted_favorite_tags AS (
			DELETE FROM public."UserTagFavorite" WHERE "tagId" = $1
		),
		deleted_topic_progress AS (
			DELETE FROM public."UserTopicProgress" WHERE "topicId" = $1
		),
		deleted_survival_questions AS (
			DELETE FROM public."SurvivalQuestion"
			WHERE "sessionId" IN (SELECT id FROM public."SurvivalSession" WHERE "tagId" = $1)
		),
		deleted_survival_sessions AS (
			DELETE FROM public."SurvivalSession" WHERE "tagId" = $1
		),
		deleted_tag_access AS (
			DELETE FROM public."TagAccess" WHERE "tagId" = $1
		),
		deleted_tag_invites AS (
			DELETE FROM public."TagInvite" WHERE "tagId" = $1
		),
		deleted_algolia_records AS (
			DELETE FROM public."AlgoliaRecord" WHERE "id" = $1
		)
		DELETE FROM public."Tag" WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, tagID)
	if err != nil {
		log.Printf("Database error deleting tag %s: %v", tagID, err)
		return status.Error(codes.Internal, fmt.Sprintf("database error deleting tag %s: %v", tagID, err))
	}

	return nil
}

func (s *SqlAdminStore) TagExistsFor(ctx context.Context, id string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM public."Tag" WHERE id = $1
		)
	`

	var exists bool
	err := s.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, status.Error(codes.Internal, "failed to check if tag exists")
	}

	return exists, nil
}

func (s *SqlAdminStore) UserExistsFor(ctx context.Context, id string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM public."User" WHERE id = $1
		)
	`

	var exists bool
	err := s.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, status.Error(codes.Internal, "failed to check if user exists")
	}

	return exists, nil
}

// UpdateMetadataFor updates the metadata for a given tag id
func (s *SqlAdminStore) UpdateMetadataFor(ctx context.Context, id string, metadata *sharedpb.Metadata) error {
	// First, get the current metadata
	query := `
		SELECT metadata 
		FROM public."Tag" 
		WHERE id = $1
	`

	var currentMetadata []byte
	err := s.db.QueryRow(ctx, query, id).Scan(&currentMetadata)
	if err != nil {
		return status.Error(codes.Internal, "failed to get current metadata")
	}

	// Parse current metadata
	var currentMap map[string]interface{}
	if len(currentMetadata) > 0 {
		if err := json.Unmarshal(currentMetadata, &currentMap); err != nil {
			return status.Error(codes.Internal, "failed to parse current metadata")
		}
	} else {
		currentMap = make(map[string]interface{})
	}

	// Merge new metadata with current metadata
	for k, v := range metadata.Metadata {
		currentMap[k] = v
	}

	// Update metadata in Tag table
	updateQuery := `
		UPDATE public."Tag" 
		SET metadata = $1, "updatedAt" = NOW() 
		WHERE id = $2
	`

	// Convert merged metadata back to JSON
	mergedMetadata, err := json.Marshal(currentMap)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal merged metadata")
	}

	_, err = s.db.Exec(ctx, updateQuery, mergedMetadata, id)
	if err != nil {
		return status.Error(codes.Internal, "failed to update metadata")
	}

	return nil
}

func (s *SqlAdminStore) ImportGob(ctx context.Context, gobPayload []byte) (bool, error) {
	// convert the gobPayload into a guideData
	var guideData sharedpb.GuideData
	err := gob.NewDecoder(bytes.NewReader(gobPayload)).Decode(&guideData)
	if err != nil {
		return false, status.Error(codes.Internal, "failed to decode gob payload")
	}

	for _, section := range guideData.Sections {
		topicId, err := s.ImportAncestry(ctx, section.Ancestor, guideData.ParserType, guideData.Title, *section)
		if err != nil {
			return false, status.Error(codes.Internal, fmt.Sprintf("failed to import ancestry for section %s", section.Title))
		}

		for _, passage := range section.Passages {
			metadata := &sharedpb.Metadata{
				Metadata: map[string]string{
					"parserType":   guideData.ParserType.String(),
					"ts":           time.Now().UTC().Format(time.RFC3339),
					"guideTitle":   guideData.Title,
					"sectionTitle": section.Title,
				},
			}
			p := NewPassage(utils.GetCUID(), passage.Title, passage.Content, topicId, metadata)
			_, err := s.UpsertPassage(ctx, p)
			if err != nil {
				return false, status.Error(codes.Internal, fmt.Sprintf("failed to upsert passage for section %s", section.Title))
			}

			for _, prompt := range passage.Prompts {
				metadata := &sharedpb.Metadata{
					Metadata: map[string]string{
						"parserType":   guideData.ParserType.String(),
						"ts":           time.Now().UTC().Format(time.RFC3339),
						"guideTitle":   guideData.Title,
						"sectionTitle": section.Title,
					},
				}
				question := NewQuestion(utils.GetCUID(), &p.Id, prompt.Hash, prompt.Question, prompt.Answer, &prompt.LearnMore, &prompt.Distractors, metadata)
				updatedQuestion, err := s.UpsertQuestion(ctx, question)
				if err != nil {
					return false, status.Error(codes.Internal, fmt.Sprintf("failed to upsert question for section %s", section.Title))
				}

				questionTag := NewQuestionTag(updatedQuestion.Id, topicId)
				err = s.UpsertQuestionTag(ctx, questionTag)
				if err != nil {
					return false, status.Error(codes.Internal, fmt.Sprintf("failed to upsert question tag for section %s", section.Title))
				}
			}
		}

		for _, prompt := range section.Prompts {
			metadata := &sharedpb.Metadata{
				Metadata: map[string]string{
					"parserType":   guideData.ParserType.String(),
					"ts":           time.Now().UTC().Format(time.RFC3339),
					"guideTitle":   guideData.Title,
					"sectionTitle": section.Title,
				},
			}
			question := NewQuestion(utils.GetCUID(), nil, prompt.Hash, prompt.Question, prompt.Answer, &prompt.LearnMore, &prompt.Distractors, metadata)
			updatedQuestion, err := s.UpsertQuestion(ctx, question)
			if err != nil {
				return false, status.Error(codes.Internal, fmt.Sprintf("failed to upsert question for section %s", section.Title))
			}

			questionTag := NewQuestionTag(updatedQuestion.Id, topicId)
			err = s.UpsertQuestionTag(ctx, questionTag)
			if err != nil {
				return false, status.Error(codes.Internal, fmt.Sprintf("failed to upsert question tag for section %s", section.Title))
			}
		}
	}

	return true, nil
}

func (s *SqlAdminStore) ancestorList(ancestor *sharedpb.Ancestor) ([]*sharedpb.Ancestor, error) {
	if ancestor == nil {
		return nil, nil
	}

	// Step 1: Collect all ancestors into a slice, starting from the youngest
	var ancestors []*sharedpb.Ancestor
	current := ancestor
	for current != nil {
		ancestors = append(ancestors, current)
		current = current.NextAncestor
	}

	return ancestors, nil
}

func (s *SqlAdminStore) ImportAncestry(ctx context.Context, ancestor *sharedpb.Ancestor, parserType sharedpb.ParserType, guideTitle string, section sharedpb.SectionData) (string, error) {
	if ancestor == nil {
		return "", status.Error(codes.Internal, "ancestor is nil")
	}

	// Step 1: Collect all ancestors into a slice, starting from the youngest
	ancestors, err := s.ancestorList(ancestor)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to get ancestors")
	}

	// Step 2: Reverse the slice to get oldest first
	for i, j := 0, len(ancestors)-1; i < j; i, j = i+1, j-1 {
		ancestors[i], ancestors[j] = ancestors[j], ancestors[i]
	}

	metadata := &sharedpb.Metadata{
		Metadata: map[string]string{
			"parserType":   parserType.String(),
			"ts":           time.Now().UTC().Format(time.RFC3339),
			"guideTitle":   guideTitle,
			"sectionTitle": section.Title,
		},
	}

	// Step 3: Process ancestors from oldest to newest
	var parentID *string
	for i, anc := range ancestors {
		// Create a new tag for this ancestor
		var tag *sharedpb.Tag
		if i == len(ancestors)-1 {
			tag = NewTag(
				utils.GetCUID(),
				anc.Name,
				anc.Hash,
				anc.TagType,
				parentID,
				section.ContentRating,
				section.ContentDescriptors,
				section.MetaTags,
				parserType,
				metadata,
			)
		} else {
			tag = NewTag(
				utils.GetCUID(),
				anc.Name,
				anc.Hash,
				anc.TagType,
				parentID,
				anc.ContentRating,
				nil,
				nil,
				parserType,
				metadata,
			)
		}

		// Upsert the tag and get the result
		updatedTag, err := s.UpsertTag(ctx, tag)
		if err != nil {
			return "", status.Error(codes.Internal, fmt.Sprintf("failed to upsert tag for ancestor %s", anc.Name))
		}

		// Set this tag's ID as the parent ID for the next iteration
		parentID = &updatedTag.Id
	}

	return *parentID, nil
}

// Tag retrieves a tag by its ID
func (s *SqlAdminStore) Tag(ctx context.Context, id string) (*sharedpb.Tag, error) {
	query := `
		SELECT * FROM public."Tag" WHERE id = $1
	`

	var tag sharedpb.Tag
	err := pgxscan.Get(ctx, s.db, &tag, query, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get tag")
	}

	return &tag, nil
}

func (s *SqlAdminStore) TagInfos(ctx context.Context, id string) ([]*sharedpb.TagInfo, error) {
	// Cache to store results
	tagCache := make(map[string][]*sharedpb.TagInfo)

	// Start the recursive process
	visitedTags := make(map[string]bool)
	return s.getParentTags(ctx, id, visitedTags, 0, tagCache)
}

// getParentTags recursively retrieves parent tags with caching
func (s *SqlAdminStore) getParentTags(ctx context.Context, tagID string, visitedTags map[string]bool, depth int, tagCache map[string][]*sharedpb.TagInfo) ([]*sharedpb.TagInfo, error) {
	// Check the cache
	if cached, exists := tagCache[tagID]; exists {
		return cached, nil
	}

	// Check for cyclic relationships
	if visitedTags[tagID] {
		return nil, nil
	}

	// Check for maximum recursion depth
	if depth > maxTagDepth {
		return nil, nil
	}

	// Add the current tag to the visited set
	visitedTags[tagID] = true

	// Fetch the tag from the database
	query := `
		SELECT 
			id, name, type, "parentTagId", "hasQuestions", "hasChildren"
		FROM public."Tag"
		WHERE id = $1 AND public = true
	`

	var tag sharedpb.TagInfo
	err := pgxscan.Get(ctx, s.db, &tag, query, tagID)
	if err != nil {
		// If no rows were found, return empty result
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, status.Error(codes.Internal, "failed to get tag")
	}

	// If there's a parent tag, recursively fetch it
	var parentTags []*sharedpb.TagInfo
	if tag.ParentTagId != "" {
		var err error
		parentTags, err = s.getParentTags(ctx, tag.ParentTagId, visitedTags, depth+1, tagCache)
		if err != nil {
			return nil, err
		}
	}

	// Combine parent tags with current tag
	result := append(parentTags, &tag)

	// Store the result in the cache
	tagCache[tagID] = result

	return result, nil
}

func (s *SqlAdminStore) User(ctx context.Context, id string) (*sharedpb.User, error) {
	query := `
		SELECT id, name, email, image, "emailVerified", "stripeCustomerId", "gamerTag"
		FROM public."User" WHERE id = $1
	`

	var user sharedpb.User
	err := pgxscan.Get(ctx, s.db, &user, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("user not found with id: %s", id))
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &user, nil
}

func (s *SqlAdminStore) UpdateIndexCache(ctx context.Context, id string) error {
	query := `
		INSERT INTO public._index_cache (id, ts) 
		VALUES ($1, now())
		ON CONFLICT (id) DO UPDATE 
		SET ts = now()
	`

	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return status.Error(codes.Internal, "failed to update index cache")
	}

	return nil
}

func (s *SqlAdminStore) ClearIndexCache(ctx context.Context) error {
	query := `
		DELETE FROM public._index_cache
	`

	_, err := s.db.Exec(ctx, query)
	if err != nil {
		return status.Error(codes.Internal, "failed to clear index cache")
	}

	return nil
}

func (s *SqlAdminStore) KillUser(ctx context.Context, email string) (bool, error) {
	// First check if the user exists
	_, err := s.UserByEmail(ctx, email)
	if err != nil {
		// If user not found, return false, nil
		if errors.Is(err, pgx.ErrNoRows) || (err.Error() == "rpc error: code = NotFound desc = user not found with email: "+email) {
			return false, nil
		}
		return false, err
	}

	// User exists, proceed with deletion
	query := `
		SELECT delete_user_data_by_email($1)
	`

	_, err = s.db.Exec(ctx, query, email)
	if err != nil {
		return false, status.Error(codes.Internal, "failed to kill user")
	}

	return true, nil
}
