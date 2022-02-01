package idgenerator

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func CreateID(size int) string {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	id := hex.EncodeToString(b)
	return id

}
