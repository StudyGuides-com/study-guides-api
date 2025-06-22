package errors

import "errors"

var ErrNotFound = errors.New("resource not found")
var ErrToolNotFound = errors.New("tool not found")
var ErrSystemPromptEmpty = errors.New("system prompt cannot be empty")
var ErrUserPromptEmpty = errors.New("user prompt cannot be empty")
var ErrNoCompletionChoicesReturned = errors.New("no completion choices returned")
var ErrFailedToCreateChatCompletionWithTools = errors.New("failed to create chat completion with tools")
