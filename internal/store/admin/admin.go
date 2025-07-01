package admin

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type AdminStore interface {
	// User retrieves a user by their ID
	User(ctx context.Context, id string) (*sharedpb.User, error)

	// UserByEmail retrieves a user by their email
	UserByEmail(ctx context.Context, email string) (*sharedpb.User, error)

	// Tag retrieves a tag by its ID
	Tag(ctx context.Context, id string) (*sharedpb.Tag, error)

	// TagInfo retrieves a tag by its ID
	TagInfos(ctx context.Context, id string) ([]*sharedpb.TagInfo, error)

	// UpsertTag saves or updates a tag in the database
	UpsertTag(ctx context.Context, tag *sharedpb.Tag) (*sharedpb.Tag, error)

	// UpsertPassage saves or updates a passage in the database
	UpsertPassage(ctx context.Context, passage *sharedpb.Passage) (*sharedpb.Passage, error)

	// UpsertQuestion saves or updates a question in the database
	UpsertQuestion(ctx context.Context, question *sharedpb.Question) (*sharedpb.Question, error)

	// UpsertQuestionTag saves a question tag if it doesn't exist
	UpsertQuestionTag(ctx context.Context, questionTag *sharedpb.QuestionTag) error

	// TagsByContextType retrieves tags by context type
	TagsByContextType(ctx context.Context, contextType sharedpb.ContextType) ([]*sharedpb.Tag, error)

	// TagsForIndexing retrieves tags for indexing
	TagsForIndexing(ctx context.Context, contextType sharedpb.ContextType) ([]*sharedpb.TagIndexResult, error)

	// Dump retrieves the entire study guide hierarchy for a given context type
	Dump(ctx context.Context, contextType sharedpb.ContextType) ([]*sharedpb.Node, error)

	// DumpId retrieves a specific study guide node by its ID
	DumpId(ctx context.Context, id string) (*sharedpb.Node, error)

	// TreeId retrieves the entire study guide hierarchy for a given id
	Tree(ctx context.Context, id string) (*sharedpb.TagNode, error)

	// KillTree kills the tree for a given id
	KillTree(ctx context.Context, id string) ([]string, error)

	// KillUser kills the user for a given email, returns true if deleted, false if not found
	KillUser(ctx context.Context, email string) (bool, error)

	// TagExistsFor checks if a tag exists for a given objectID
	TagExistsFor(ctx context.Context, id string) (bool, error)

	// UserExistsFor checks if a user exists for a given objectID
	UserExistsFor(ctx context.Context, id string) (bool, error)

	// UpdateMetadataFor updates the metadata for a given id
	UpdateMetadataFor(ctx context.Context, id string, metadata *sharedpb.Metadata) error

	// ImportGob imports a gob payload into the database
	ImportGob(ctx context.Context, gobPayload []byte) (bool, error)

	// BuildIndexCache builds the index cache for a given id
	UpdateIndexCache(ctx context.Context, id string) error

	// ClearIndexCache clears the index cache
	ClearIndexCache(ctx context.Context) error
}

func NewSqlAdminStore(ctx context.Context, dbURL string) (*SqlAdminStore, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to connect to postgres: "+err.Error())
	}
	return &SqlAdminStore{db: db}, nil
}
