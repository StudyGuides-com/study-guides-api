package kpi

import (
	"github.com/studyguides-com/study-guides-api/internal/repository"
)

// GetResourceSchema returns the MCP schema for KPI operations
func GetResourceSchema() repository.ResourceSchema {
	return repository.ResourceSchema{
		EntityType: KPIExecution{},
		FilterType: KPIFilter{},
		UpdateType: KPIUpdate{},
	}
}

// GetToolDescriptions provides detailed descriptions for AI tool generation
func GetToolDescriptions() map[string]string {
	return map[string]string{
		"kpi_find": `Run KPI calculations or check status of running calculations.
		
Examples:
- "run all KPIs" → {"filter": {"run_all": true}}
- "calculate monthly interactions" → {"filter": {"group": "MonthlyInteractions"}}
- "update tags statistics" → {"filter": {"group": "Tags"}}
- "check running KPIs" → {"filter": {"status": "running"}}
- "run user stats" → {"filter": {"group": "Users"}}

Available KPI groups:
- MonthlyInteractions: Calculate time-based interaction statistics
- Tags: Update tag-related statistics
- TagTypes: Update tag type statistics
- Reports: Update report statistics
- Topics: Update topic statistics
- MissingData: Identify and update missing data statistics
- Ratings: Update rating statistics
- Questions: Update question statistics
- Users: Update user statistics
- UserContent: Update user content statistics
- Contacts: Update contact statistics

Note: These calculations run in the background and may take several minutes to complete.`,

		"kpi_findById": `Check the status of a specific KPI execution by its ID.
Returns current status, duration, and any errors.`,

		"kpi_create": `Start a new KPI calculation.
Specify the group to calculate: MonthlyInteractions, Tags, TagTypes, Reports, Topics, MissingData, Ratings, Questions, Users, UserContent, or Contacts.
The calculation runs in the background and returns immediately with an execution ID.`,

		"kpi_count": `Count KPI executions.
Use with filter to count running executions: {"filter": {"status": "running"}}`,

		"kpi_delete": `Cancel a running KPI execution by its ID.
Only running executions can be cancelled.`,
	}
}

// GetAIPromptAdditions returns additional context for the AI system prompt
func GetAIPromptAdditions() string {
	return `
KPI Operations:
- KPIs are long-running statistical calculations that execute in the background
- When user asks to "run KPIs", "update stats", "calculate metrics", use kpi_find with run_all:true
- When user mentions specific stats like "user stats", "tag stats", map to the appropriate group
- Always inform user that KPI calculations run in background and may take several minutes
- Suggest checking status with "check running KPIs" after starting calculations
- Monthly interactions uses a different procedure (calculate_time_stats) than others (update_calculated_stats)

KPI Group Mappings:
- "monthly stats", "interactions" → MonthlyInteractions
- "tag stats", "tag metrics" → Tags
- "tag type stats" → TagTypes
- "report stats", "reports" → Reports
- "topic stats", "topics" → Topics
- "missing data", "data quality" → MissingData
- "ratings", "rating stats" → Ratings
- "question stats", "questions" → Questions
- "user stats", "user metrics" → Users
- "user content", "content stats" → UserContent
- "contact stats", "contacts" → Contacts
`
}