package cmd

import "github.com/spf13/cobra"

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "输入消息",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	sendCmd.Flags().Int32P("id", "i", 5562261, "群聊id")
	sendCmd.Flags().StringP("token", "t", "d319b005726cce8b646a201999283d9eb", "机器人token")
	sendCmd.Flags().StringP("msg", "m", "", "发送的消息")
	rootCmd.AddCommand(receiveCmd)
}
