package provider

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

func GenerateHash(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not compute file '%s' checksum: %s", filename, err)
	}
	h := sha1.New()
	h.Write(data)
	sha := hex.EncodeToString(h.Sum(nil))

	return sha, nil
}
