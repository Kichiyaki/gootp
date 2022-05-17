package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"

	"golang.org/x/crypto/pbkdf2"
)

const (
	ivLen         = 12
	keyLen        = 32
	iterationLen  = 4
	saltLen       = 12
	maxIterations = 160000
	minIterations = 140000
)

func Encrypt(password, plaintext []byte) ([]byte, error) {
	iter := make([]byte, iterationLen)
	iv := make([]byte, ivLen)
	salt := make([]byte, saltLen)

	iterations, err := randIterations()
	if err != nil {
		return nil, fmt.Errorf("randIterations: %w", err)
	}

	binary.BigEndian.PutUint32(iter, uint32(iterations))

	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("rand.Read(iv): %w", err)
	}

	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("rand.Read(salt): %w", err)
	}

	secretKey := pbkdf2.Key(password, salt, iterations, keyLen, sha1.New)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}

	cipherText := aesGCM.Seal(nil, iv, plaintext, nil)

	finalCipher := make([]byte, 0, len(iter)+len(salt)+len(iv)+len(cipherText))
	finalCipher = append(finalCipher, iter...)
	finalCipher = append(finalCipher, salt...)
	finalCipher = append(finalCipher, iv...)
	finalCipher = append(finalCipher, cipherText...)

	return finalCipher, nil
}

func randIterations() (int, error) {
	iterations, err := rand.Int(rand.Reader, big.NewInt(int64(maxIterations-minIterations)))
	if err != nil {
		return 0, fmt.Errorf("rand.Int: %w", err)
	}

	return int(iterations.Int64() + minIterations), nil
}

func Decrypt(password, text []byte) ([]byte, error) {
	iterations := text[:iterationLen]
	salt := text[iterationLen : iterationLen+saltLen]
	iv := text[iterationLen+saltLen : iterationLen+saltLen+ivLen]
	cipherText := text[iterationLen+saltLen+ivLen:]
	iter := int(binary.BigEndian.Uint32(iterations))
	secretKey := pbkdf2.Key(password, salt, iter, keyLen, sha1.New)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}

	plaintext, err := aesGCM.Open(nil, iv, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("aesGCM.Open: %w", err)
	}

	return plaintext, nil
}

type Entry struct {
	Algorithm     string   `json:"algorithm"`
	Digits        uint     `json:"digits"`
	Issuer        string   `json:"issuer"`
	Label         string   `json:"label"`
	LastUsed      uint     `json:"last_used"`
	Period        uint     `json:"period"`
	Secret        string   `json:"secret"`
	Tags          []string `json:"tags"`
	Thumbnail     string   `json:"thumbnail"`
	Type          string   `json:"type"`
	UsedFrequency uint     `json:"usedFrequency"`
}

func DecryptAsEntries(password, text []byte) ([]Entry, error) {
	result, err := Decrypt(password, text)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	var entries []Entry
	if err := json.Unmarshal(result, &entries); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return entries, nil
}
