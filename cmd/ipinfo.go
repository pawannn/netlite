package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/pawannn/netlite/internal/tools"
	"github.com/pawannn/netlite/pkg"
	"github.com/spf13/cobra"
)

var ipinfoCmd = &cobra.Command{
	Use:   "ipinfo",
	Short: "Display IP information for all local network interfaces",
	Run:   runIPInfo,
}

func runIPInfo(cmd *cobra.Command, args []string) {
	infos, err := tools.IfInfo()
	if err != pkg.NoErr {
		fmt.Println(err.ClientMessage)
		return
	}

	heading := color.New(color.FgHiCyan, color.Bold)
	label := color.New(color.FgHiWhite)
	value := color.New(color.FgHiGreen)

	fmt.Println()
	heading.Println("Network Interfaces:")
	fmt.Println("────────────────────")

	for _, inf := range infos {
		label.Printf("Name: ")
		value.Printf("%s\n", inf.Name)
		label.Printf("  Type: ")
		value.Printf("%s\n", inf.Type)

		if len(inf.IPv4) > 0 {
			label.Printf("  IPv4: ")
			fmt.Println("  ", value.Sprint(inf.IPv4))
		}
		if len(inf.IPv6) > 0 {
			label.Printf("  IPv6: ")
			fmt.Println("  ", value.Sprint(inf.IPv6))
		}
		label.Printf("  MAC: ")
		value.Printf("%s\n", inf.MAC)
		label.Printf("  Status: ")
		if inf.IsUp {
			value.Println("UP")
		} else {
			color.New(color.FgHiRed).Println("DOWN")
		}
		fmt.Println()
	}
}
