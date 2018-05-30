package paperfishGo

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
)

func NewTrackId() (string, error) {
	var b []byte
	var err error

	b = make([]byte, 256)
	_, err = rand.Read(b)
	if err != nil {
		Goose.Fetch.Logf(1, "Error generating random string")
		return "", err
	}

	return fmt.Sprintf("PF%x", sha256.Sum256(b)), nil
}
