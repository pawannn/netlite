package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/pawannn/netlite/internal/scanner"
	"github.com/pawannn/netlite/pkg"
	"github.com/spf13/cobra"
)

var (
	flagIP          string
	flagStart       int
	flagEnd         int
	flagWatch       bool
	flagInterval    int
	flagJSON        string
	flagKill        bool
	flagVerbose     bool
	flagConcurrency int
	flagNoBanner    bool
)

var rootCmd = &cobra.Command{
	Use:   "netlite",
	Short: "Netlite — lightweight port scanner",
	Long:  "Netlite scans TCP ports on localhost.",
	Run:   Netlite,
}

func Netlite(cmd *cobra.Command, args []string) {
	if flagNoBanner {
		pkg.PrintBanner()
	}

	var host string
	if flagIP != "" {
		if !pkg.IsValidIP(flagIP) {
			fmt.Fprintln(os.Stderr, "Invalid IP")
			os.Exit(1)
		}
		host = flagIP
	} else {
		var err pkg.NetliteErr
		host, err = pkg.GetIP()
		if err != pkg.NoErr {
			log.Fatalf("%s", err.ClientMessage)
		}
	}

	fmt.Printf("\nNetlite scanning %s ports %d-%d\n\n", host, flagStart, flagEnd)

	results, err := scanner.ScanRange(host, flagStart, flagEnd, flagConcurrency)

	if err != pkg.NoErr {
		fmt.Fprintf(os.Stderr, "scan error: %v\n", err.ClientMessage)
		os.Exit(1)
	}

	openCount := 0
	for _, r := range results {
		if r.Open {
			openCount++
			fmt.Printf("Port %d: OPEN (%s)\n", r.Port, r.Protocol)
		} else if flagVerbose {
			if r.Error != "" {
				fmt.Printf("Port %d: CLOSED — %s\n", r.Port, r.Error)
			} else {
				fmt.Printf("Port %d: CLOSED\n", r.Port)
			}
		}
	}
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(&flagIP, "IP", "", "Host to scan")
	rootCmd.PersistentFlags().IntVar(&flagStart, "start", 1, "Start port (default 1)")
	rootCmd.PersistentFlags().IntVar(&flagEnd, "end", 65535, "End port (default 1024)")
	rootCmd.PersistentFlags().BoolVar(&flagWatch, "watch", false, "Continuously scan every --interval seconds")
	rootCmd.PersistentFlags().IntVar(&flagInterval, "interval", 3, "Interval seconds for --watch (default 3)")
	rootCmd.PersistentFlags().StringVar(&flagJSON, "json", "", "Export results to JSON file (provide path)")
	rootCmd.PersistentFlags().BoolVar(&flagKill, "kill", false, "Attempt to kill process using a given port (requires sudo/admin)")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "Verbose output (show errors/timeouts)")
	rootCmd.PersistentFlags().IntVar(&flagConcurrency, "concurrency", 200, "Max concurrent connections (default 200)")
	rootCmd.PersistentFlags().BoolVar(&flagNoBanner, "no-banner", true, "Hides the banner")

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
