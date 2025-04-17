package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration values
type Config struct {
	DBPath        string
	JWTSecret     string
	Port          int
	EncryptionKey string
	RateLimit     int
}

// Default values for when environment variables are not set
const (
	DefaultJWTSecretKey     = "key-place-holder-for-production"
	DefaultJWTExpiration    = 15 * time.Minute
	DefaultJWTRefreshExpiry = 24 * time.Hour
	DefaultDatabasePath     = "./data.db"
	DefaultServerPort       = 8080
	DefaultEncryptionKey    = "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
)

func New() *Config {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = DefaultServerPort
	}

	rateLimit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if rateLimit == 0 {
		rateLimit = 100 // Default to 100 requests per minute
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = DefaultDatabasePath
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		encryptionKey = DefaultEncryptionKey
	}

	return &Config{
		DBPath:        dbPath,
		JWTSecret:     os.Getenv("JWT_SECRET"),
		Port:          port,
		EncryptionKey: encryptionKey,
		RateLimit:     rateLimit,
	}
}

func (c *Config) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(c.EncryptionKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *Config) Decrypt(encrypted string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(c.EncryptionKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
