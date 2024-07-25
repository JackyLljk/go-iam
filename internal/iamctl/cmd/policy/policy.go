package policy

import (
	"github.com/spf13/cobra"

	cmdutil "j-iam/internal/iamctl/cmd/util"
	"j-iam/internal/iamctl/util/templates"
	"j-iam/pkg/cli/genericclioptions"
)

var policyLong = templates.LongDesc(`
	Authorization policy management commands.

	This commands allow you to manage your authorization policy on iam platform.`)

// NewCmdPolicy returns new initialized instance of 'policy' sub command.
func NewCmdPolicy(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "policy SUBCOMMAND",
		DisableFlagsInUseLine: true,
		Short:                 "Manage authorization policies on iam platform",
		Long:                  policyLong,
		Run:                   cmdutil.DefaultSubCommandRun(ioStreams.ErrOut),
	}

	cmd.AddCommand(NewCmdCreate(f, ioStreams))
	cmd.AddCommand(NewCmdGet(f, ioStreams))
	cmd.AddCommand(NewCmdList(f, ioStreams))
	cmd.AddCommand(NewCmdDelete(f, ioStreams))
	cmd.AddCommand(NewCmdUpdate(f, ioStreams))

	return cmd
}
