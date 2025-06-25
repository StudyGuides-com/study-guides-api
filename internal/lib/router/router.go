package router

import (
	"context"

	"github.com/studyguides-com/study-guides-api/internal/lib/tools"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

type Router interface {
	Route(ctx context.Context, op string, params map[string]string) (string, error)
}

type OperationRouter struct {
	Handlers map[string]OperationHandler
}

type OperationHandler func(ctx context.Context, params map[string]string) (string, error)

func NewRouter(store store.Store) *OperationRouter {
	return &OperationRouter{
		Handlers: map[string]OperationHandler{
			string(tools.ToolNameTagCount): func(ctx context.Context, params map[string]string) (string, error) {
				return handleTagCount(ctx, store, params)
			},
			string(tools.ToolNameListTags): func(ctx context.Context, params map[string]string) (string, error) {
				return handleListTags(ctx, store, params)
			},
			string(tools.ToolNameListRootTags): func(ctx context.Context, params map[string]string) (string, error) {
				return handleListRootTags(ctx, store, params)
			},
			string(tools.ToolNameUniqueTagTypes): func(ctx context.Context, params map[string]string) (string, error) {
				return handleUniqueTagTypes(ctx, store, params)
			},
			string(tools.ToolNameUnknown): func(ctx context.Context, params map[string]string) (string, error) {
				return handleUnknown(ctx, store, params)
			},
		},
	}
}

func (r *OperationRouter) Route(ctx context.Context, op string, params map[string]string) (string, error) {
	unknownHandler := r.Handlers[string(tools.ToolNameUnknown)]
	handler, ok := r.Handlers[op]
	if !ok {
		return unknownHandler(ctx, params)
	}
	return handler(ctx, params)
}
