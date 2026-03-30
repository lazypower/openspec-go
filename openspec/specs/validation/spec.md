# validation

## Purpose

## Requirements

### Requirement: Validation Constants
The validator SHALL use configurable thresholds for validation checks.

#### Scenario: Threshold values
- **WHEN** validation checks are performed
- **THEN** the following thresholds apply:
  - MIN_WHY_SECTION_LENGTH: 50
  - MAX_WHY_SECTION_LENGTH: 1000
  - MIN_PURPOSE_LENGTH: 50
  - MAX_REQUIREMENT_TEXT_LENGTH: 500
  - MAX_DELTAS_PER_CHANGE: 10 (warning only)

