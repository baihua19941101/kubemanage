package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

const encryptedPrefix = "enc:v1:"

type sensitiveCodec struct {
	enabled bool
	key     []byte
}

func newSensitiveCodec(secretKey string) *sensitiveCodec {
	trimmed := strings.TrimSpace(secretKey)
	if trimmed == "" {
		return &sensitiveCodec{enabled: false}
	}
	sum := sha256.Sum256([]byte(trimmed))
	return &sensitiveCodec{
		enabled: true,
		key:     sum[:],
	}
}

func (c *sensitiveCodec) Encrypt(value string) (string, error) {
	if !c.enabled || strings.TrimSpace(value) == "" {
		return value, nil
	}
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("build aes cipher failed: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("build gcm failed: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce failed: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *sensitiveCodec) Decrypt(value string) (string, error) {
	if !c.enabled || strings.TrimSpace(value) == "" {
		return value, nil
	}
	if !strings.HasPrefix(value, encryptedPrefix) {
		return value, nil
	}
	raw := strings.TrimPrefix(value, encryptedPrefix)
	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext failed: %w", err)
	}
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("build aes cipher failed: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("build gcm failed: %w", err)
	}
	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt ciphertext failed: %w", err)
	}
	return string(plaintext), nil
}
