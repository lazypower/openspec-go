# parsing

## Purpose

## Requirements

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

