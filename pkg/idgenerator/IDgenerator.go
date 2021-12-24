package idgenerator

import (
	"crypto/rand"
	"fmt"
	"log"
)

func CreateID() *string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	id := fmt.Sprintf("%x%x",
		b[0:4], b[4:6])
	return &id

}
