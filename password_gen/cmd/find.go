package cmd

import (
	"fmt"

	"github.com/IIIcecream/go-tools/password_gen/storage"
	"github.com/spf13/cobra"
)

var (
	find_key string
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "find saved passwords by key",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := storage.New()
		if err != nil {
			fmt.Printf("Failed to init: %v\n", err)
			return
		}

		entries, err := s.FindPasswords(find_key)
		if err != nil {
			fmt.Printf("Failed to list passwords: %v\n", err)
			return
		}
		for i, entry := range entries {
			fmt.Printf("%d.\n", i+1)
			fmt.Printf("   Key: %s\n", entry.Key)
			fmt.Printf("   Password: %s\n", entry.Password)
			fmt.Printf("   Created: %s\n", entry.Timestamp)
			fmt.Println()
		}
	},
}

func init() {
	findCmd.Flags().StringVarP(&find_key, "key", "k", "", "the password key to find")

	rootCmd.AddCommand(findCmd)
}
