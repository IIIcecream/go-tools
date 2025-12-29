package cmd

import (
	"fmt"
	"os"

	"github.com/IIIcecream/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "log_parser",
	Short: "parse log, add timestamp",
	Long: `If you want to parse log.zip once and then exit, use [log_parser parse].
If you want to listen dir and watch new zip all the time, use [log_parser listen]
If you just want to merge some log files, use [log_parser merge]`,
	Example: `log_parser listen -s D:\\log_parse -m APP_PROXY,EVENT_TASK`,
	Version: version.FullVersion(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
