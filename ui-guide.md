# Block UI Components Guide for Study Guides API

## Overview

This guide provides comprehensive instructions for building Block UI components to display various JSON response types from the Study Guides API. All API responses follow a universal wrapper format with consistent structure and content types.

## Universal Response Structure

All API responses follow this structure:

```json
{
  "type": "count|list|single|csv|table",
  "data": "varies by type",
  "message": "Human-readable description",
  "content_type": "application/json|text/plain|text/csv",
  "filters": { "optional": "filter parameters" },
  "pagination": { "optional": "pagination info" },
  "metadata": { "optional": "additional data" }
}
```

## Response Types and Components

### 1. Count Response (`type: "count"`)

**Purpose**: Display numerical counts with optional filters

**JSON Structure**:

```json
{
  "type": "count",
  "data": 42,
  "message": "Found 42 tags total.",
  "content_type": "application/json",
  "filters": {
    "type": "Course",
    "contextType": "University"
  }
}
```

**Block UI Component Requirements**:

- **Primary Display**: Large, prominent number showing the count
- **Message**: Descriptive text explaining what was counted
- **Filters Badge**: If filters present, show as small badges/tags
- **Visual Style**: Use cards or stat blocks with clear hierarchy
- **Color Coding**: Use different colors for different count types (tags, users, etc.)

**Example Layout**:

```
┌─────────────────────────────┐
│          42                 │  ← Large number
│     Tags Found              │  ← Message
│  [Course] [University]      │  ← Filter badges
└─────────────────────────────┘
```

### 2. List Response (`type: "list"`)

**Purpose**: Display lists of items in various formats

#### 2.1 Simple List (`content_type: "text/plain"`)

**JSON Structure**:

```json
{
  "type": "list",
  "data": "Category\nSubCategory\nUniversity\nRegion",
  "message": "Found 4 unique tag types",
  "content_type": "text/plain"
}
```

**Block UI Component Requirements**:

- **List Display**: Simple bulleted or numbered list
- **Item Styling**: Clean, readable text with proper spacing
- **Message Header**: Show the descriptive message above the list
- **Count Badge**: Display the count as a small badge

**Example Layout**:

```
┌─────────────────────────────┐
│ Found 4 unique tag types    │  ← Message
│                             │
│ • Category                  │  ← List items
│ • SubCategory               │
│ • University                │
│ • Region                    │
└─────────────────────────────┘
```

#### 2.2 JSON List (`content_type: "application/json"`)

**JSON Structure**:

```json
{
  "type": "list",
  "data": [
    { "id": "123", "name": "Math 101", "type": "Course" },
    { "id": "124", "name": "Physics 201", "type": "Course" }
  ],
  "message": "Found 2 courses",
  "content_type": "application/json"
}
```

**Block UI Component Requirements**:

- **Table/Grid Layout**: Display as structured table or card grid
- **Column Headers**: Show field names as headers
- **Data Formatting**: Format each field appropriately
- **Interactive Elements**: Make items clickable if they have IDs
- **Responsive Design**: Adapt to different screen sizes

**Example Layout**:

```
┌─────────────────────────────┐
│ Found 2 courses             │  ← Message
│                             │
│ ┌─────┬─────────┬─────────┐ │
│ │ ID  │ Name    │ Type    │ │  ← Headers
│ ├─────┼─────────┼─────────┤ │
│ │ 123 │ Math 101│ Course  │ │  ← Data rows
│ │ 124 │Physics  │ Course  │ │
│ └─────┴─────────┴─────────┘ │
└─────────────────────────────┘
```

### 3. CSV Response (`content_type: "text/csv"`)

**Purpose**: Display structured data in CSV format

**JSON Structure**:

```json
{
  "type": "list",
  "data": "tag_type\nCategory\nSubCategory\nUniversity\nRegion",
  "message": "Found 4 unique tag types",
  "content_type": "text/csv"
}
```

**Block UI Component Requirements**:

- **Table Display**: Convert CSV to formatted table
- **Header Row**: Style the first row as headers
- **Data Rows**: Display subsequent rows as data
- **Download Option**: Provide CSV download button
- **Copy Functionality**: Allow copying to clipboard

**Example Layout**:

```
┌─────────────────────────────┐
│ Found 4 unique tag types    │  ← Message
│ [Download CSV] [Copy]       │  ← Action buttons
│                             │
│ ┌─────────────┐             │
│ │ Tag Type    │             │  ← Header
│ ├─────────────┤             │
│ │ Category    │             │  ← Data rows
│ │ SubCategory │             │
│ │ University  │             │
│ │ Region      │             │
│ └─────────────┘             │
└─────────────────────────────┘
```

### 4. Table Response (`content_type: "text/plain"` with markdown table)

**Purpose**: Display pre-formatted markdown tables

**JSON Structure**:

```json
{
  "type": "list",
  "data": "| Tag Type | Count |\n|----------|-------|\n| Course   | 15    |\n| User     | 8     |",
  "message": "Tag type breakdown",
  "content_type": "text/plain"
}
```

**Block UI Component Requirements**:

- **Markdown Parsing**: Parse and render markdown table syntax
- **Table Styling**: Apply consistent table styling
- **Responsive Tables**: Handle overflow on small screens
- **Sorting**: Optional column sorting functionality

### 5. Single Item Response (`type: "single"`)

**Purpose**: Display detailed information about a single item

**JSON Structure**:

```json
{
  "type": "single",
  "data": {
    "id": "123",
    "name": "Math 101",
    "description": "Introduction to Mathematics",
    "type": "Course",
    "created_at": "2024-01-15T10:30:00Z"
  },
  "message": "Found user 'john@example.com'",
  "content_type": "application/json"
}
```

