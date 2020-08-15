package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "Run managers install",
	Long:    "Run managers install",
	Example: "stm install arch rust pip",
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
				if t.Installed() {
					continue
				}
				pkgs = append(pkgs, t.Package)
			}

			if err := m.InstallPackages(pkgs); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
