package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "发送消息",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt32("id")
		if err != nil {
			log.Fatalln(err)
			return
		}
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			log.Fatalln(err)
			return
		}
		msg, err := cmd.Flags().GetString("msg")
		if err != nil {
			log.Fatalln(err)
			return
		}
		baseUrl := "http://apiin.im.baidu.com/api/msg/groupmsgsend"

		params := url.Values{}
		params.Add("access_token", token)
		fullURL := fmt.Sprintf("%s?%s", baseUrl, params.Encode())

		requestBody := map[string]interface{}{
			"message": map[string]interface{}{
				"header": map[string]interface{}{
					"toid": id,
				},
				"body": []map[string]string{
					{
						"type":    "TEXT",
						"content": msg,
					},
				},
			},
		}

		// 将请求体编码为 JSON
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatalln(err)
			return
		}

		// 创建 HTTP POST 请求
		req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalln(err)
			return
		}

		// 设置请求头
		req.Header.Set("Content-Type", "application/json")

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
			return
		}
		defer resp.Body.Close()

		// 检查响应状态码
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error sending message:", resp.Status)
			return
		}

		fmt.Println("Message sent successfully!")
	},
}

func init() {
	sendCmd.Flags().Int32P("id", "i", 5562261, "群聊id")
	sendCmd.Flags().StringP("token", "t", "d319b005726cce8b646a201999283d9eb", "机器人token")
	sendCmd.Flags().StringP("msg", "m", "", "发送的消息")
	rootCmd.AddCommand(sendCmd)
}
