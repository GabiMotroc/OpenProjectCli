package appCommands

import (
	"github.com/spf13/cobra"
)

var AppCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage apps",
	Long:  "Create, list, and manage applications.",
}
