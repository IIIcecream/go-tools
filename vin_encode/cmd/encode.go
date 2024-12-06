package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/IIIcecream/go-tools/base_convert"
)

func encode(input string) string {
	ret, _ := base_convert.baseConovert(input, 36, 62)
	return ret
}

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode vin",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please input the string you want to encode")
			return
		}
		fmt.Println(encode(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)
}
