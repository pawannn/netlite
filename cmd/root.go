package cmd

import (
	"fmt"
	"os"

	"github.com/pawannn/netlite/pkg"
	"github.com/spf13/cobra"
)

var FlagNoBanner bool

var rootCmd = &cobra.Command{
	Use:              "netlite",
	Short:            "Netlite — lightweight port scanner & utilities",
	Long:             "Netlite — lightweight port scanner & utilities",
	PersistentPreRun: DisplayBanner,
}

func DisplayBanner(cmd *cobra.Command, args []string) {
	pkg.DisplayBanner()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().BoolVar(&FlagNoBanner, "no-banner", false, "")

	rootCmd.AddCommand(ScannerCommand)
	rootCmd.AddCommand(ipinfoCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpTemplate(pkg.HelpTemplate)
}
