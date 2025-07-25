# internal/errors Package

## Purpose and Responsibility
The `internal/errors` package defines application-specific error types for the Study Guides API. These custom errors provide semantic meaning and enable user-friendly error message transformation through middleware.

## Key Architectural Decisions

### Custom Error Variables
Uses Go's standard error package to define specific error instances that can be compared directly using `==` operator.

### Semantic Error Types
Errors are categorized by their meaning rather than implementation details:
- **Resource errors**: `ErrNotFound`, `ErrToolNotFound`
- **Validation errors**: `ErrSystemPromptEmpty`, `ErrUserPromptEmpty`  
- **Integration errors**: `ErrNoCompletionChoicesReturned`, `ErrFailedToCreateChatCompletionWithTools`

### Middleware Integration
These errors are specifically designed to be caught by the error middleware and transformed into user-friendly messages.

## Implementation Details

### Error Definitions

#### Resource Errors
- `ErrNotFound`: Generic resource not found error
- `ErrToolNotFound`: Specific to AI tool routing when tool cannot be found

#### Validation Errors  
- `ErrSystemPromptEmpty`: AI system prompt validation failure
- `ErrUserPromptEmpty`: AI user prompt validation failure

#### AI Integration Errors
- `ErrNoCompletionChoicesReturned`: OpenAI API returned empty choices array
- `ErrFailedToCreateChatCompletionWithTools`: OpenAI function calling request failed

## Usage Patterns

### Error Creation
```go
return "", errors.ErrToolNotFound
```

### Error Comparison
```go
if err == errors.ErrNotFound {
    // Handle not found case
}
```

### Middleware Processing
The error middleware catches these errors and maps them to user-friendly messages:
- `ErrToolNotFound` → "I couldn't understand how to help with that request..."
- `ErrNotFound` → "The requested resource was not found."

## Key File

### errors.go
Single file containing all application error definitions using Go's standard error package.

## Dependencies

### External Libraries
- Standard library `errors` package only

### Internal Dependencies
None - this is a leaf package that other packages depend on.

## Design Considerations

### Why Variable Errors vs Custom Types
Uses error variables instead of custom error types for:
- **Simplicity**: No need for complex error type hierarchies
- **Performance**: Direct pointer comparison is fast
- **Compatibility**: Works seamlessly with standard error handling

### Error Message Strategy  
- **Technical errors**: Defined here for internal use
- **User-friendly messages**: Handled by middleware transformation
- **Separation of concerns**: Errors focus on classification, middleware handles presentation

## Gotchas & Non-Obvious Behaviors

1. **Direct Comparison**: These errors should be compared using `==`, not string comparison

2. **Middleware Dependency**: These errors are only meaningful when the error middleware is active

3. **No Wrapping**: These are sentinel errors, not meant to wrap other errors

4. **Limited Scope**: Only covers common application-level errors, not all possible error conditions

5. **User Message Mapping**: The actual user-facing messages are defined in middleware, not here