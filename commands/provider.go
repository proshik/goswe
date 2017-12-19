package commands

import "github.com/spf13/cobra"

var provider = &cobra.Command{
	Use:   "provider []",
	Short: "open translate screen",
	Run: translateRun,
}