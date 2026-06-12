package projCommands

import (
	"github.com/spf13/cobra"
)

var ProjCmd = &cobra.Command{
	Use:   "proj",
	Short: "Manage projects",
	Long:  "Create, list, and manage projects.",
}
