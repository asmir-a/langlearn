package passwords

import (
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

const ARGON_TIME = 1
const ARGON_MEMORY = 64 * 1024
const ARGON_THREADS = 4
const ARGON_KEYLEN = 128

func Hash(password string, salt string) string {
	hashBytes := argon2.IDKey([]byte(password), []byte(salt), ARGON_TIME, ARGON_MEMORY, ARGON_THREADS, ARGON_KEYLEN)
	hashString := base64.StdEncoding.EncodeToString(hashBytes)
	return hashString
}
