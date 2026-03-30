## ADDED Requirements

### Requirement: Spec File Parsing
The parser SHALL extract structured data from spec markdown files using line-by-line scanning with regexp matching.

#### Scenario: Parse spec with requirements and scenarios
- **WHEN** parser reads a spec.md containing `## Purpose`, `## Requirements`, `### Requirement:` headers, and `#### Scenario:` blocks
- **THEN** it returns a Spec struct with title, overview (purpose text), and a slice of Requirements each containing text and Scenarios

#### Scenario: Parse requirement text
- **WHEN** parser encounters `### Requirement: Feature Name` followed by body text
- **THEN** it captures "Feature Name" as the requirement name and all text until the next header as the requirement body

#### Scenario: Parse scenario content
- **WHEN** parser encounters `#### Scenario: Happy path` followed by `- **WHEN**`, `- **THEN**`, `- **AND**` lines
- **THEN** it captures the scenario name and preserves the raw text of all bullet lines

#### Scenario: Whitespace tolerance
- **WHEN** headers contain leading/trailing whitespace (e.g., `###  Requirement:  Name  `)
- **THEN** the parser normalizes by trimming and collapsing whitespace for matching purposes

### Requirement: Change File Parsing
The parser SHALL extract structured data from change proposal markdown files.

#### Scenario: Parse proposal.md
- **WHEN** parser reads a proposal.md with `## Why`, `## What Changes`, and `## Impact` sections
- **THEN** it returns a Change struct with the extracted section content

#### Scenario: Extract title from heading
- **WHEN** proposal.md starts with `# Change: Brief description`
- **THEN** the parser extracts "Brief description" as the change title

### Requirement: Delta Parsing
The parser SHALL extract delta operations from change spec files.

#### Scenario: Parse ADDED requirements
- **WHEN** parser encounters `## ADDED Requirements` followed by requirement blocks
- **THEN** it returns Delta structs with operation "ADDED" and the parsed requirements

#### Scenario: Parse MODIFIED requirements
- **WHEN** parser encounters `## MODIFIED Requirements` followed by full requirement blocks
- **THEN** it returns Delta structs with operation "MODIFIED" and complete requirement content

#### Scenario: Parse REMOVED requirements
- **WHEN** parser encounters `## REMOVED Requirements` with requirement names and reasons
- **THEN** it returns Delta structs with operation "REMOVED" and the requirement names

#### Scenario: Parse RENAMED requirements
- **WHEN** parser encounters `## RENAMED Requirements` with FROM/TO pairs
- **THEN** it returns Delta structs with operation "RENAMED" and the old/new names

#### Scenario: Multiple delta sections in one file
- **WHEN** a delta file contains both `## ADDED Requirements` and `## MODIFIED Requirements`
- **THEN** the parser returns all deltas from all sections

#### Scenario: Header normalization
- **WHEN** delta headers have inconsistent whitespace (e.g., `##  ADDED   Requirements`)
- **THEN** the parser matches them correctly using trimmed, normalized comparison

### Requirement: Task Progress Parsing
The parser SHALL count task completion status from tasks.md files.

#### Scenario: Count completed tasks
- **WHEN** parser reads a tasks.md containing `- [x] Done task` and `- [ ] Pending task`
- **THEN** it returns TaskStatus with total=2, completed=1

#### Scenario: Case-insensitive checkbox
- **WHEN** tasks.md contains `- [X] Task` (uppercase X)
- **THEN** it counts as completed

#### Scenario: Ignore non-task lines
- **WHEN** tasks.md contains headers, paragraphs, or non-checkbox list items
- **THEN** they are not counted as tasks

### Requirement: Line Scanner Architecture
The parser SHALL use a single-pass line scanner pattern for all markdown parsing.

#### Scenario: Single pass parsing
- **WHEN** parsing any markdown file
- **THEN** the parser reads the file line-by-line in a single pass, using a state machine to track current section, requirement, and scenario context

#### Scenario: Header detection
- **WHEN** scanner encounters a line matching `^(#{1,6})\s+(.+)$`
- **THEN** it determines the header level and text, updating the parser state accordingly

#### Scenario: Content accumulation
- **WHEN** scanner encounters non-header lines between two headers
- **THEN** it accumulates them as content belonging to the current context (section body, requirement text, or scenario steps)
