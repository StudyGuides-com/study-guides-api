# internal/utils Package

## Purpose and Responsibility
The `internal/utils` package provides common utility functions used throughout the application. It includes unique ID generation and type mapping utilities for protobuf enums.

## Key Architectural Decisions

### Utility Function Approach
Provides standalone utility functions rather than methods on types, making them easily accessible across the codebase.

### CUID for Unique IDs
Uses CUID (Collision-resistant Unique Identifiers) instead of UUIDs for better properties:
- **URL-safe**: No special characters requiring encoding
- **Shorter**: More compact than UUIDs
- **Collision-resistant**: Extremely low probability of collisions
- **Sortable**: Timestamp-based prefix enables chronological sorting

### Bidirectional Type Mapping
Provides mapping utilities between related protobuf enum types (ParserType ↔ ContextType) with existence checking.

## Implementation Details

### Unique ID Generation
```go
func GetCUID() string {
    return cuid.New()
}
```
- Generates 25-character alphanumeric identifiers
- Used for all entity IDs in the database
- Thread-safe and suitable for concurrent use

### Parser-Context Type Mapping
Maps different educational content types to their organizational contexts:

#### Mapping Table
- `PARSER_TYPE_COLLEGES` → `ContextType_Colleges`
- `PARSER_TYPE_CERTIFICATIONS` → `ContextType_Certifications`
- `PARSER_TYPE_ENTRANCE_EXAMS` → `ContextType_EntranceExams`
- `PARSER_TYPE_AP_EXAMS` → `ContextType_APExams`
- `PARSER_TYPE_DOD` → `ContextType_DoD`

#### Bidirectional Conversion
```go
// Parser → Context
func GetContextTypeForParser(parserType sharedpb.ParserType) (sharedpb.ContextType, bool)

// Context → Parser  
func GetParserTypeForContext(contextType sharedpb.ContextType) (sharedpb.ParserType, bool)
```

Both functions return the mapped type and a boolean indicating whether the mapping exists.

## Usage Patterns

### ID Generation
```go
tagID := utils.GetCUID()
userID := utils.GetCUID()
```

### Type Mapping with Validation
```go
contextType, exists := utils.GetContextTypeForParser(parserType)
if !exists {
    // Handle unknown parser type
}
```

### Reverse Mapping
```go
parserType, exists := utils.GetParserTypeForContext(contextType)
if !exists {
    // Handle unknown context type
}
```

## Key File

### utils.go
Single file containing:
- CUID generation function
- Parser-Context mapping table
- Bidirectional conversion functions

## Dependencies

### External Libraries
- `github.com/lucsky/cuid` - CUID generation library

### Internal Dependencies
- `api/v1/shared` - Protobuf enum types (ParserType, ContextType)

## Design Considerations

### Why CUID over UUID
CUIDs provide several advantages for this application:
- **Database-friendly**: Better for primary keys and indexing
- **URL-safe**: Can be used directly in URLs without encoding
- **Human-readable**: Easier to debug and work with
- **Collision-resistant**: Suitable for distributed systems

### Mapping Strategy
The parser-context mapping reflects the domain model:
- **ParserType**: Represents how content is processed/imported
- **ContextType**: Represents the organizational context where content is used
- **One-to-One Mapping**: Each parser type corresponds to exactly one context type

### Error Handling
Functions return boolean flags rather than errors for missing mappings, allowing callers to decide how to handle unknown types.

## Gotchas & Non-Obvious Behaviors

1. **CUID Format**: CUIDs are 25-character strings, not the 36-character format of UUIDs

2. **Case Sensitivity**: Protobuf enum comparisons are case-sensitive

3. **Missing Mappings**: Not all possible enum values have mappings - always check the boolean return value

4. **Reverse Lookup Performance**: Context→Parser conversion uses linear search through the map

5. **Concurrent Safety**: CUID generation is thread-safe, but the mapping table is read-only after initialization

6. **Default Values**: Functions return `UNSPECIFIED` enum values when mappings don't exist

7. **Import Dependencies**: This package depends on protobuf types, creating coupling to the API layer

## Future Considerations
- Consider adding more robust error handling for type conversions
- Could benefit from caching for reverse lookups if performance becomes an issue
- May need additional utility functions as the application grows