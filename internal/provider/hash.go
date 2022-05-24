package provider

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

func GenerateHash(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not compute file '%s' checksum: %s", filename, err)
	}
	h := sha1.New()
	h.Write(data)
	sha1 := hex.EncodeToString(h.Sum(nil))

	return sha1, nil
}
