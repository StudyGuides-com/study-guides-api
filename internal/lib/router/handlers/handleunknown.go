package handlers

import (
	"context"

	"github.com/studyguides-com/study-guides-api/internal/store"
)

func HandleUnknown(ctx context.Context, store store.Store, params map[string]string) (string, error) {
	return "I'm not sure how to help with that request. Could you please rephrase or ask about something else?", nil
}
