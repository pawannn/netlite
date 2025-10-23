package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pawannn/netlite/internal/tools"
	"github.com/spf13/cobra"
)

var killPorts string

var KillCommand = &cobra.Command{
	Use:   "kill",
	Short: "Kill a process running on a port",
	Long:  "Kill processes that are listening on specified TCP ports. Requires sudo/admin privileges.",
	Run:   runKillCommand,
}

func init() {
	KillCommand.Flags().StringVarP(&killPorts, "ports", "p", "", "Comma-separated list of ports to kill (required)")
	KillCommand.MarkFlagRequired("ports")
}

func runKillCommand(cmd *cobra.Command, args []string) {
	if killPorts == "" {
		fmt.Fprintln(os.Stderr, "Please specify ports to kill using --ports")
		return
	}

	// Parse ports string into []int
	portsStr := strings.Split(killPorts, ",")
	var ports []int
	for _, ps := range portsStr {
		ps = strings.TrimSpace(ps)
		p, err := strconv.Atoi(ps)
		if err != nil || p < 1 || p > 65535 {
			fmt.Fprintf(os.Stderr, "Invalid port: %s\n", ps)
			return
		}
		ports = append(ports, p)
	}

	// Call the tools function
	if err := tools.KillOpenPorts(ports); err != nil {
		fmt.Printf("Error killing ports: %v\n", err)
	} else {
		fmt.Println("Ports killed successfully")
	}
}
