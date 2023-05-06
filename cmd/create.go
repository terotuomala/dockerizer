package cmd

import (
	"github.com/spf13/cobra"
	"github.com/terotuomala/dockerizer/pkg/prompt"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Dockerfile for Go or Node.js application",
	Run: func(cmd *cobra.Command, args []string) {
		prompt.StartPrompt()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
