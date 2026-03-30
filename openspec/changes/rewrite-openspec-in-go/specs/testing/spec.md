## ADDED Requirements

### Requirement: Test-to-Implementation Cross-Reference Matrix
The project SHALL maintain a cross-reference matrix mapping every test function to the requirement it verifies and the implementation file it exercises.

#### Scenario: Matrix format
- **WHEN** the test matrix is consulted
- **THEN** it provides a table with columns: Test Function, File, Requirement (spec:requirement name), Implementation File

#### Scenario: Coverage completeness
- **WHEN** a new requirement is added to any spec
- **THEN** the matrix MUST be updated to include corresponding test functions before the implementation is considered complete

### Requirement: Parser Unit Tests
The parser package SHALL have unit tests covering all parsing paths.

#### Scenario: Spec parsing tests
- **WHEN** `parser_test.go` runs
- **THEN** it verifies:
  - Extracting title from `# Spec Title` header → `TestParseSpec_Title` → validates Spec File Parsing
  - Extracting purpose from `## Purpose` section → `TestParseSpec_Purpose` → validates Spec File Parsing
  - Extracting requirements from `### Requirement:` blocks → `TestParseSpec_Requirements` → validates Spec File Parsing
  - Extracting scenarios from `#### Scenario:` blocks → `TestParseSpec_Scenarios` → validates Spec File Parsing
  - Handling whitespace in headers → `TestParseSpec_WhitespaceNormalization` → validates Spec File Parsing (whitespace tolerance)
  - Empty/missing sections → `TestParseSpec_EmptySections` → validates Spec File Parsing

#### Scenario: Change parsing tests
- **WHEN** `parser_test.go` runs
- **THEN** it verifies:
  - Extracting change title → `TestParseChange_Title` → validates Change File Parsing
  - Extracting Why/What/Impact sections → `TestParseChange_Sections` → validates Change File Parsing
  - Missing required sections → `TestParseChange_MissingSections` → validates Change File Parsing

#### Scenario: Delta parsing tests
- **WHEN** `parser_test.go` runs
- **THEN** it verifies:
  - ADDED section extraction → `TestParseDelta_Added` → validates Delta Parsing (ADDED)
  - MODIFIED section extraction → `TestParseDelta_Modified` → validates Delta Parsing (MODIFIED)
  - REMOVED section extraction → `TestParseDelta_Removed` → validates Delta Parsing (REMOVED)
  - RENAMED section extraction → `TestParseDelta_Renamed` → validates Delta Parsing (RENAMED)
  - Multiple sections in one file → `TestParseDelta_MultipleSections` → validates Delta Parsing (multiple sections)
  - Header normalization → `TestParseDelta_HeaderNormalization` → validates Delta Parsing (header normalization)

#### Scenario: Task progress tests
- **WHEN** `parser_test.go` runs
- **THEN** it verifies:
  - Counting completed/total tasks → `TestParseTaskProgress_Counts` → validates Task Progress Parsing
  - Case-insensitive checkbox → `TestParseTaskProgress_CaseInsensitive` → validates Task Progress Parsing (case-insensitive)
  - Non-task lines ignored → `TestParseTaskProgress_IgnoresNonTasks` → validates Task Progress Parsing (ignore non-task)

### Requirement: Validator Unit Tests
The validator package SHALL have unit tests covering all validation rules.

#### Scenario: Change validation tests
- **WHEN** `validator_test.go` runs
- **THEN** it verifies:
  - No deltas → ERROR → `TestValidateChange_NoDeltasError` → validates Change Validation Rules (must have delta)
  - Missing scenarios → ERROR → `TestValidateChange_NoScenariosError` → validates Change Validation Rules (scenarios required)
  - Missing SHALL/MUST → ERROR → `TestValidateChange_NoShallMustError` → validates Change Validation Rules (SHALL/MUST)
  - Missing Why section → ERROR → `TestValidateChange_MissingWhyError` → validates Change Validation Rules (required sections)
  - Why too short → WARNING → `TestValidateChange_WhyTooShort` → validates Change Validation Rules (why length)
  - Duplicate requirements → ERROR → `TestValidateChange_DuplicateRequirements` → validates Change Validation Rules (no duplicates)
  - Valid change passes → `TestValidateChange_ValidPasses` → validates Change Validation Rules

