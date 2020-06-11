package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "notice",
}

func Execute() error {
	return rootCmd.Execute()
}
