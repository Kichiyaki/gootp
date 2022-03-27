package andotp_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Kichiyaki/gootp/internal/andotp"
	"github.com/stretchr/testify/assert"
)

type entry struct {
	Algorithm     string   `json:"algorithm"`
	Digits        uint8    `json:"digits"`
	Issuer        string   `json:"issuer"`
	Label         string   `json:"label"`
	LastUsed      uint64   `json:"last_used"`
	Period        uint32   `json:"period"`
	Secret        string   `json:"secret"`
	Tags          []string `json:"tags"`
	Thumbnail     string   `json:"thumbnail"`
	Type          string   `json:"type"`
	UsedFrequency uint64   `json:"usedFrequency"`
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	entries := []entry{
		{
			Algorithm:     "SHA1",
			Digits:        6,
			Issuer:        "TestIssuer",
			Label:         "TestLabel",
			LastUsed:      uint64(time.Now().Unix()),
			Period:        30,
			Secret:        "secret",
			Thumbnail:     "Default",
			Type:          "TOTP",
			UsedFrequency: 0,
		},
	}
	entriesJSON, err := json.Marshal(entries)
	assert.Nil(t, err)
	password := []byte("password22231")

	encrypted, err := andotp.Encrypt(password, entriesJSON)
	assert.Nil(t, err)

	decrypted, err := andotp.Decrypt(password, encrypted)
	assert.Nil(t, err)

	var result []entry
	err = json.Unmarshal(decrypted, &result)
	assert.Nil(t, err)
	assert.Equal(t, entries, result)
}
