package handlers

import (
	"context"
	"fmt"

	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleUniqueContextTypes(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	contextTypes, err := store.TagStore().UniqueContextTypes(ctx)
	if err != nil {
		return "", err
	}

	if len(contextTypes) == 0 {
		return "No context types found.", nil
	}

	response := "Available context types:\n"
	for _, contextType := range contextTypes {
		response += fmt.Sprintf("- %s\n", contextType)
	}

	return response, nil
}
