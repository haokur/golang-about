package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "MyApp is a CLI tool",
	Long:  `MyApp is a CLI tool to demonstrate the usage of Cobra`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to MyApp!")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
