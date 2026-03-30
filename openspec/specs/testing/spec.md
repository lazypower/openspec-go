# testing

## Purpose

## Requirements

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

