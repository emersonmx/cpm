package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Run managers update",
	Long:    "Run managers update",
	Example: "stm update arch rust pip",
	Args:    ValidManagers,
	Run: func(cmd *cobra.Command, args []string) {
		for _, a := range args {
			m, err := config.FindManager(a)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			pkgs := []string{}
			for _, t := range config.Tools.FilterByManager(m.Name) {
				pkgs = append(pkgs, t.Package)
			}

			if err := m.UpdatePackages(pkgs); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
