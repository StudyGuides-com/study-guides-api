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