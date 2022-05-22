package internal_test

import (
	"encoding/json"
	"testing"

	"github.com/Kichiyaki/gootp/internal"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	entries := []internal.Entry{
		{
			Algorithm: "SHA1",
			Digits:    6,
			Issuer:    "TestIssuer",
			Label:     "TestLabel",
			Period:    30,
			Secret:    "secret",
			Thumbnail: "Default",
			Type:      "TOTP",
		},
	}
	entriesJSON, err := json.Marshal(entries)
	assert.Nil(t, err)
	password := []byte("password22231")

	encrypted, err := internal.Encrypt(entriesJSON, password)
	assert.Nil(t, err)

	decrypted, err := internal.Decrypt(encrypted, password)
	assert.Nil(t, err)
	var result []internal.Entry
	err = json.Unmarshal(decrypted, &result)
	assert.Nil(t, err)
	assert.Equal(t, entries, result)

	decryptedEntries, err := internal.DecryptAsEntries(encrypted, password)
	assert.Nil(t, err)
	assert.Equal(t, entries, decryptedEntries)
}
