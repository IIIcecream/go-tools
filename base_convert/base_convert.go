package base_convert

import (
	"fmt"
	"math/big"
	"strings"
)

// 定义进制的字符表，需要保证<=len(baseChars)进制
// int.go
// For bases <= 36, lower and upper case letters are considered the same:
// The letters 'a' to 'z' and 'A' to 'Z' represent digit values 10 to 35.
// For bases > 36, the upper case letters 'A' to 'Z' represent the digit
// values 36 to 61.
const (
	base36Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func getChar(n int64, base int) (byte, error) {
	if base <= 0 || base > len(base62Chars) {
		return 0, fmt.Errorf("invalid base: %d", base)
	}
	if base <= 36 {
		return base36Chars[n], nil
	}
	return base62Chars[n], nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// 将 input 从 inputBase 进制转成 outputBase 进制，如36进制->62进制
func baseConovert(input string, inputBase int, outputBase int) (string, error) {
	if inputBase <= 0 || inputBase > 62 {
		return "", fmt.Errorf("invalid inputBase: %d", inputBase)
	}
	if outputBase <= 0 || outputBase > 62 {
		return "", fmt.Errorf("invalid outputBase: %d", outputBase)
	}

	// 解析 inputBase 进制为大整数
	value := new(big.Int)
	_, ok := value.SetString(input, inputBase)
	if !ok {
		return "", fmt.Errorf("invalid input: %s", input)
	}

	// 转换为 outputBase 进制字符串
	baseDest := big.NewInt(int64(outputBase))
	var result strings.Builder
	zero := big.NewInt(0)
	remainder := new(big.Int)

	for value.Cmp(zero) > 0 {
		value.DivMod(value, baseDest, remainder) // value = value / outputBase, remainder = value % outputBase
		c, _ := getChar(remainder.Int64(), outputBase)
		result.WriteByte(c)
	}

	converted := result.String()
	return reverseString(converted), nil
}
