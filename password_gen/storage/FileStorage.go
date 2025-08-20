package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	passwordStoreDir  = ".passgen"      // 存储目录名
	passwordStoreFile = "passwords.enc" // 存储文件名
	encryptionKey     = "tangwei"
)

type PasswordEntry struct {
	Key       string `json:"key"`
	Password  string `json:"password"`
	Timestamp string `json:"timestamp"`
}

type FileStorage struct {
	storagePath string
}

func New() (FileStorage, error) {
	fs := FileStorage{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fs, fmt.Errorf("failed to get user home directory: %v", err)
	}

	storageDir := filepath.Join(homeDir, passwordStoreDir)
	if err := os.MkdirAll(storageDir, 0700); err != nil {
		return fs, fmt.Errorf("failed to create storage directory: %v", err)
	}

	fs.storagePath = storageDir + "/" + passwordStoreFile
	return fs, nil
}

func deriveKey(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// SavePassword 保存加密后的密码
func (ps *FileStorage) SavePassword(key, value string) error {
	entry := PasswordEntry{
		Key:       key,
		Password:  value,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal password entry: %v", err)
	}

	derivedKey := deriveKey(encryptionKey)
	encryptedData, err := encrypt(jsonData, derivedKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %v", err)
	}

	encodedData := base64.StdEncoding.EncodeToString(encryptedData)

	// 追加到文件或创建新文件
	f, err := os.OpenFile(ps.storagePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open password store file: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(encodedData + "\n"); err != nil {
		return fmt.Errorf("failed to write to password store: %v", err)
	}

	return nil
}

func (ps *FileStorage) getPasswords(filterKey string) ([]PasswordEntry, error) {
	derivedKey := deriveKey(encryptionKey)

	data, err := os.ReadFile(ps.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []PasswordEntry{}, nil
		}
		return nil, fmt.Errorf("failed to read password store: %v", err)
	}

	var entries []PasswordEntry
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		encryptedData, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return nil, fmt.Errorf("failed to decode password entry: %v", err)
		}

		decryptedData, err := decrypt(encryptedData, derivedKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt password entry: %v", err)
		}

		var entry PasswordEntry
		if err := json.Unmarshal(decryptedData, &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal password entry: %v", err)
		}

		if filterKey == "" {
			entries = append(entries, entry)
		} else if filterKey == entry.Key {
			entries = append(entries, entry)
			break
		}
	}

	return entries, nil
}

func (ps *FileStorage) FindPassword(key string) (PasswordEntry, error) {
	entries, e := ps.getPasswords(key)
	if len(entries) != 1 {
		return PasswordEntry{}, fmt.Errorf("can not find key: %s", key)
	}
	return entries[0], e
}

// ListPasswords 列出所有保存的密码(需要解密)
func (ps *FileStorage) ListPasswords() ([]PasswordEntry, error) {
	return ps.getPasswords("")
}
