package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	tool "github.com/pawannn/netlite/internal/tools"
	"github.com/pawannn/netlite/pkg"
	"github.com/spf13/cobra"
)

var (
	scanHost        string
	scanStart       int
	scanEnd         int
	scanConcurrency int
	scanVerbose     bool
	scanJSON        string
	scanWatch       bool
	scanInterval    int
)

var ScannerCommand = &cobra.Command{
	Use:   "scan",
	Short: "Scan TCP ports on a host",
	Long:  "Scan TCP ports on a given host with options for concurrency, verbose output, JSON export, and continuous watch mode.",
	Run:   RunScan,
}

func RunScan(cmd *cobra.Command, args []string) {
	if scanHost == "" {
		host, err := pkg.GetIP()
		if err != pkg.NoErr {
			fmt.Fprintln(os.Stderr, err.ClientMessage)
			return
		}
		scanHost = host
	}

	run := func() {
		fmt.Printf("\nNetlite scanning %s ports %d-%d\n\n", scanHost, scanStart, scanEnd)
		progressCh := make(chan int)
		go func() {
			for p := range progressCh {
				if scanVerbose {
					fmt.Printf("Scanned port %d\n", p)
				}
			}
		}()
		results, err := tool.ScanRange(scanHost, scanStart, scanEnd, scanConcurrency, progressCh)
		close(progressCh)
		if err != pkg.NoErr {
			fmt.Fprintf(os.Stderr, "scan error: %v\n", err.ClientMessage)
			return
		}

		openCount := 0
		for _, r := range results {
			if r.Open {
				openCount++
				fmt.Printf("Port %d: OPEN (%s)\n", r.Port, r.Protocol)
			} else if scanVerbose {
				if r.Error != "" {
					fmt.Printf("Port %d: CLOSED — %s\n", r.Port, r.Error)
				} else {
					fmt.Printf("Port %d: CLOSED\n", r.Port)
				}
			}
		}

		fmt.Printf("\nTotal open ports: %d\n", openCount)

		if scanJSON != "" {
			file, err := os.Create(scanJSON)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create JSON file: %v\n", err)
				return
			}
			defer file.Close()
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(results); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write JSON: %v\n", err)
			} else {
				fmt.Printf("Results saved to %s\n", scanJSON)
			}
		}
	}

	if scanWatch {
		for {
			run()
			time.Sleep(time.Duration(scanInterval) * time.Second)
		}
	} else {
		run()
	}
}

func init() {
	ScannerCommand.Flags().StringVarP(&scanHost, "host", "", "", "Host to scan (default: localhost)")
	ScannerCommand.Flags().IntVarP(&scanStart, "start", "s", pkg.START_RANGE, "Start port (default 1)")
	ScannerCommand.Flags().IntVarP(&scanEnd, "end", "e", pkg.END_RANGE, "End port (default 1024)")
	ScannerCommand.Flags().IntVarP(&scanConcurrency, "concurrency", "c", pkg.CONCURRENCY, "Max concurrent connections (default 200)")
	ScannerCommand.Flags().BoolVarP(&scanVerbose, "verbose", "v", false, "Verbose output showing closed ports and errors")
	ScannerCommand.Flags().StringVarP(&scanJSON, "json", "j", "", "Export results to JSON file")
	ScannerCommand.Flags().BoolVarP(&scanWatch, "watch", "w", false, "Continuously scan every --interval seconds")
	ScannerCommand.Flags().IntVarP(&scanInterval, "interval", "i", pkg.INTERVAL, "Interval seconds for --watch (default 3)")
}
