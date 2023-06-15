package auth

import (
	"github.com/asmir-a/langlearn/backend/auth/dbwrappers"
	"github.com/asmir-a/langlearn/backend/httperrors"
)

func Logout(currentSessionKey string) *httperrors.HttpError {
	//todo: maybe should check if the session is valid
	err := dbwrappers.DeleteSession(currentSessionKey)
	if err != nil {
		return httperrors.NewHttp500Error(err)
	}
	return nil
}
