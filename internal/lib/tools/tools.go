package tools

import (
	"github.com/sashabaranov/go-openai"
)

type ToolNames string

const (
	ToolNameTagCount            ToolNames = "TagCount"
	ToolNameListTags            ToolNames = "ListTags"
	ToolNameListRootTags        ToolNames = "ListRootTags"
	ToolNameGetTag              ToolNames = "GetTag"
	ToolNameUniqueTagTypes      ToolNames = "UniqueTagTypes"
	ToolNameUniqueContextTypes  ToolNames = "UniqueContextTypes"
	ToolNameUserCount           ToolNames = "UserCount"
	ToolNameGetUser             ToolNames = "GetUser"
	ToolNameDeploy              ToolNames = "Deploy"
	ToolNameRollback            ToolNames = "Rollback"
	ToolNameListDeployments     ToolNames = "ListDeployments"
	ToolNameGetDeploymentStatus ToolNames = "GetDeploymentStatus"
	ToolNameUnknown             ToolNames = "Unknown"
)

// ClassificationToolDefinitions contains all available tool definitions for classification
var ClassificationToolDefinitions = []ToolDefinition{
	NewToolDefinition(
		string(ToolNameTagCount),
		"Returns the number of tags. Use 'type' for tag categories (Course, Subject, etc.) and 'contextType' for organizational contexts (College, DoD, etc.). Use 'name' for partial name searches. Optional filters: type, contextType, name, and public status.",
	).WithParameters(NoRequiredParams, typeProperty, contextProperty, nameProperty, publicProperty, formatProperty),

	NewToolDefinition(
		string(ToolNameListTags),
		"Returns a list of tags. Use 'type' for tag categories (Course, Subject, etc.) and 'contextType' for organizational contexts (College, DoD, etc.). Use 'name' for partial name searches. Optional filters: type, contextType, name, and public status. Choose format based on user intent: 'list' for human reading, 'json' for data/API use, 'csv' for spreadsheets, 'table' for markdown.",
	).WithParameters(NoRequiredParams, typeProperty, contextProperty, nameProperty, publicProperty, formatProperty, limitProperty),

	NewToolDefinition(
		string(ToolNameListRootTags),
		"Returns a list of root tags (tags with no parent). Optional filters: name and public status. Choose format based on user intent: 'list' for human reading, 'json' for data/API use, 'csv' for spreadsheets, 'table' for markdown.",
	).WithParameters(NoRequiredParams, nameProperty, publicProperty, formatProperty, limitProperty),

	NewToolDefinition(
		string(ToolNameGetTag),
		"Returns detailed information about a specific tag by its ID. When user refers to 'tag number X', first use ListTags to get a list of tags with their IDs, then use the actual tag ID from that list. Tag IDs are CUIDs (25-character alphanumeric strings like 'cmav63fwp03ef1jmtqkh9wnvv'). The tagId parameter must be the actual CUID, not a number. Choose format based on user intent: 'list' for human reading, 'json' for data/API use, 'csv' for spreadsheets, 'table' for markdown.",
	).WithParameters([]string{"tagId"}, NewProperty("tagId", "string", "The actual CUID of the tag to retrieve (25-character alphanumeric string like 'cmav63fwp03ef1jmtqkh9wnvv'). Use ListTags first to get the CUID if user refers to a tag by number."), formatProperty),

	NewToolDefinition(
		string(ToolNameUniqueTagTypes),
		"Returns a list of all unique tag types available in the system. Choose format based on user intent: 'list' for human reading, 'json' for data/API use, 'csv' for spreadsheets, 'table' for markdown.",
	).WithParameters(NoRequiredParams, formatProperty),

	NewToolDefinition(
		string(ToolNameUniqueContextTypes),
		"Returns a list of all unique context types (organizational contexts) available in the system. Choose format based on user intent: 'list' for human reading, 'json' for data/API use, 'csv' for spreadsheets, 'table' for markdown.",
	).WithParameters(NoRequiredParams, formatProperty),

	NewToolDefinition(
		string(ToolNameUserCount),
		"Returns the number of users. The system uses intelligent date parsing to handle relative time expressions. For 'this month', 'last year', '3 months ago', etc., extract the appropriate time parameters. The system will automatically correct outdated cached dates. Examples: 'this month' should use current month and year, 'last year' should use previous year, 'last week' should use days=7. Use time-based filters: 'days' for recent users, 'months' for quarterly/annual counts, 'month' and 'year' for specific time periods, or 'since'/'until' for custom date ranges.",
	).WithParameters(NoRequiredParams, sinceProperty, untilProperty, daysProperty, monthsProperty, yearsProperty, monthProperty, yearProperty),

	NewToolDefinition(
		string(ToolNameGetUser),
		"Returns detailed information about a specific user by their email address. The userEmail parameter must be a valid email address. Choose format based on user intent: 'list' for human reading, 'json' for data/API use, 'csv' for spreadsheets, 'table' for markdown.",
	).WithParameters([]string{"userEmail"}, userEmailProperty, formatProperty),

	NewToolDefinition(
		string(ToolNameDeploy),
		"Deploys an application to the production environment. Use when user wants to deploy, push to production, or release an app. The appId parameter is required and should be the DigitalOcean App Platform app ID. The force parameter is optional and defaults to false - set to true to force a rebuild even if unchanged.",
	).WithParameters([]string{"appId"}, appIdProperty, forceProperty),

	NewToolDefinition(
		string(ToolNameRollback),
		"Rolls back an application to a previous deployment. Use when user wants to rollback, revert, or go back to a previous version. Both appId and deploymentId parameters are required. The deploymentId should be the ID of the deployment to roll back to.",
	).WithParameters([]string{"appId", "deploymentId"}, appIdProperty, deploymentIdProperty),

	NewToolDefinition(
		string(ToolNameListDeployments),
		"Lists all deployments for an application. Use when user wants to see deployment history, list deployments, or check deployment status. The appId parameter is required. Choose format based on user intent: 'list' for human reading, 'json' for data/API use, 'csv' for spreadsheets, 'table' for markdown.",
	).WithParameters([]string{"appId"}, appIdProperty, formatProperty),

	NewToolDefinition(
		string(ToolNameGetDeploymentStatus),
		"Gets the status of a specific deployment. Use when user wants to check deployment status, see if deployment is complete, or get deployment details. Both appId and deploymentId parameters are required. Choose format based on user intent: 'list' for human reading, 'json' for data/API use, 'csv' for spreadsheets, 'table' for markdown.",
	).WithParameters([]string{"appId", "deploymentId"}, appIdProperty, deploymentIdProperty, formatProperty),

	NewToolDefinition(
		string(ToolNameUnknown),
		"Use when the user's request doesn't match any other available operations.",
	).WithParameters(NoRequiredParams),
}

// ClassificationToolMap provides efficient access to both tools and names
var ClassificationToolMap = func() map[string]openai.Tool {
	toolMap := make(map[string]openai.Tool)
	for _, toolDef := range ClassificationToolDefinitions {
		toolMap[toolDef.Name] = toolDef.AsTool()
	}
	return toolMap
}()

// GetClassificationData returns both the tool definitions and the tools map
func GetClassificationDefinitions() []ToolDefinition {
	return ClassificationToolDefinitions
}

// GetClassificationTools returns the tools as a slice
func GetClassificationTools() []openai.Tool {
	tools := make([]openai.Tool, 0, len(ClassificationToolMap))
	for _, tool := range ClassificationToolMap {
		tools = append(tools, tool)
	}
	return tools
}
