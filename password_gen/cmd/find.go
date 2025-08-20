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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(find_key) == 0 {
			return fmt.Errorf("--key is empty")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		s, err := storage.New()
		if err != nil {
			fmt.Printf("Failed to init: %v\n", err)
			return
		}

		entry, err := s.FindPassword(find_key)
		if err != nil {
			fmt.Printf("Failed to list passwords: %v\n", err)
			return
		}
		fmt.Printf("password is : %s\n", entry.Password)
	},
}

func init() {
	findCmd.Flags().StringVarP(&find_key, "key", "k", "", "the password key to find")

	rootCmd.AddCommand(findCmd)
}
