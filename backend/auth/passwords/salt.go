package passwords

import (
	"crypto/rand"
	"log"
)

const SALT_LENGTH = 128

func Salt(username string) (randomSalt string, err error) {
	randomSaltBytes := make([]byte, SALT_LENGTH)
	_, err = rand.Read(randomSaltBytes)
	if err != nil {
		log.Println("rand.Read failed")
	}
	randomSalt = username + string(randomSaltBytes)
	return
}
