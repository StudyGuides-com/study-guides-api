package router

import (
	"context"

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
      "GetTagCount": func(ctx context.Context, params map[string]string) (string, error) {
        return handleTagCount(ctx, store, params)
      },
    },
  }
}

func (r *OperationRouter) Route(ctx context.Context, op string, params map[string]string) (string, error) {
  handler, ok := r.Handlers[op]
  if !ok {
    return "Sorry, I can't help with that.", nil
  }
  return handler(ctx, params)
}

