package cmd

import (
	"fmt"
	"os"

	"github.com/IIIcecream/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "robot-cli",
	Long:    `机器人`,
	Version: version.FullVersion(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
