package tools

var contextProperty = NewProperty(
	"contextType",
	"string",
	"Filter by the context/organization where the tag is used (e.g. 'College', 'DoD', 'University', 'Company'). This is different from tag type - context refers to the organizational context.",
)

var typeProperty = NewProperty(
	"type",
	"string",
	"Filter by the tag's category/classification (e.g. 'Course', 'Subject', 'Topic', 'Department'). This is the tag's inherent type, not its organizational context.",
)

var nameProperty = NewProperty(
	"name",
	"string",
	"Search for tags by name using partial matching (e.g. 'math' will find 'Mathematics', 'Math 101', etc.)",
)

var publicProperty = NewProperty(
	"public",
	"string",
	"Filter by public status: 'true' for public tags, 'false' for private tags",
)

var formatProperty = NewProperty(
	"format",
	"string",
	"Output format: 'list' (default, human-readable), 'json' (machine-readable), 'csv' (spreadsheet), or 'table' (markdown table)",
)

var limitProperty = NewProperty(
	"limit",
	"integer",
	"Maximum number of tags to return (e.g. 10 for first 10 results)",
)

// Time-based properties for user counting
var sinceProperty = NewProperty(
	"since",
	"string",
	"Count users created since this date (ISO format: YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)",
)

var untilProperty = NewProperty(
	"until",
	"string",
	"Count users created until this date (ISO format: YYYY-MM-DD or YYYY-MM-DDTHH:MM:SS)",
)

var daysProperty = NewProperty(
	"days",
	"string",
	"Count users created in the last N days (e.g. '7' for last week, '30' for last month)",
)

var monthsProperty = NewProperty(
	"months",
	"string",
	"Count users created in the last N months (e.g. '3' for last quarter, '12' for last year)",
)

var yearsProperty = NewProperty(
	"years",
	"string",
	"Count users created in the last N years (e.g. '1' for last year, '5' for last 5 years)",
)

var monthProperty = NewProperty(
	"month",
	"string",
	"Count users created in a specific month (1-12, e.g. '7' for July)",
)

var yearProperty = NewProperty(
	"year",
	"string",
	"Count users created in a specific year (e.g. '2023' for year 2023)",
)

var userEmailProperty = NewProperty(
	"userEmail",
	"string",
	"The email address of the user to retrieve (e.g. 'user@example.com')",
)

// DevOps properties
var appIdProperty = NewProperty(
	"appId",
	"string",
	"The DigitalOcean App Platform app ID (e.g. 'abc123def456')",
)

var deploymentIdProperty = NewProperty(
	"deploymentId",
	"string",
	"The deployment ID to reference (e.g. 'def789ghi012')",
)

var forceProperty = NewProperty(
	"force",
	"string",
	"Force a rebuild even if unchanged: 'true' to force rebuild, 'false' for normal deployment (default: 'false')",
)
