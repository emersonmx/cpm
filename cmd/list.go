package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available managers",
	Long:  "List all available managers",
	Run: func(cmd *cobra.Command, args []string) {
		for _, m := range config.Managers {
			fmt.Println(m.Name)
			for _, t := range config.Tools {
				if t.Manager != m.Name {
					continue
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
