package passwords

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/asmir-a/langlearn/backend/httperrors"
)

const SALT_LENGTH = 128

func Salt(username string) (string, *httperrors.HttpError) { //todo: maybe using bytes instead of strings in function signature is better
	randomSaltBytes := make([]byte, SALT_LENGTH)
	if _, err := rand.Read(randomSaltBytes); err != nil {
		return "", httperrors.NewHttp500Error(err)
	}
	randomSaltString := base64.StdEncoding.EncodeToString(randomSaltBytes)
	return username + randomSaltString, nil
}
