# Prisma Schema Architecture

## Overview
The Prisma schema folder contains a comprehensive data model for a study guides platform. This folder is **READ ONLY** as these files are controlled by the study-guides-com project and serve as reference for building database code in this API service.

## Architecture Principles

### Modular Schema Design
The schema is organized into separate `.prisma` files, each focused on a specific domain:
- **Core Entities**: User, Tag, Question, Interaction
- **Authentication**: Account, Session, VerificationToken
- **Study Features**: Test, Survival, Challenge, Badge
- **Metadata**: Rating, Report, Subscription, Import
- **Supporting**: Browser, Passage, Prompt, Article

### Schema Folder Pattern
Uses Prisma's `prismaSchemaFolder` preview feature enabling:
- Better organization through file separation
- Domain-driven schema structure
- Easier maintenance and collaboration
- Clear separation of concerns

## Key Architectural Decisions

### Hierarchical Tag System
Tags form a self-referencing tree structure:
- **Parent-Child Relationships**: `parentTagId` creates hierarchical organization
- **Multiple Tag Types**: Category, Course, Topic, UserContent, etc.
- **Context Segregation**: Colleges, Certifications, EntranceExams, APExams
- **Access Control**: Public/private tags with owner-based permissions

### Dual User/Anonymous Support
Most user-centric models support both authenticated and anonymous users:
- **User ID**: For authenticated users
- **Browser ID**: For anonymous session tracking
- **Graceful Migration**: Anonymous data can be transferred to user accounts

### Content Access Control
Implements flexible permission system:
- **AccessType Enum**: Public, Private, ReadOnly, ReadWrite
- **Access Lists**: Tag/Question specific permissions
- **Invite System**: Email-based collaboration with expiration
- **Ownership Model**: Content creators maintain control

### Rich Interaction Tracking
Comprehensive user behavior analytics:
- **Interaction Types**: Answer patterns, reveals, difficulty ratings
- **Study Methods**: StudyGuide, Flashcards, MultipleChoice, etc.
- **Strength Scoring**: Weighted interaction values
- **Progress Tracking**: Topic-level completion status

### Gamification Architecture
Multi-layered engagement system:
- **Experience Points**: Activity-based rewards
- **Challenges**: Action-target completion tracking
- **Badges**: Achievement recognition system
- **Progress Metrics**: Completion tracking across study methods

## Core Entity Relationships

### User Entity (`user.prisma`)
Central hub connecting all user-related data:
- **Authentication**: Accounts, Sessions via NextAuth.js pattern
- **Content Creation**: Tags, Questions with ownership
- **Learning Progress**: Interactions, Topic Progress, Test/Survival sessions
- **Social Features**: Favorites, Recents, Ratings, Reports
- **Gamification**: Badges, Challenges, Experience Points
- **Monetization**: Subscriptions with Stripe integration

### Tag Entity (`tag.prisma`)
Hierarchical content organization:
- **Tree Structure**: Self-referencing parent-child relationships
- **Rich Metadata**: Description, context, content ratings
- **Access Control**: Public/private with granular permissions
- **Usage Tracking**: hasChildren, hasQuestions optimization flags
- **Search Optimization**: GIN indexes on metaTags arrays

### Question Entity (`question.prisma`)
Core learning content with rich features:
- **Content**: Question/answer text with multimedia support
- **Difficulty Tracking**: Automatic ratio calculation from interactions
- **Tag Associations**: Many-to-many through QuestionTag junction
- **Access Control**: Similar to tags with ownership model
- **Analytics**: Correct/incorrect counts for adaptive learning

### Interaction Tracking (`interaction.prisma`)
Detailed user behavior capture:
- **Activity Types**: Correct/incorrect answers, reveals, difficulty ratings
- **Study Context**: Which study method was used
- **Scoring**: Strength-based weighting for spaced repetition
- **Temporal Data**: When interactions occurred for analytics

## Data Modeling Patterns

### Composite Primary Keys
Junction tables use composite keys:
```prisma
@@id([questionId, tagId])  // QuestionTag
@@id([userId, tagId])      // UserTagRating
```

### Timestamp Patterns
Consistent temporal tracking:
- **createdAt**: Record creation (always present)
- **updatedAt**: Modification tracking (@updatedAt)
- **occurredAt**: Event-specific timestamps

