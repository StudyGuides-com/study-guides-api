# Claude's Notes on MCP Integration

## 🧠 Mental Model

The MCP (Model Context Protocol) system is a **generic CRUD framework** that automatically generates OpenAI function-calling tools from repository interfaces. Think of it as:

1. **Repository Pattern** → Defines what operations are available
2. **Schema Definition** → Describes parameters and types for AI
3. **Auto-generated Tools** → AI can call these without manual tool definitions
4. **Type Safety** → No more `map[string]string` - real typed structs

## 🏗️ Architecture Overview

```
ChatService (entry point)
    ↓
MCP Processor (AI coordination)
    ↓
Tag Repository Adapter (type conversion)
    ↓
Existing SqlTagStore (unchanged)
```

## 🔑 Key Files to Remember

### Core MCP System
- `/internal/repository/repository.go` - Generic interfaces (created to avoid import cycles)
- `/internal/mcp/server.go` - Main processor, AI integration, system prompt
- `/internal/mcp/handler.go` - Reflection-based command execution
- `/internal/mcp/tools_simple.go` - OpenAI tool generation

### Tag Integration
- `/internal/mcp/tag/types.go` - Domain types with **custom UnmarshalJSON** for enums
- `/internal/mcp/tag/adapter.go` - Bridges old SqlTagStore to new interface
- `/internal/mcp/tag/schema.go` - Defines available operations for AI

### Service Integration
- `/internal/services/chat.go` - Completely replaced to use MCP
- `/cmd/server/server.go` - Updated constructor call

## 🐛 Issues I Fixed

1. **OpenAI API Error**: Removed JSON response format when using tools (incompatible)
2. **Enum Marshaling**: AI sends "Category" as string, but we need `sharedv1.TagType` enum
   - Solution: Custom `UnmarshalJSON` method on `TagFilter`
3. **Import Cycles**: Created separate `/internal/repository` package
4. **System Prompt**: Made it explicit that AI MUST use tools, not text responses

## 🎯 Current State

### What Works
- All tag CRUD operations through natural language
- Enum filtering (type="Category", etc.)
- Boolean filters (public, hasChildren, isRoot)
- Count operations
- Complex combined filters

### What's Not Done
- User repository (started but removed to focus on tags)
- Question, Interaction, etc. repositories
- Create/Update/Delete operations (only Find/Count implemented)
- Batch operations
- Advanced search integration

## 🚀 Next Domain Template

When adding a new domain (e.g., User), follow this pattern:

1. **Create types** (`/internal/mcp/user/types.go`):
   ```go
   type User struct { /* fields */ }
   type UserFilter struct { /* filter fields */ }
   type UserUpdate struct { /* update fields */ }
   ```

2. **Create schema** (`/internal/mcp/user/schema.go`):
   ```go
   func GetResourceSchema() mcp.ResourceSchema {
       // Define operations and parameters
   }
   ```

3. **Create adapter** (`/internal/mcp/user/adapter.go`):
   ```go
   type UserRepositoryAdapter struct {
       store user.UserStore
   }
   // Implement Repository interface methods
   ```

4. **Register in ChatService**:
   ```go
   userRepo := user.NewUserRepositoryAdapter(store.UserStore())
   mcpProcessor.Register(user.ResourceName, userRepo, user.GetResourceSchema())
   ```

## ⚠️ Gotchas to Remember

1. **Enum Handling**: Always need custom JSON unmarshaling for protobuf enums
2. **Nil Pointers**: Protobuf optional fields need careful nil checking
3. **Type Conversion**: Lots of back-and-forth between domain types and protobuf
4. **AI Prompting**: Must be very explicit about tool usage in system prompt
5. **Payload Conversion**: The generic `convertPayload` uses JSON marshal/unmarshal

## 💭 Design Decisions

1. **Why separate repository package?** - To avoid import cycles between mcp and domain packages
2. **Why custom UnmarshalJSON?** - OpenAI sends strings, we need enums
3. **Why keep old stores?** - Minimal changes to existing codebase, just adapt
4. **Why reflection?** - Enables truly generic handlers without code generation

## 🔮 Future Improvements

1. **Code Generation**: Could generate adapters from protobuf definitions
2. **Caching**: MCP responses could be cached for common queries  
3. **Validation**: Add request validation before hitting database
4. **Metrics**: Track which operations are most used
5. **Testing**: Add unit tests for adapters and handlers

## 📝 Quick Test Command

To verify everything still works:
```bash
go run cmd/test-tag-operations/main.go
```

## 🎉 Summary

The MCP system successfully replaces the old string-based router with a type-safe, extensible framework. The pattern is established - adding new domains is now formulaic. The hardest parts (enum handling, type conversion, AI integration) are solved.

When you come back to this, start by running the test command above to make sure everything still works, then follow the "Next Domain Template" to add User or Question support.

Good luck exploring! The code should be much cleaner and more maintainable now. 🚀