package cmd

import (
	"fmt"

	"github.com/IIIcecream/go-tools/password_gen/storage"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved passwords",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := storage.New()
		if err != nil {
			fmt.Printf("Failed to init: %v\n", err)
			return
		}

		entries, err := s.ListPasswords()
		if err != nil {
			fmt.Printf("Failed to list passwords: %v\n", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("No passwords stored")
			return
		}

		fmt.Println("Saved passwords:")
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
	rootCmd.AddCommand(listCmd)
}
