package router

import (
	"context"

	"github.com/studyguides-com/study-guides-api/internal/lib/router/handlers"
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
				return handlers.HandleTagCount(ctx, store, params)
			},
			string(tools.ToolNameListTags): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleListTags(ctx, store, params)
			},
			string(tools.ToolNameListRootTags): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleListRootTags(ctx, store, params)
			},
			string(tools.ToolNameGetTag): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleGetTag(ctx, store, params)
			},
			string(tools.ToolNameUniqueTagTypes): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleUniqueTagTypes(ctx, store, params)
			},
			string(tools.ToolNameUniqueContextTypes): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleUniqueContextTypes(ctx, store, params)
			},
			string(tools.ToolNameUserCount): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleUserCount(ctx, store, params)
			},
			string(tools.ToolNameGetUser): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleGetUser(ctx, store, params)
			},
			string(tools.ToolNameDeploy): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleDeploy(ctx, store, params)
			},
			string(tools.ToolNameRollback): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleRollback(ctx, store, params)
			},
			string(tools.ToolNameListDeployments): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleListDeployments(ctx, store, params)
			},
			string(tools.ToolNameGetDeploymentStatus): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleGetDeploymentStatus(ctx, store, params)
			},
			string(tools.ToolNameUnknown): func(ctx context.Context, params map[string]string) (string, error) {
				return handlers.HandleUnknown(ctx, store, params)
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
