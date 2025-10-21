package cmd

import (
	"fmt"

	"github.com/pawannn/netlite/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Netlite version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Netlite", version.Version)
	},
}
