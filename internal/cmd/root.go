package cmd

import (
	"github.com/chuck/openspec-go/internal/output"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root openspec command.
func NewRootCmd(version string) *cobra.Command {
	var noColor bool

	cmd := &cobra.Command{
		Use:     "openspec",
		Short:   "Spec-driven development tool",
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if noColor {
				output.SetNoColor(true)
			}
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	cmd.AddCommand(
		newInitCmd(),
		newUpdateCmd(),
		newListCmd(),
		newShowCmd(),
		newValidateCmd(),
		newArchiveCmd(),
		newViewCmd(),
	)

	return cmd
}
