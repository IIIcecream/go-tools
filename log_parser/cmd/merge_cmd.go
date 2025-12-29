package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge same module log files",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("todo")
	},
}

func init() {
	mergeCmd.Flags().StringP("srcDir", "s", ".", "src log file path")
	mergeCmd.Flags().StringP("destDir", "d", ".", "dest path")
	mergeCmd.Flags().StringSliceP(
		"modules",
		"m",
		[]string{},
		"log modules to parse (e.g. app_proxy,eventtask)",
	)
	rootCmd.AddCommand(mergeCmd)
}
