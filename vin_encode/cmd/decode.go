package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/IIIcecream/go-tools/convert"
)

func decode(output string) string {
	ret, _ := convert.baseConovert(output, 62, 36)
	return ret
}

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode vin",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please input the string you want to decode")
			return
		}
		fmt.Println(decode(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)
}
