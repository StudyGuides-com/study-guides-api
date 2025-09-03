package indexing

import (
	"github.com/studyguides-com/study-guides-api/internal/repository"
)

// GetResourceSchema returns the schema for indexing operations
func GetResourceSchema() repository.ResourceSchema {
	return repository.ResourceSchema{
		EntityType: IndexingExecution{},
		FilterType: IndexingFilter{},
		UpdateType: IndexingUpdate{},
	}
}

// GetToolDescriptions provides detailed descriptions for AI tool generation
func GetToolDescriptions() map[string]string {
	return map[string]string{
		"indexing_find": `Trigger indexing operations or check status of running indexing jobs.
		
Examples:
- "reindex tags" → {"filter": {"triggerReindex": true, "objectType": "Tag"}}
- "force reindex tags" → {"filter": {"triggerReindex": true, "objectType": "Tag", "force": true}}
- "sync tags to algolia" → {"filter": {"triggerReindex": true, "objectType": "Tag"}}
- "index all tags" → {"filter": {"triggerReindex": true, "objectType": "Tag"}}
- "check indexing status" → {"filter": {"status": "running"}}
- "check indexing jobs" → {"filter": {}}

Available object types:
- Tag: Index tag data to Algolia search

Options:
- triggerReindex: Set to true to start a new indexing job
- force: Set to true to reindex even if content hasn't changed
- objectType: Specify "Tag" to index tags (required for triggering)
- status: Filter by job status ("running", "complete", "failed")

Note: Indexing runs in the background and may take several minutes to complete.`,

		"indexing_findById": `Check the status of a specific indexing job by its ID.
Returns current status, progress, duration, and any errors.`,

		"indexing_count": `Count indexing jobs.
Use with filter to count running jobs: {"filter": {"status": "running"}}`,
	}
}

// GetAIPromptAdditions returns additional context for the AI system prompt
func GetAIPromptAdditions() string {
	return `
Indexing Operations:
- When user asks to "reindex", "index", "sync to algolia", use indexing_find with triggerReindex:true
- Always specify objectType:"Tag" when triggering indexing
- Use force:true when user mentions "force" or "even if unchanged"
- Indexing jobs run in background and may take several minutes
- Suggest checking status after starting indexing operations

Indexing Trigger Phrases:
- "reindex tags", "index tags", "sync tags" → triggerReindex:true, objectType:"Tag"
- "force reindex", "reindex all" → triggerReindex:true, objectType:"Tag", force:true
- "algolia sync", "sync to search" → triggerReindex:true, objectType:"Tag"
- "check indexing", "indexing status" → status filter or empty filter
`
}