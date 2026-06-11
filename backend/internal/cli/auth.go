package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/spf13/cobra"

	"github.com/srikarbuddhiraju/CloudProject/backend/internal/azure"
)

func newAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
	}
	authCmd.AddCommand(newAuthCheckCmd())
	return authCmd
}

func newAuthCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Verify Azure credentials can be resolved and used to get a token",
		Long: `Resolves Azure credentials via the default credential chain (environment
variables, Azure CLI / Azure Developer CLI session, managed identity, etc.)
and requests a token for Azure Resource Manager, to confirm the
Reader-only service principal (or other credential) is set up correctly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			auth := azure.DefaultAuthenticator{}

			cred, err := auth.Credential()
			if err != nil {
				return fmt.Errorf("resolving credential: %w", err)
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
			defer cancel()

			token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
				Scopes: []string{azure.ARMScope},
			})
			if err != nil {
				return fmt.Errorf("getting token: %w", err)
			}

			cmd.Printf("OK — token acquired, expires %s\n", token.ExpiresOn.Format(time.RFC3339))
			return nil
		},
	}
}