### Index Strategy
Performance-optimized indexing:
- **Single Column**: Primary lookup fields
- **Composite**: Multi-field queries (userId + tagId)
- **GIN Indexes**: Array field searches (metaTags, contentDescriptors)
- **Conditional**: Type-specific optimizations

### JSON Metadata Fields
Flexible schema extension:
- **metadata**: Domain-specific additional data
- **progressDetails**: Challenge progress tracking
- **contentDescriptors**: Rich content classification

### Cascade Deletion Patterns
Data integrity through relationships:
- **User Deletion**: Cascades to owned content and interactions
- **Content Deletion**: Preserves analytics but removes access
- **Batch Operations**: Import tracking with cascading cleanup

## Enum Definitions

### Study System Enums
- **StudyMethod**: StudyGuide, Flashcards, MultipleChoice, Survival, Test
- **InteractionType**: Answer patterns, reveals, difficulty feedback
- **TagType**: 27 different content organization types
- **ContextType**: Platform-wide content categorization

### Access Control Enums
- **AccessType**: Public, Private, ReadOnly, ReadWrite
- **TagInviteStatus**: Pending, Accepted, Rejected
- **ContentRatingType**: ESRB-style content ratings

### Session Management Enums
- **TestSessionType**: Quiz vs Exam modes
- **AnswerStatus**: Correct, Incorrect, Unanswered tracking
- **SubscriptionType/Status**: Monetization states

## Performance Considerations

### Strategic Indexing
- **User Lookups**: Optimized for dashboard queries
- **Content Discovery**: Tag/question browsing patterns
- **Analytics**: Time-series interaction data
- **Search**: GIN indexes for array field searches

### Denormalization Patterns
- **hasChildren/hasQuestions**: Avoid expensive COUNT queries
- **correctCount/incorrectCount**: Pre-calculated statistics
- **accessCount**: Usage metrics without joins

### Batch Operations
- **ImportBatch**: Bulk content operations with tracking
- **Cascade Controls**: Selective data cleanup strategies

## Key Files Structure

### Domain Entities
- **main.prisma**: Database config, global enums, helper models
- **user.prisma**: User accounts, progress, preferences
- **tag.prisma**: Content organization hierarchy
- **question.prisma**: Learning content with associations
- **interaction.prisma**: User behavior tracking

### Feature Modules
- **test.prisma**: Formal assessment sessions
- **survival.prisma**: Timed challenge mode
- **challenge.prisma**: Gamification objectives
- **badge.prisma**: Achievement system
- **experiencePoint.prisma**: Reward tracking

### Infrastructure
- **account.prisma**: OAuth provider integration
- **session.prisma**: Authentication state
- **subscription.prisma**: Monetization tracking
- **import.prisma**: Content migration system

## Gotchas & Non-Obvious Behaviors

1. **Schema Folder Feature**: Requires `prismaSchemaFolder` preview feature - not yet stable

2. **Dual Identity Support**: Most models support both userId AND browserId for anonymous user tracking

3. **Cascade Complexity**: User deletion cascades extensively - review before implementing

4. **Access Control Inheritance**: Tags inherit permissions but questions have independent access control

5. **Import Batch Tracking**: All imported content links to batch for rollback capabilities

6. **Content Rating System**: Implements ESRB-style ratings with content descriptors array

7. **Challenge Snapshots**: UserChallenge stores snapshots of challenge properties to prevent data inconsistency

8. **Expiring Invites**: TagInvite has database-generated expiration (24 hours default)

9. **Progress Uniqueness**: Complex unique constraints prevent duplicate progress entries

10. **GIN Index Requirements**: Array field searches require PostgreSQL GIN indexes

## Database Configuration

### Provider Setup
- **Database**: PostgreSQL (required for advanced features)
- **Shadow Database**: Required for migrations
- **Client Generation**: JavaScript/TypeScript target
- **Preview Features**: Schema folder organization

### Required Extensions
- **Arrays**: Native PostgreSQL array support
- **JSON**: Rich metadata storage
- **GIN Indexes**: Full-text and array searching
- **Generated Columns**: For computed expiration dates

## Integration Notes

This schema is designed for a full-stack learning platform with:
- **NextAuth.js**: Authentication pattern compliance
- **Stripe**: Payment processing integration
- **Algolia**: Search indexing support
- **Anonymous Users**: Pre-authentication feature usage
- **Data Migration**: Import/export capabilities

The schema serves as the single source of truth for the database structure, with this API service consuming but not modifying these definitions.