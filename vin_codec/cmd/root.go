package cmd

import (
	"fmt"
	"os"

	"github.com/IIIcecream/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "vin_cli",
	Long:    `vin号加密解密工具`,
	Version: version.FullVersion(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
