package projCommands

import (
	"OpenCli/data"

	"github.com/spf13/cobra"
)

func clearApps(c *cobra.Command, args []string) {
	err := data.SaveApps([]data.App{})
	if err != nil {
		return
	}
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear projects",
	Long:  `Clear all projects from the system. This command will remove all projects data and configurations.`,
	Run:   clearApps,
}

func init() {
	ProjCmd.AddCommand(clearCmd)
}
