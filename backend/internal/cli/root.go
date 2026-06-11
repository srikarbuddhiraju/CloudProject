// Package cli wires up the CloudProject command-line interface.
package cli

import (
	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags "-X .../cli.version=...".
var version = "dev"

// NewRootCmd builds the root "cloudproject" command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "cloudproject",
		Short: "Self-hosted Azure Landing Zone conformance dashboard",
		Long: `cloudproject discovers your Azure tenant (management groups, subscriptions,
policy assignments, RBAC, network topology, cost data) using a Reader-only
service principal, and evaluates it against the Azure Landing Zone (ALZ)
reference architecture.`,
	}

	root.AddCommand(newVersionCmd())

	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the cloudproject version",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(version)
			return nil
		},
	}
}
