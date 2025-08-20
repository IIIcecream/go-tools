package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/IIIcecream/go-tools/password_gen/storage"
	"github.com/spf13/cobra"
)

var (
	length          int
	min             int
	special         bool
	special_charset string
	lower           bool
	lower_charset   string
	upper           bool
	upper_charset   string
	number          bool
	number_charset  string
	save            bool
	key             string
)

func create_random_password() string {
	var allCharsets []string
	var requiredCharsets [][]rune
	if special {
		allCharsets = append(allCharsets, special_charset)
		requiredCharsets = append(requiredCharsets, []rune(special_charset))
	}
	if lower {
		allCharsets = append(allCharsets, lower_charset)
		requiredCharsets = append(requiredCharsets, []rune(lower_charset))
	}
	if upper {
		allCharsets = append(allCharsets, upper_charset)
		requiredCharsets = append(requiredCharsets, []rune(upper_charset))
	}
	if number {
		allCharsets = append(allCharsets, number_charset)
		requiredCharsets = append(requiredCharsets, []rune(number_charset))
	}
	// 字符集
	fullCharset := strings.Join(allCharsets, "")

	// 生成密码
	password := make([]byte, length)
	fullCharsetRunes := []rune(fullCharset)

	// 1. 首先确保每个启用的字符集至少有一个字符出现在密码中
	for i, charset := range requiredCharsets {
		// 从该字符集中随机选择一个字符
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}
		password[i] = byte(charset[idx.Int64()])
	}

	// 2. 填充剩余位置
	for i := len(requiredCharsets); i < length; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(fullCharsetRunes))))
		if err != nil {
			return ""
		}
		password[i] = byte(fullCharsetRunes[idx.Int64()])
	}

	// 3. 打乱密码中的字符顺序（避免前面总是按特定字符集顺序排列）
	for i := range password {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {
			return ""
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	ret := string(password)
	if save {
		storage, err := storage.New()
		if err != nil {
			fmt.Printf("Failed to initialize password storage: %v\n", err)
			return ""
		}

		if err := storage.SavePassword(key, ret); err != nil {
			fmt.Printf("Failed to save password: %v\n", err)
			return ""
		}
	}
	return ret
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create random password",
	Long: `Create a secure random password with customizable character sets.

This command generates cryptographically secure random passwords with the following features:

- Specify exact password length (required)
- Include/exclude character types:
  * Special characters (!@#$%^&* etc.)
  * Lowercase letters (a-z)
  * Uppercase letters (A-Z)
  * Numbers (0-9)
  
- Customize character sets for each type
- Optionally save password with encryption key

Examples:
  # Basic 12-character password with all character types
  passgen create --length 12 --special --lower --upper --number
  
  # Custom special characters only
  passgen create --length 16 --special --special_charset "@#$%"
  
  # Save password with encryption key
  passgen create --length 24 --save --key "my-secret-key"`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if length == 0 {
			return fmt.Errorf("--length must be greater than zero")
		}
		if special && len(special_charset) == 0 {
			return fmt.Errorf("--special_charset must be not empty when --special is true")
		}
		if special && len(special_charset) == 0 {
			return fmt.Errorf("--special_charset must be not empty when --special is true")
		}
		if lower && len(lower_charset) == 0 {
			return fmt.Errorf("--lower_charset must be not empty when --lower is true")
		}
		if upper && len(upper_charset) == 0 {
			return fmt.Errorf("--upper_charset must be not empty when --upper is true")
		}
		if number && len(number_charset) == 0 {
			return fmt.Errorf("--number_charset must be not empty when --number is true")
		}
		if save && len(key) == 0 {
			return fmt.Errorf("--key must be set when --save is true")
		}
		cnt := 0
		if special {
			cnt += 1
		}
		if lower {
			cnt += 1
		}
		if upper {
			cnt += 1
		}
		if number {
			cnt += 1
		}
		if cnt < min {
			return fmt.Errorf("at least %d character type must be enabled (--special, --lower, --upper or --number)", min)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("password is ", create_random_password())
	},
}

func init() {
	createCmd.Flags().IntVar(&length, "length", 10, "password length")
	createCmd.Flags().IntVarP(&min, "min", "m", 2, "at least [min] character type must be enabled")

	createCmd.Flags().BoolVarP(&special, "special", "s", false, "contain special char")
	createCmd.Flags().StringVar(&special_charset, "special_charset", "!@#$%^&*()_+-=[]{};':\",./<>?", "special charset")

	createCmd.Flags().BoolVarP(&lower, "lower", "l", true, "contain lower char")
	createCmd.Flags().StringVar(&lower_charset, "lower_charset", "abcdefghijklmnopqrstuvwxyz", "lower charset")

	createCmd.Flags().BoolVarP(&upper, "upper", "u", false, "contain upper char")
	createCmd.Flags().StringVar(&upper_charset, "upper_charset", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "upper charset")

	createCmd.Flags().BoolVarP(&number, "number", "n", true, "contain number")
	createCmd.Flags().StringVar(&number_charset, "number_charset", "0123456789", "number charset")

	createCmd.Flags().BoolVar(&save, "save", false, "save password with key")
	createCmd.Flags().StringVar(&key, "key", "0123456789", "number charset")

	rootCmd.AddCommand(createCmd)
}