#### Scenario: Spec validation tests
- **WHEN** `validator_test.go` runs
- **THEN** it verifies:
  - Missing Purpose → ERROR → `TestValidateSpec_MissingPurposeError` → validates Spec Validation Rules (purpose/requirements)
  - Missing Requirements section → ERROR → `TestValidateSpec_MissingRequirementsError` → validates Spec Validation Rules
  - Purpose too short → ERROR → `TestValidateSpec_PurposeTooShortError` → validates Spec Validation Rules (purpose length)
  - Requirement without scenario → ERROR → `TestValidateSpec_RequirementNoScenario` → validates Spec Validation Rules (scenarios)
  - Valid spec passes → `TestValidateSpec_ValidPasses` → validates Spec Validation Rules

#### Scenario: Strict mode tests
- **WHEN** `validator_test.go` runs with strict fixtures
- **THEN** it verifies:
  - Cross-spec conflicts detected → `TestValidateStrict_CrossSpecConflict` → validates Strict Validation Mode (conflicts)
  - Malformed RENAMED pairs → `TestValidateStrict_MalformedRenamed` → validates Strict Validation Mode (RENAMED)
  - MODIFIED completeness warning → `TestValidateStrict_ModifiedCompleteness` → validates Strict Validation Mode (completeness)

### Requirement: Archive Unit Tests
The archive package SHALL have unit tests covering merge operations and workflow.

#### Scenario: Merge operation tests
- **WHEN** `archive_test.go` runs
- **THEN** it verifies:
  - ADDED appends to spec → `TestMerge_AddedAppendsRequirement` → validates Delta Merge Engine (ADDED)
  - MODIFIED replaces requirement → `TestMerge_ModifiedReplacesRequirement` → validates Delta Merge Engine (MODIFIED)
  - REMOVED deletes requirement → `TestMerge_RemovedDeletesRequirement` → validates Delta Merge Engine (REMOVED)
  - RENAMED updates header → `TestMerge_RenamedUpdatesHeader` → validates Delta Merge Engine (RENAMED)
  - Missing requirement for MODIFIED → ERROR → `TestMerge_ModifiedNotFoundError` → validates Delta Merge Engine (not found)
  - Missing requirement for REMOVED → ERROR → `TestMerge_RemovedNotFoundError` → validates Delta Merge Engine (not found)
  - New spec creation from ADDED → `TestMerge_NewSpecCreation` → validates Delta Merge Engine (new spec)

#### Scenario: Archive workflow tests
- **WHEN** `archive_test.go` runs
- **THEN** it verifies:
  - Full archive with merge → `TestArchive_FullWorkflow` → validates Archive Workflow (pre-archive validation, directory move)
  - Skip-specs mode → `TestArchive_SkipSpecs` → validates Archive Workflow (skip spec updates)
  - Validation failure aborts → `TestArchive_ValidationAborts` → validates Archive Workflow (pre-archive validation)
  - Post-merge validation failure rolls back → `TestArchive_PostMergeRollback` → validates Archive Workflow (post-merge validation)

#### Scenario: Atomic write tests
- **WHEN** `archive_test.go` runs
- **THEN** it verifies:
  - Temp file write-then-rename → `TestAtomicWrite_Success` → validates Atomic File Operations (write-then-move)
  - Multi-spec all-or-nothing → `TestAtomicWrite_MultiSpecRollback` → validates Atomic File Operations (multi-spec)

### Requirement: Editor Integration Tests
The editor package SHALL have unit tests covering configurator behavior.

