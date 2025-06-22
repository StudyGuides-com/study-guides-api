package router

import (
	"context"
  "fmt"

	"github.com/studyguides-com/study-guides-api/internal/store"
)


func handleTagCount(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	count, err := store.TagStore().CountTags(ctx, params)
	if err != nil {
		return "", err
	}
  
	return fmt.Sprintf("You have %d tags matching the criteria.", count), nil
  }
  
