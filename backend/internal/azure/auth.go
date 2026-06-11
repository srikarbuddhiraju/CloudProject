// Package azure provides Azure-specific authentication and collectors.
package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// ARMScope is the OAuth scope for Azure Resource Manager (Azure Public Cloud).
const ARMScope = "https://management.azure.com/.default"

// Authenticator provides an Azure credential for collectors.
//
// A small per-cloud interface (rather than one shared across clouds) so that
// future AWS/GCP collectors can define their own native credential chains
// without constraining this one.
type Authenticator interface {
	Credential() (azcore.TokenCredential, error)
}

// DefaultAuthenticator resolves credentials via azidentity's default chain:
// environment variables, Azure CLI / Azure Developer CLI session, managed
// identity, and so on, in order.
type DefaultAuthenticator struct{}

func (DefaultAuthenticator) Credential() (azcore.TokenCredential, error) {
	return azidentity.NewDefaultAzureCredential(nil)
}
