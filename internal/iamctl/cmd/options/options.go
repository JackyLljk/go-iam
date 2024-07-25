package options

import (
	"io"

	"github.com/spf13/cobra"

	"j-iam/internal/iamctl/util/templates"
)

var optionsExample = templates.Examples(`
		# Print flags inherited by all commands
		iamctl options`)

// NewCmdOptions implements the options command.
func NewCmdOptions(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "options",
		Short:   "Print the list of flags inherited by all commands",
		Long:    "Print the list of flags inherited by all commands",
		Example: optionsExample,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	// The `options` command needs write its output to the `out` stream
	// (typically stdout). Without calling SetOutput here, the Usage()
	// function call will fall back to stderr.
	cmd.SetOutput(out)

	templates.UseOptionsTemplates(cmd)

	return cmd
}