#### Scenario: Claude Code tests
- **WHEN** `editor_test.go` runs
- **THEN** it verifies:
  - Slash commands created → `TestClaudeCode_Configure` → validates Claude Code Integration (slash commands)
  - CLAUDE.md injection → `TestClaudeCode_ClaudeMdInjection` → validates Claude Code Integration (CLAUDE.md injection)
  - Update refreshes files → `TestClaudeCode_Update` → validates Claude Code Integration (update)
  - Marker block replaced correctly → `TestClaudeCode_MarkerBlockReplace` → validates Claude Code Integration (CLAUDE.md injection)

#### Scenario: OpenCode tests
- **WHEN** `editor_test.go` runs
- **THEN** it verifies:
  - Configuration created → `TestOpenCode_Configure` → validates OpenCode Integration
  - Update refreshes → `TestOpenCode_Update` → validates OpenCode Integration (update)

#### Scenario: Codex tests
- **WHEN** `editor_test.go` runs
- **THEN** it verifies:
  - Prompts created → `TestCodex_Configure` → validates Codex Integration
  - Update refreshes → `TestCodex_Update` → validates Codex Integration (update)

#### Scenario: Goose tests
- **WHEN** `editor_test.go` runs
- **THEN** it verifies:
  - Configuration created → `TestGoose_Configure` → validates Goose Integration
  - Update refreshes → `TestGoose_Update` → validates Goose Integration (update)

### Requirement: Integration Tests
The project SHALL include integration tests that exercise commands end-to-end against a temporary filesystem.

#### Scenario: Init integration test (non-interactive)
- **WHEN** `integration_test.go` runs TestInit_NonInteractive
- **THEN** it executes `openspec init --tools all` in a temp directory and verifies the complete directory structure and file contents
- **AND** validates Init Command (non-interactive scenarios)

#### Scenario: Init integration test (no tools)
- **WHEN** `integration_test.go` runs TestInit_NoTools
- **THEN** it executes `openspec init --tools none` and verifies directory structure without editor configs

#### Scenario: Full lifecycle integration test
- **WHEN** `integration_test.go` runs TestFullLifecycle
- **THEN** it executes init → create change files → validate → archive in sequence
- **AND** verifies specs are correctly updated after archive
- **AND** verifies the change directory is moved to archive
- **AND** validates Archive Workflow, Validate Command, and Delta Merge Engine end-to-end

#### Scenario: View integration test
- **WHEN** `integration_test.go` runs TestView
- **THEN** it creates a project with specs and changes, runs `openspec view`, and verifies the output contains expected metrics and formatting
- **AND** validates Terminal Dashboard Rendering

### Requirement: Testdata Fixtures
The project SHALL maintain testdata directories with representative fixtures for all test scenarios.

#### Scenario: Fixture organization
- **WHEN** tests reference fixture files
- **THEN** fixtures are organized under `testdata/` as:
  - `testdata/valid-spec/` — well-formed spec files
  - `testdata/valid-change/` — well-formed change proposals with deltas
  - `testdata/invalid-spec/` — spec files with known validation errors
  - `testdata/invalid-change/` — change proposals with known validation errors
  - `testdata/merge/` — before/after pairs for merge testing

#### Scenario: Fixture documentation
- **WHEN** a fixture file is created
- **THEN** it includes a comment or companion `.expected` file documenting what the fixture tests

### Requirement: Test Cross-Reference Matrix
The following matrix SHALL be maintained as the authoritative mapping between tests, requirements, and implementations.

#### Scenario: Matrix contents

