package cmd

import (
	"fmt"

	"github.com/pawannn/netlite/internal/tools"
	"github.com/spf13/cobra"
)

var killPorts string

var KillCommand = &cobra.Command{
	Use:   "kill",
	Short: "Kill a process running on a port",
	Run:   KillPort,
}

func KillPort(cmd *cobra.Command, args []string) {
	ports := []int{8080, 3000}
	err := tools.KillOpenPorts(ports)
	if err != nil {
		fmt.Printf("Error killing ports: %v\n", err)
	} else {
		fmt.Println("Ports killed successfully")
	}
}

func init() {
	KillCommand.PersistentFlags().StringVarP(&killPorts, "ports", "p", "", "Comma-separated list of ports to kill")
}
