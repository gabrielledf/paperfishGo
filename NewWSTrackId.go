package paperfishGo

import (
	"crypto/sha256"
	"math/rand"
	"time"
)

func NewWSTrackId() (uint32, error) {
	var buf []byte
	var c byte
	var err error
	var i int
	var n uint32
	var sum [sha256.Size]byte
	var rnd *rand.Rand

	rnd = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	buf = make([]byte, 256)
	_, err = rnd.Read(buf)
	if err != nil {
		Goose.Fetch.Logf(1, "Error generating random string")
		return 0, err
	}

	sum = sha256.Sum256(buf)
	for i, c = range sum {
		n ^= uint32(c) << (uint(i%4) * 8)
	}

	return n, nil
}
