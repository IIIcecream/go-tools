package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	encodeString string
	decodeString string
)

var rootCmd = &cobra.Command{
	Use:  "vin_cli",
	Long: `vin号加密解密工具`,
}

func Execute() {
	rootCmd.Flags().StringVarP(&encodeString, "string", "e", "L6T79P4N7RD091732", "vin")
	rootCmd.Flags().StringVarP(&decodeString, "string", "d", "dB4dgIsvjWgbZLw", "encoded vin")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
