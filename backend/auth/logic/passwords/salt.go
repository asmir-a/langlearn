package passwords

import (
	"crypto/rand"
	"encoding/base64"
)

const SALT_LENGTH = 128

func Salt(username string) (randomSalt string, err error) {
	randomSaltBytes := make([]byte, SALT_LENGTH)
	_, err = rand.Read(randomSaltBytes)
	randomSaltString := base64.StdEncoding.EncodeToString(randomSaltBytes)
	return username + randomSaltString, nil
}
