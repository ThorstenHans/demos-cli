package demo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/manifoldco/promptui"
)

type Config struct {
	JumpBoxAddress  string `json:"address"`
	JumpBoxPort     int    `json:"port"`
	JumpBoxUser     string `json:"user"`
	JumpBoxPassword string `json:"password"`
}

func Configure() error {
	cfg := &Config{}

	// Address (IPv4 validation)
	prompt := promptui.Prompt{
		Label: "JumpBox IPv4 Address",
		Validate: func(input string) error {
			ip := net.ParseIP(input)
			if ip == nil || ip.To4() == nil {
				return fmt.Errorf("invalid IPv4 address")
			}
			return nil
		},
	}
	addr, err := prompt.Run()
	if err != nil {
		return err
	}
	cfg.JumpBoxAddress = addr

	// Port (default 22)
	portPrompt := promptui.Prompt{
		Label:   "JumpBox Port",
		Default: "22",
		Validate: func(input string) error {
			port, err := strconv.Atoi(input)
			if err != nil || port < 1 || port > 65535 {
				return fmt.Errorf("invalid port number")
			}
			return nil
		},
	}
	portStr, err := portPrompt.Run()
	if err != nil {
		return err
	}
	port, _ := strconv.Atoi(portStr)
	cfg.JumpBoxPort = port

	// Username
	userPrompt := promptui.Prompt{
		Label: "JumpBox Username",
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("username cannot be empty")
			}
			return nil
		},
	}
	user, err := userPrompt.Run()
	if err != nil {
		return err
	}
	cfg.JumpBoxUser = user

	// Password (masked, with confirmation)
	passPrompt := promptui.Prompt{
		Label: "JumpBox Password",
		Mask:  '*',
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("password cannot be empty")
			}
			return nil
		},
	}
	pass1, err := passPrompt.Run()
	if err != nil {
		return err
	}

	passPrompt.Label = "Confirm Password"
	pass2, err := passPrompt.Run()
	if err != nil {
		return err
	}
	if pass1 != pass2 {
		return fmt.Errorf("passwords do not match")
	}

	cfg.JumpBoxPassword = pass1
	return SaveConfig(cfg)
}

// Salt is compiled into the binary via ldflags at build time.
// It must be a string for ldflags to work.
var Salt string

// getConfigPath builds the path to the config file in a subfolder of the home directory
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".demo")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "demo.config"), nil
}

// deriveKey creates a 32-byte AES key from the salt
func deriveKey() []byte {
	hash := sha256.Sum256([]byte(Salt))
	return hash[:]
}

// encrypt encrypts data using AES-GCM
func encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(deriveKey())
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM
func decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(deriveKey())
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("invalid ciphertext")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// SaveConfig serializes and encrypts the config struct, then writes it to disk
func SaveConfig(cfg *Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(jsonData)
	if err != nil {
		return err
	}

	path, err := getConfigPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, encrypted, 0600)
}

// LoadConfig reads, decrypts, and unmarshals the config from disk
func LoadConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	decrypted, err := decrypt(data)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(decrypted, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
