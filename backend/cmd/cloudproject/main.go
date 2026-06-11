// Command cloudproject is the entry point for the CloudProject CLI.
package main

import (
	"os"

	"github.com/srikarbuddhiraju/CloudProject/backend/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
