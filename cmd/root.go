package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/emersonmx/stm/internal/app/stm"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "stm",
		Short:   "STM is a System Tool Manager",
		Long:    "STM is a System Tool Manager",
		Version: "1.0.2",
	}
	config stm.Config
)

func Execute() {
	var err error
	config, err = stm.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err = rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ValidManagers(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires a managers argument")
	}

	for _, a := range args {
		valid := false
		for _, m := range config.Managers.Names() {
			if a == m {
				valid = true
				break
			}
		}
		if !valid {
			return errors.New(fmt.Sprintf("invalid manager %s", a))
		}
	}

	return nil
}
