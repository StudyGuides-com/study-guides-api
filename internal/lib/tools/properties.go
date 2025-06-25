package tools

var contextProperty = NewProperty(
	"contextType",
	"string",
	"Context type filter, e.g. 'College'",
)

var typeProperty = NewProperty(
	"type",
	"string",
	"Tag type filter, e.g. 'Course'",
)

var formatProperty = NewProperty(
	"format",
	"string",
	"Output format: 'list' (default, human-readable), 'json' (machine-readable), 'csv' (spreadsheet), or 'table' (markdown table)",
)