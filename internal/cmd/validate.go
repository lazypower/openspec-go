package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/chuck/openspec-go/internal/model"
	"github.com/chuck/openspec-go/internal/output"
	"github.com/chuck/openspec-go/internal/validator"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	var (
		strict      bool
		jsonOutput  bool
		all         bool
		changes     bool
		specs       bool
		concurrency int
		itemType    string
	)

	cmd := &cobra.Command{
		Use:   "validate [item]",
		Short: "Check correctness of changes and specs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ospPath, err := findOpenSpecPath("")
			if err != nil {
				return err
			}

			// Determine concurrency
			if concurrency <= 0 {
				if env := os.Getenv("OPENSPEC_CONCURRENCY"); env != "" {
					if n, err := strconv.Atoi(env); err == nil && n > 0 {
						concurrency = n
					}
				}
				if concurrency <= 0 {
					concurrency = runtime.NumCPU()
				}
			}

			if len(args) > 0 {
				return validateSingle(ospPath, args[0], itemType, strict, jsonOutput)
			}
			if all {
				return validateAll(ospPath, strict, jsonOutput, concurrency)
			}
			if specs {
				return validateSpecs(ospPath, strict, jsonOutput, concurrency)
			}
			// Default: validate changes
			return validateChanges(ospPath, strict, jsonOutput, concurrency)
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "Enable strict validation mode")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON report")
	cmd.Flags().BoolVar(&all, "all", false, "Validate all changes and specs")
	cmd.Flags().BoolVar(&changes, "changes", false, "Validate changes only")
	cmd.Flags().BoolVar(&specs, "specs", false, "Validate specs only")
	cmd.Flags().IntVar(&concurrency, "concurrency", 0, "Maximum concurrent validations")
	cmd.Flags().StringVar(&itemType, "type", "", "Force item type: change or spec")
	return cmd
}

func validateSingle(ospPath, id, forceType string, strict, jsonOut bool) error {
	itemType := forceType
	if itemType == "" {
		itemType = autoDetectType(ospPath, id)
	}

	var item model.ValidationItem
	switch itemType {
	case "change":
		change, err := loadChange(ospPath, id)
		if err != nil {
			return err
		}
		issues := validator.ValidateChange(change)
		if strict {
			allChanges, _ := loadChanges(ospPath)
			allSpecs, _ := loadSpecs(ospPath)
			issues = append(issues, validator.ValidateStrict(allChanges, allSpecs)...)
		}
		item = toValidationItem(id, "change", issues)
	case "spec":
		allSpecs, _ := loadSpecs(ospPath)
		spec, ok := allSpecs[id]
		if !ok {
			return fmt.Errorf("spec %q not found", id)
		}
		issues := validator.ValidateSpec(spec, id)
		item = toValidationItem(id, "spec", issues)
	default:
		return fmt.Errorf("item %q not found", id)
	}

	report := buildReport([]model.ValidationItem{item})
	return outputReport(report, jsonOut)
}

func validateChanges(ospPath string, strict, jsonOut bool, concurrency int) error {
	allChanges, err := loadChanges(ospPath)
	if err != nil {
		return err
	}

	var items []validator.ValidationFunc
	for _, ch := range allChanges {
		ch := ch
		items = append(items, func() (model.ValidationItem, error) {
			issues := validator.ValidateChange(ch)
			return toValidationItem(ch.ID, "change", issues), nil
		})
	}

	results := validator.ValidateConcurrent(items, concurrency)

	if strict {
		allSpecs, _ := loadSpecs(ospPath)
		strictIssues := validator.ValidateStrict(allChanges, allSpecs)
		if len(strictIssues) > 0 {
			results = append(results, toValidationItem("strict", "meta", strictIssues))
		}
	}

	report := buildReport(results)
	return outputReport(report, jsonOut)
}

func validateSpecs(ospPath string, strict, jsonOut bool, concurrency int) error {
	allSpecs, err := loadSpecs(ospPath)
	if err != nil {
		return err
	}

	var items []validator.ValidationFunc
	for name, spec := range allSpecs {
		name, spec := name, spec
		items = append(items, func() (model.ValidationItem, error) {
			issues := validator.ValidateSpec(spec, name)
			return toValidationItem(name, "spec", issues), nil
		})
	}

	results := validator.ValidateConcurrent(items, concurrency)
	report := buildReport(results)
	return outputReport(report, jsonOut)
}

func validateAll(ospPath string, strict, jsonOut bool, concurrency int) error {
	allChanges, _ := loadChanges(ospPath)
	allSpecs, _ := loadSpecs(ospPath)

	var items []validator.ValidationFunc
	for _, ch := range allChanges {
		ch := ch
		items = append(items, func() (model.ValidationItem, error) {
			issues := validator.ValidateChange(ch)
			return toValidationItem(ch.ID, "change", issues), nil
		})
	}
	for name, spec := range allSpecs {
		name, spec := name, spec
		items = append(items, func() (model.ValidationItem, error) {
			issues := validator.ValidateSpec(spec, name)
			return toValidationItem(name, "spec", issues), nil
		})
	}

	results := validator.ValidateConcurrent(items, concurrency)

	if strict {
		strictIssues := validator.ValidateStrict(allChanges, allSpecs)
		if len(strictIssues) > 0 {
			results = append(results, toValidationItem("strict", "meta", strictIssues))
		}
	}

	report := buildReport(results)
	return outputReport(report, jsonOut)
}

func toValidationItem(id, itemType string, issues []model.Issue) model.ValidationItem {
	valid := true
	for _, iss := range issues {
		if iss.Level == model.LevelError {
			valid = false
			break
		}
	}
	return model.ValidationItem{
		ID:     id,
		Type:   itemType,
		Valid:  valid,
		Issues: issues,
	}
}

func buildReport(items []model.ValidationItem) model.ValidationReport {
	report := model.ValidationReport{
		Items: items,
		Summary: model.ValidationSummary{
			ByType: make(map[string]int),
		},
	}

	for _, item := range items {
		report.Summary.Totals.Items++
		if item.Valid {
			report.Summary.Totals.Passed++
		} else {
			report.Summary.Totals.Failed++
		}
		report.Summary.ByType[item.Type]++
	}

	return report
}

func outputReport(report model.ValidationReport, jsonOut bool) error {
	if jsonOut {
		return outputJSON(report)
	}

	// Text output to stderr
	for _, item := range report.Items {
		if len(item.Issues) == 0 {
			fmt.Fprintf(os.Stderr, "%s %s: %s\n", output.Green("✓"), item.ID, "valid")
			continue
		}
		for _, iss := range item.Issues {
			prefix := output.ErrorStyle.Render(string(iss.Level))
			if iss.Level == model.LevelWarning {
				prefix = output.WarnStyle.Render(string(iss.Level))
			}
			fmt.Fprintf(os.Stderr, "%s [%s] %s\n", prefix, item.ID, iss.Message)
		}
	}

	fmt.Fprintf(os.Stderr, "\n%d items: %d passed, %d failed\n",
		report.Summary.Totals.Items,
		report.Summary.Totals.Passed,
		report.Summary.Totals.Failed,
	)

	if report.HasErrors() {
		return fmt.Errorf("validation failed")
	}
	return nil
}