**Block UI Component Requirements**:

- **Detail Card**: Display as a detailed information card
- **Field Labels**: Show field names as labels
- **Data Formatting**: Format dates, IDs, and other fields appropriately
- **Action Buttons**: Include edit, delete, or navigation actions
- **Responsive Layout**: Adapt to different screen sizes

**Example Layout**:

```
┌─────────────────────────────┐
│ Found user 'john@example.com'│  ← Message
│                             │
│ ┌─────────────────────────┐ │
│ │ ID: 123                 │ │  ← Field labels
│ │ Name: Math 101          │ │
│ │ Description: Intro...   │ │
│ │ Type: Course            │ │
│ │ Created: Jan 15, 2024   │ │
│ └─────────────────────────┘ │
│                             │
│ [Edit] [Delete] [View]      │  ← Action buttons
└─────────────────────────────┘
```

## Component Design Principles

### 1. Consistency

- Use consistent spacing, typography, and color schemes
- Maintain visual hierarchy across all component types
- Follow established design patterns for similar data types

### 2. Accessibility

- Ensure proper contrast ratios
- Include ARIA labels for screen readers
- Provide keyboard navigation support
- Use semantic HTML elements

### 3. Responsiveness

- Design for mobile-first approach
- Ensure components work on all screen sizes
- Handle data overflow gracefully
- Provide appropriate touch targets

### 4. Performance

- Implement virtual scrolling for large lists
- Use lazy loading for images or heavy content
- Optimize re-renders with proper state management
- Cache frequently accessed data

## Color Scheme Guidelines

### Primary Colors

- **Blue**: Primary actions, links, and interactive elements
- **Green**: Success states and positive metrics
- **Red**: Error states and destructive actions
- **Orange**: Warning states and attention-grabbing elements
- **Gray**: Neutral text, borders, and secondary information

### Semantic Colors

- **Count Components**: Use blue for general counts, green for positive metrics
- **List Components**: Use neutral colors with blue accents for interactive elements
- **Error States**: Use red for errors, orange for warnings
- **Success States**: Use green for successful operations

## Typography Guidelines

### Font Hierarchy

- **H1**: Component titles and major headings (24px, bold)
- **H2**: Section headers (20px, semi-bold)
- **H3**: Subsection headers (18px, medium)
- **Body**: Regular text content (16px, regular)
- **Small**: Secondary information and captions (14px, regular)
- **Micro**: Metadata and timestamps (12px, regular)

### Font Weights

- **Bold (700)**: Primary headings and important information
- **Semi-bold (600)**: Secondary headings
- **Medium (500)**: Emphasis and interactive elements
- **Regular (400)**: Body text and general content
- **Light (300)**: Subtle text and decorative elements

## Interactive Elements

### 1. Buttons

- **Primary**: Solid background, high contrast
- **Secondary**: Outlined or ghost style
- **Tertiary**: Text-only buttons
- **Icon**: Buttons with icons only

### 2. Links

- Use consistent link styling
- Include hover and focus states
- Distinguish between internal and external links

### 3. Form Elements

- Clear labels and placeholders
- Proper validation states
- Accessible error messages
- Consistent input styling

## Error Handling

### 1. Empty States

- Provide helpful empty state messages
- Include suggested actions
- Use appropriate illustrations or icons

### 2. Loading States

- Show skeleton loaders for content
- Use spinners for quick operations
- Provide progress indicators for long operations

### 3. Error States

- Display clear error messages
- Provide retry options
- Include fallback content when possible

## Implementation Examples

### React Component Structure

```jsx
// Example Count Component
const CountBlock = ({ response }) => {
  const { data, message, filters } = response;

  return (
    <div className="count-block">
      <div className="count-number">{data}</div>
      <div className="count-message">{message}</div>
      {filters && <FilterBadges filters={filters} />}
    </div>
  );
};

// Example List Component
const ListBlock = ({ response }) => {
  const { data, message, content_type } = response;

  if (content_type === "text/csv") {
    return <CSVTable data={data} message={message} />;
  }

  if (content_type === "application/json") {
    return <JSONTable data={data} message={message} />;
  }

  return <SimpleList data={data} message={message} />;
};
```

### CSS Structure

```css
/* Base component styles */
.api-response-block {
  border-radius: 8px;
  padding: 16px;
  margin: 8px 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* Count component styles */
.count-block {
  text-align: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.count-number {
  font-size: 48px;
  font-weight: bold;
  margin-bottom: 8px;
}

/* List component styles */
.list-block {
  background: white;
  border: 1px solid #e1e5e9;
}

.list-item {
  padding: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.list-item:last-child {
  border-bottom: none;
}
```

## Testing Guidelines

### 1. Visual Testing

- Test on different screen sizes
- Verify color contrast ratios
- Check for visual consistency
- Test with different data sets

### 2. Functional Testing

- Test all interactive elements
- Verify data display accuracy
- Test error handling scenarios
- Check accessibility compliance

### 3. Performance Testing

- Measure component render times
- Test with large data sets
- Verify memory usage
- Check bundle size impact

## Best Practices

### 1. Code Organization

- Separate concerns (data, presentation, logic)
- Use consistent naming conventions
- Implement proper error boundaries
- Follow component composition patterns

### 2. State Management

- Use appropriate state management solutions
- Implement proper loading states
- Handle error states gracefully
- Cache data when appropriate

### 3. Documentation

- Document component props and usage
- Include usage examples
- Maintain changelog
- Provide migration guides

This guide provides a comprehensive foundation for building Block UI components that effectively display the various JSON response types from the Study Guides API. Follow these guidelines to create consistent, accessible, and performant user interfaces.
