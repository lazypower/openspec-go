package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	var (
		itemType     string
		jsonOutput   bool
		deltasOnly   bool
		requirements bool
		reqIndex     int
	)

	cmd := &cobra.Command{
		Use:   "show <item>",
		Short: "Display a change or spec",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ospPath, err := findOpenSpecPath("")
			if err != nil {
				return err
			}

			itemID := args[0]
			resolvedType := itemType

			if resolvedType == "" {
				resolvedType = autoDetectType(ospPath, itemID)
			}

			if resolvedType == "" {
				// Suggest nearest match
				var candidates []string
				if entries, err := os.ReadDir(filepath.Join(ospPath, "changes")); err == nil {
					for _, e := range entries {
						if e.IsDir() && e.Name() != "archive" {
							candidates = append(candidates, e.Name())
						}
					}
				}
				if entries, err := os.ReadDir(filepath.Join(ospPath, "specs")); err == nil {
					for _, e := range entries {
						if e.IsDir() {
							candidates = append(candidates, e.Name())
						}
					}
				}
				suggestion := suggestNearestMatch(itemID, candidates)
				if suggestion != "" {
					return fmt.Errorf("item %q not found. Did you mean %q?", itemID, suggestion)
				}
				return fmt.Errorf("item %q not found", itemID)
			}

			if resolvedType == "change" {
				return showChange(ospPath, itemID, jsonOutput, deltasOnly)
			}
			return showSpec(ospPath, itemID, jsonOutput, requirements, reqIndex)
		},
	}

	cmd.Flags().StringVar(&itemType, "type", "", "Force item type: change or spec")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&deltasOnly, "deltas-only", false, "Show only deltas (changes, JSON mode)")
	cmd.Flags().BoolVar(&requirements, "requirements", false, "Show requirements without scenarios (specs, JSON mode)")
	cmd.Flags().IntVarP(&reqIndex, "requirement", "r", 0, "Show single requirement by 1-based index (specs, JSON mode)")
	return cmd
}

func autoDetectType(ospPath, id string) string {
	changePath := filepath.Join(ospPath, "changes", id)
	specPath := filepath.Join(ospPath, "specs", id)

	changeExists := dirExists(changePath)
	specExists := dirExists(specPath)

	if changeExists && !specExists {
		return "change"
	}
	if specExists && !changeExists {
		return "spec"
	}
	if changeExists && specExists {
		return "change" // default to change when ambiguous
	}
	return ""
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func showChange(ospPath, id string, jsonOut, deltasOnly bool) error {
	if !jsonOut {
		data, err := os.ReadFile(filepath.Join(ospPath, "changes", id, "proposal.md"))
		if err != nil {
			return err
		}
		fmt.Print(string(data))
		return nil
	}

	change, err := loadChange(ospPath, id)
	if err != nil {
		return err
	}

	if deltasOnly {
		return outputJSON(change.Deltas)
	}

	type jsonChange struct {
		ID         string `json:"id"`
		Title      string `json:"title"`
		DeltaCount int    `json:"deltaCount"`
		Deltas     any    `json:"deltas"`
	}
	return outputJSON(jsonChange{
		ID:         change.ID,
		Title:      change.Title,
		DeltaCount: len(change.Deltas),
		Deltas:     change.Deltas,
	})
}

func showSpec(ospPath, id string, jsonOut, requirementsOnly bool, reqIndex int) error {
	if !jsonOut {
		data, err := os.ReadFile(filepath.Join(ospPath, "specs", id, "spec.md"))
		if err != nil {
			return err
		}
		fmt.Print(string(data))
		return nil
	}

	specs, _ := loadSpecs(ospPath)
	spec, ok := specs[id]
	if !ok {
		return fmt.Errorf("spec %q not found", id)
	}

	if reqIndex > 0 {
		if reqIndex > len(spec.Requirements) {
			return fmt.Errorf("requirement index %d out of range (spec has %d requirements)", reqIndex, len(spec.Requirements))
		}
		return outputJSON(spec.Requirements[reqIndex-1])
	}

	if requirementsOnly {
		type briefReq struct {
			Name string `json:"name"`
			Text string `json:"text"`
		}
		var reqs []briefReq
		for _, r := range spec.Requirements {
			reqs = append(reqs, briefReq{Name: r.Name, Text: r.Text})
		}
		return outputJSON(reqs)
	}

	type jsonSpec struct {
		ID               string `json:"id"`
		Title            string `json:"title"`
		Overview         string `json:"overview"`
		RequirementCount int    `json:"requirementCount"`
		Requirements     any    `json:"requirements"`
	}
	return outputJSON(jsonSpec{
		ID:               id,
		Title:            spec.Title,
		Overview:         spec.Overview,
		RequirementCount: len(spec.Requirements),
		Requirements:     spec.Requirements,
	})
}

func outputJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