| Test Function | Test File | Requirement | Implementation File |
|---|---|---|---|
| TestParseSpec_Title | parser/parser_test.go | parsing:Spec File Parsing | parser/spec.go |
| TestParseSpec_Purpose | parser/parser_test.go | parsing:Spec File Parsing | parser/spec.go |
| TestParseSpec_Requirements | parser/parser_test.go | parsing:Spec File Parsing | parser/spec.go |
| TestParseSpec_Scenarios | parser/parser_test.go | parsing:Spec File Parsing | parser/spec.go |
| TestParseSpec_WhitespaceNormalization | parser/parser_test.go | parsing:Spec File Parsing | parser/spec.go |
| TestParseSpec_EmptySections | parser/parser_test.go | parsing:Spec File Parsing | parser/spec.go |
| TestParseChange_Title | parser/parser_test.go | parsing:Change File Parsing | parser/change.go |
| TestParseChange_Sections | parser/parser_test.go | parsing:Change File Parsing | parser/change.go |
| TestParseChange_MissingSections | parser/parser_test.go | parsing:Change File Parsing | parser/change.go |
| TestParseDelta_Added | parser/parser_test.go | parsing:Delta Parsing | parser/delta.go |
| TestParseDelta_Modified | parser/parser_test.go | parsing:Delta Parsing | parser/delta.go |
| TestParseDelta_Removed | parser/parser_test.go | parsing:Delta Parsing | parser/delta.go |
| TestParseDelta_Renamed | parser/parser_test.go | parsing:Delta Parsing | parser/delta.go |
| TestParseDelta_MultipleSections | parser/parser_test.go | parsing:Delta Parsing | parser/delta.go |
| TestParseDelta_HeaderNormalization | parser/parser_test.go | parsing:Delta Parsing | parser/delta.go |
| TestParseTaskProgress_Counts | parser/parser_test.go | parsing:Task Progress Parsing | parser/task.go |
| TestParseTaskProgress_CaseInsensitive | parser/parser_test.go | parsing:Task Progress Parsing | parser/task.go |
| TestParseTaskProgress_IgnoresNonTasks | parser/parser_test.go | parsing:Task Progress Parsing | parser/task.go |
| TestValidateChange_NoDeltasError | validator/validator_test.go | validation:Change Validation Rules | validator/change.go |
| TestValidateChange_NoScenariosError | validator/validator_test.go | validation:Change Validation Rules | validator/change.go |
| TestValidateChange_NoShallMustError | validator/validator_test.go | validation:Change Validation Rules | validator/change.go |
| TestValidateChange_MissingWhyError | validator/validator_test.go | validation:Change Validation Rules | validator/change.go |
| TestValidateChange_WhyTooShort | validator/validator_test.go | validation:Change Validation Rules | validator/change.go |
| TestValidateChange_DuplicateRequirements | validator/validator_test.go | validation:Change Validation Rules | validator/change.go |
| TestValidateChange_ValidPasses | validator/validator_test.go | validation:Change Validation Rules | validator/change.go |
| TestValidateSpec_MissingPurposeError | validator/validator_test.go | validation:Spec Validation Rules | validator/spec.go |
| TestValidateSpec_MissingRequirementsError | validator/validator_test.go | validation:Spec Validation Rules | validator/spec.go |
| TestValidateSpec_PurposeTooShortError | validator/validator_test.go | validation:Spec Validation Rules | validator/spec.go |
| TestValidateSpec_RequirementNoScenario | validator/validator_test.go | validation:Spec Validation Rules | validator/spec.go |
| TestValidateSpec_ValidPasses | validator/validator_test.go | validation:Spec Validation Rules | validator/spec.go |
| TestValidateStrict_CrossSpecConflict | validator/validator_test.go | validation:Strict Validation Mode | validator/strict.go |
| TestValidateStrict_MalformedRenamed | validator/validator_test.go | validation:Strict Validation Mode | validator/strict.go |
| TestValidateStrict_ModifiedCompleteness | validator/validator_test.go | validation:Strict Validation Mode | validator/strict.go |
| TestValidateConcurrent_BatchResults | validator/validator_test.go | validation:Concurrent Validation | validator/concurrent.go |
| TestValidateConcurrent_DeterministicOrder | validator/validator_test.go | validation:Concurrent Validation | validator/concurrent.go |
| TestValidateConcurrent_SingleItemNoop | validator/validator_test.go | validation:Concurrent Validation | validator/concurrent.go |
| TestMerge_AddedAppendsRequirement | archive/archive_test.go | archive:Delta Merge Engine | archive/archive.go |
| TestMerge_ModifiedReplacesRequirement | archive/archive_test.go | archive:Delta Merge Engine | archive/archive.go |
| TestMerge_RemovedDeletesRequirement | archive/archive_test.go | archive:Delta Merge Engine | archive/archive.go |
| TestMerge_RenamedUpdatesHeader | archive/archive_test.go | archive:Delta Merge Engine | archive/archive.go |
| TestMerge_ModifiedNotFoundError | archive/archive_test.go | archive:Delta Merge Engine | archive/archive.go |
| TestMerge_RemovedNotFoundError | archive/archive_test.go | archive:Delta Merge Engine | archive/archive.go |
| TestMerge_NewSpecCreation | archive/archive_test.go | archive:Delta Merge Engine | archive/archive.go |
| TestArchive_FullWorkflow | archive/archive_test.go | archive:Archive Workflow | archive/archive.go |
| TestArchive_SkipSpecs | archive/archive_test.go | archive:Archive Workflow | archive/archive.go |
| TestArchive_ValidationAborts | archive/archive_test.go | archive:Archive Workflow | archive/archive.go |
| TestArchive_PostMergeRollback | archive/archive_test.go | archive:Archive Workflow | archive/archive.go |
| TestAtomicWrite_Success | archive/archive_test.go | archive:Atomic File Operations | archive/archive.go |
| TestAtomicWrite_MultiSpecRollback | archive/archive_test.go | archive:Atomic File Operations | archive/archive.go |
| TestClaudeCode_Configure | editor/editor_test.go | editor-integration:Claude Code Integration | editor/claude.go |
| TestClaudeCode_ClaudeMdInjection | editor/editor_test.go | editor-integration:Claude Code Integration | editor/claude.go |
| TestClaudeCode_Update | editor/editor_test.go | editor-integration:Claude Code Integration | editor/claude.go |
| TestClaudeCode_MarkerBlockReplace | editor/editor_test.go | editor-integration:Claude Code Integration | editor/claude.go |
| TestOpenCode_Configure | editor/editor_test.go | editor-integration:OpenCode Integration | editor/opencode.go |
| TestOpenCode_Update | editor/editor_test.go | editor-integration:OpenCode Integration | editor/opencode.go |
| TestCodex_Configure | editor/editor_test.go | editor-integration:Codex Integration | editor/codex.go |
| TestCodex_Update | editor/editor_test.go | editor-integration:Codex Integration | editor/codex.go |
| TestGoose_Configure | editor/editor_test.go | editor-integration:Goose Integration | editor/goose.go |
| TestGoose_Update | editor/editor_test.go | editor-integration:Goose Integration | editor/goose.go |
| TestTaskStatus_Percent | model/model_test.go | model:TaskStatus | model/change.go |
| TestValidationReport_HasErrors | model/model_test.go | model:ValidationReport | model/validation.go |
| TestProgressBar | output/output_test.go | dashboard:Terminal Dashboard Rendering | output/color.go |
| TestFormatProgress | output/output_test.go | dashboard:Terminal Dashboard Rendering | output/progress.go |
| TestWriteJSON | output/output_test.go | validation:Validation Report Format | output/json.go |
| TestSetNoColor | output/output_test.go | dashboard:Color and Formatting | output/color.go |
| TestColorFunctions | output/output_test.go | dashboard:Color and Formatting | output/color.go |
| TestTable_Render | output/output_test.go | dashboard:Terminal Dashboard Rendering | output/table.go |
| TestTable_Empty | output/output_test.go | dashboard:Terminal Dashboard Rendering | output/table.go |
| TestGet | template/templates_test.go | editor-integration:Template Embedding | template/templates.go |
| TestGet_NotFound | template/templates_test.go | editor-integration:Template Embedding | template/templates.go |
| TestRender | template/templates_test.go | editor-integration:Template Embedding | template/templates.go |
| TestMustRender | template/templates_test.go | editor-integration:Template Embedding | template/templates.go |
| TestMustRender_Panic | template/templates_test.go | editor-integration:Template Embedding | template/templates.go |
| TestRaw | template/templates_test.go | editor-integration:Template Embedding | template/templates.go |
| TestRaw_NotFound | template/templates_test.go | editor-integration:Template Embedding | template/templates.go |
| TestAllTemplatesExist | template/templates_test.go | editor-integration:Template Embedding | template/templates.go |
| TestInit_NonInteractive | cmd/cmd_test.go | cli:Init Command | cmd/init.go |
| TestInit_NoTools | cmd/cmd_test.go | cli:Init Command | cmd/init.go |
| TestInit_AlreadyInitialized | cmd/cmd_test.go | cli:Init Command | cmd/init.go |
| TestInit_CustomPath | cmd/cmd_test.go | cli:Init Command | cmd/init.go |
| TestUpdate | cmd/cmd_test.go | cli:Update Command | cmd/update.go |
| TestList_Changes | cmd/cmd_test.go | cli:List Command | cmd/list.go |
| TestList_Specs | cmd/cmd_test.go | cli:List Command | cmd/list.go |
| TestList_Empty | cmd/cmd_test.go | cli:List Command | cmd/list.go |
| TestShow_ChangeText | cmd/cmd_test.go | cli:Show Command | cmd/show.go |
| TestShow_ChangeJSON | cmd/cmd_test.go | cli:Show Command | cmd/show.go |
| TestShow_SpecText | cmd/cmd_test.go | cli:Show Command | cmd/show.go |
| TestShow_SpecJSON | cmd/cmd_test.go | cli:Show Command | cmd/show.go |
| TestShow_NotFound | cmd/cmd_test.go | cli:Show Command | cmd/show.go |
| TestShow_DeltasOnly | cmd/cmd_test.go | cli:Show Command | cmd/show.go |
| TestShow_RequirementsOnly | cmd/cmd_test.go | cli:Show Command | cmd/show.go |
| TestShow_SingleRequirement | cmd/cmd_test.go | cli:Show Command | cmd/show.go |
| TestShow_TypeFlag | cmd/cmd_test.go | cli:Show Command | cmd/show.go |
| TestValidate_SingleChange | cmd/cmd_test.go | cli:Validate Command | cmd/validate.go |
| TestValidate_SingleSpec | cmd/cmd_test.go | cli:Validate Command | cmd/validate.go |
| TestValidate_All | cmd/cmd_test.go | cli:Validate Command | cmd/validate.go |
| TestValidate_JSON | cmd/cmd_test.go | validation:Validation Report Format | cmd/validate.go |
| TestValidate_Strict | cmd/cmd_test.go | cli:Validate Command | cmd/validate.go |
| TestValidate_Concurrency | cmd/cmd_test.go | validation:Concurrent Validation | cmd/validate.go |
| TestValidate_Changes | cmd/cmd_test.go | cli:Validate Command | cmd/validate.go |
| TestValidate_Specs | cmd/cmd_test.go | cli:Validate Command | cmd/validate.go |
| TestValidate_NotFound | cmd/cmd_test.go | cli:Validate Command | cmd/validate.go |
| TestView_Dashboard | cmd/cmd_test.go | dashboard:Terminal Dashboard Rendering | cmd/view.go |
| TestView_EmptyProject | cmd/cmd_test.go | dashboard:Terminal Dashboard Rendering | cmd/view.go |
| TestView_WithArchived | cmd/cmd_test.go | dashboard:Terminal Dashboard Rendering | cmd/view.go |
| TestArchive_WithYes | cmd/cmd_test.go | cli:Archive Command | cmd/archive.go |
| TestArchive_SkipSpecs | cmd/cmd_test.go | cli:Archive Command | cmd/archive.go |
| TestArchive_NoValidate | cmd/cmd_test.go | cli:Archive Command | cmd/archive.go |
| TestArchive_NotFound | cmd/cmd_test.go | cli:Archive Command | cmd/archive.go |
| TestVersion | cmd/cmd_test.go | cli:Root Command | cmd/root.go |
| TestNoColor | cmd/cmd_test.go | cli:Root Command | cmd/root.go |
